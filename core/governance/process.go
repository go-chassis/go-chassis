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

package governance

import (
	"github.com/go-chassis/go-chassis/core/match"
	"github.com/go-chassis/go-chassis/pkg/rate"
	"github.com/go-mesh/openlogging"
	"gopkg.in/yaml.v2"
	"strings"
)

//ProcessMatch saves all policy to match module
//then match module is able to mark invocation
func ProcessMatch(key string, value string) {
	s := strings.Split(key, ".")
	if len(s) != 3 {
		openlogging.Warn("invalid key:" + key)
		return
	}
	name := s[2]
	match.SaveMatchPolicy(value, key, name)
}

type limiterPolicy struct {
	Matcher string `json:"match"`
	Quota   int    `json:"quota"`
}

//ProcessLimiter saves limiter, after a invocation is marked,
//go chassis will get correspond limiter with mark name
func ProcessLimiter(key string, value string) {
	s := strings.Split(key, ".")
	if len(s) != 3 {
		openlogging.Warn("invalid key:" + key)
		return
	}
	policy := &limiterPolicy{}
	err := yaml.Unmarshal([]byte(value), policy)
	if err != nil {
		openlogging.Error("invalid limiter: " + key)
		return
	}

	//key is match rule name, value is qps
	rate.GetRateLimiters().UpdateRateLimit(policy.Matcher, policy.Quota)
}
