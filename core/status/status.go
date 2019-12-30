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

package status

import "net/http"

//status key const
const (
	Unauthorized = "Unauthorized"

	InternalServerError = "InternalServerError"
	ServiceUnavailable = "ServiceUnavailable"
	//TODO more status key
)

var defaultStatus = map[string]int{
	Unauthorized: http.StatusUnauthorized,

	InternalServerError: http.StatusInternalServerError,
	ServiceUnavailable: http.StatusServiceUnavailable,
	//TODO more default status
}

var protocolStatus = make(map[string]map[string]int, 1)

func init() {
	protocolStatus["rest"] = defaultStatus
}

//Register allows you custom a status map for a protocol plugin
func Register(protocol string, status map[string]int) {
	//TODO map key list check
	protocolStatus[protocol] = status

}

//Status return a status, if protocol do not has dedicated status map will use http status as standard map
func Status(protocol, statusKey string) int {
	s, ok := protocolStatus[protocol]
	if !ok {
		return defaultStatus[statusKey]
	}
	return s[statusKey]
}
