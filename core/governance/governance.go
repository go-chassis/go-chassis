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
	"github.com/go-chassis/go-archaius"
	"github.com/go-mesh/openlogging"
	"strings"
)

//prefix const
const (
	KindMatchPrefix        = "servicecomb.match"
	KindRateLimitingPrefix = "servicecomb.rateLimiting"
)

var processFuncMap = map[string]ProcessFunc{
	//build-in
	KindMatchPrefix:        ProcessMatch,
	KindRateLimitingPrefix: ProcessLimiter,
}

//ProcessFunc process a config
type ProcessFunc func(key string, value string)

//InstallProcessor install a func to process config,
//if a config key matches the key prefix, then the func will process the config
func InstallProcessor(keyPrefix string, process ProcessFunc) {
	processFuncMap[keyPrefix] = process
}

//Init go through all governance configs
//and call process func according to key prefix
func Init() {
	configMap := archaius.GetConfigs()
	openlogging.Info("process all governance rules")
	for k, v := range configMap {
		value, ok := v.(string)
		if !ok {
			openlogging.Warn("not string format,key:" + k)
		}
		openlogging.Debug(k + ":" + value)
		for prefix, f := range processFuncMap {
			if strings.HasPrefix(k, prefix) {
				f(k, value)
				break
			}
		}
	}
}
