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
	"github.com/go-mesh/openlogging"
	"gopkg.in/yaml.v2"
	"sync"
)

var matches sync.Map

type Method func(value, expression string) bool

var matchPlugin = map[string]Method{
	"exact": exact,
}

//Install a strategy
func Install(name string, m Method) {
	matchPlugin[name] = m
}
func mark(inv *invocation.Invocation) {
	matches.Range(func(k, v interface{}) bool {
		return false
	})
}

//compare value and expression
func match(strategy, value, expression string) (bool, error) {
	f, ok := matchPlugin[strategy]
	if !ok {
		return false, fmt.Errorf("invalid match method")
	}
	return f(value, expression), nil
}
func exact(value, express string) bool {
	return value == express
}

//SaveMatchPolicy saves match policy
func SaveMatchPolicy(value string, k string, name string) {
	m := &config.MatchPolicy{}
	err := yaml.Unmarshal([]byte(value), m)
	if err != nil {
		openlogging.Warn("invalid policy:" + k)
	}
	matches.Store(name, m)
}
