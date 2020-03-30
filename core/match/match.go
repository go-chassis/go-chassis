/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package match

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/go-mesh/openlogging"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
	"sync"
)

var matches sync.Map

//Method decide value Match expression oor not
type Method func(value, expression string) bool

var matchPlugin = map[string]Method{
	"exact":    exact,
	"contains": contains,
	"regex":    regex,
}

//Install a strategy
func Install(name string, m Method) {
	matchPlugin[name] = m
}

func mark(inv *invocation.Invocation) {
	matchName := ""
	matches.Range(func(k, v interface{}) bool {
		mp, ok := v.(*config.MatchPolicy)
		if ok {
			if isMatch(inv, mp) {
				if name, ok := k.(string); ok {
					matchName = name
					return false
				}
			}
		}
		return true
	})
	if matchName != "" {
		//openlogging.GetLogger().Infof("the invocation math policy %s", matchName)
		inv.Mark(matchName)
	}
}

func isMatch(inv *invocation.Invocation, matchPolicy *config.MatchPolicy) bool {
	if !headsMatch(inv.Headers(), matchPolicy.Headers) {
		return false
	}

	req, err := httputil.HTTPRequest(inv)
	if err != nil {
		openlogging.Warn("get request error: " + err.Error())
		return false
	}

	if len(matchPolicy.APIPaths) != 0 && !apiMatch(req.URL.Path, matchPolicy.APIPaths) {
		return false
	}

	if matchPolicy.Method != "" && strings.ToUpper(matchPolicy.Method) != req.Method {
		return false
	}
	return true
}

func apiMatch(apiPath string, apiPolicy map[string]string) bool {
	for strategy, exp := range apiPolicy {
		if ok, _ := Match(strategy, apiPath, exp); ok {
			return true
		}
	}
	return false
}

func headsMatch(headers map[string]string, headPolicy map[string]map[string]string) bool {

	for key, policy := range headPolicy {
		val := headers[key]
		if val == "" {
			return false
		}
		for strategy, exp := range policy {
			if o, err := Match(strategy, val, exp); err != nil || !o {
				return false
			}
		}
	}
	return true
}

//Match compare value and expression
func Match(strategy, value, expression string) (bool, error) {
	f, ok := matchPlugin[strategy]
	if !ok {
		return false, fmt.Errorf("invalid Match method")
	}
	return f(value, expression), nil
}

func exact(value, express string) bool {
	return value == express
}

func contains(value, express string) bool {
	return strings.Contains(value, express)
}

func regex(value, express string) bool {
	reg := regexp.MustCompilePOSIX(express)
	if !reg.Match([]byte(value)) {
		return false
	}
	return true
}

//SaveMatchPolicy saves Match policy
func SaveMatchPolicy(value string, k string, name string) {
	m := &config.MatchPolicy{}
	err := yaml.Unmarshal([]byte(value), m)
	if err != nil {
		openlogging.Warn("invalid policy:" + k)
		return
	}
	openlogging.GetLogger().Debugf("get policy %s %v", name, m)
	matches.Store(name, m)
}
