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

package marker

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/openlog"
	"gopkg.in/yaml.v2"
	"net/http"
	"strings"
	"sync"
)

const (
	Once       = "once"
	PerService = "perService"
)

var matches sync.Map

//Operate decide value match expression or not
type Operate func(value, expression string) bool

var operatorPlugin = map[string]Operate{
	"exact":     exact,
	"contains":  contains,
	"regex":     regex,
	"noEqu":     noEqu,
	"less":      less,
	"noLess":    noLess,
	"greater":   greater,
	"noGreater": noGreater,
}

//Install a strategy
func Install(name string, m Operate) {
	operatorPlugin[name] = m
}

//Mark mark an invocation with matchName by match policy
func Mark(inv *invocation.Invocation) {
	matchName := ""
	policy := "once"
	matches.Range(func(k, v interface{}) bool {
		mps, ok := v.(*config.MatchPolicies)
		if !ok {
			return true
		}
		for _, mp := range mps.Matches {
			if isMatch(inv, mp) {
				if name, ok := k.(string); ok {
					matchName = name
					policy = mp.TrafficMarkPolicy
					return false
				}
			}
		}

		return true
	})
	if matchName != "" {
		//the invocation math policy
		if policy == Once {
			inv.SetHeader(common.HeaderMark, matchName)
		}
		inv.Mark(matchName)
	}
}

func isMatch(inv *invocation.Invocation, matchPolicy config.MatchPolicy) bool {
	if !headsMatch(inv.Headers(), matchPolicy.Headers) {
		return false
	}
	var req *http.Request
	switch r := inv.Args.(type) {
	case *http.Request:
		req = r
	case *restful.Request:
		req = r.Request
	default:
		return false
	}

	if len(matchPolicy.APIPaths) != 0 && !apiMatch(req.URL.Path, matchPolicy.APIPaths) {
		return false
	}
	if len(matchPolicy.Method) != 0 {
		if !methodMatch(req.Method, matchPolicy.Method) {
			return false
		}
	}
	return true
}
func methodMatch(reqMethod string, methods []string) bool {
	matchMethod := false
	for _, m := range methods {
		if strings.ToUpper(reqMethod) == m {
			matchMethod = true
		}
	}
	return matchMethod
}
func apiMatch(apiPath string, apiPolicy map[string]string) bool {
	if len(apiPolicy) == 0 {
		return true
	}

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

//match compare value and expression
func Match(operator, value, expression string) (bool, error) {
	f, ok := operatorPlugin[operator]
	if !ok {
		return false, fmt.Errorf("invalid match method")
	}
	return f(value, expression), nil
}

//SaveMatchPolicy saves match policy
func SaveMatchPolicy(name, value string, k string) error {
	m := &config.MatchPolicies{}
	err := yaml.Unmarshal([]byte(value), m)
	if err != nil {
		openlog.Error("invalid policy " + k + ":" + err.Error())
		return err
	}
	openlog.Info("add match policy", openlog.WithTags(openlog.Tags{
		"module": "marker",
		"event":  "update",
	}))
	matches.Store(name, m)
	return nil
}

//Policy return policy
func Policy(name string) *config.MatchPolicies {
	i, ok := matches.Load(name)
	if !ok {
		return nil
	}
	m, ok := i.(*config.MatchPolicies)
	if !ok {
		return nil
	}
	return m
}
