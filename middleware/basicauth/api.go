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

//Package basicauth supply basicAuth middleware abstraction
package basicauth

import (
	"github.com/go-chassis/go-chassis/core/handler"
	"net/http"
)

var auth *BasicAuth

//BasicAuth should implement basic auth server side logic
//it is singleton
type BasicAuth struct {
	Realm        string                                     //required
	Authorize    func(user, pwd string) error               //required
	Authenticate func(user string, req *http.Request) error //optional
}

//Use put a custom basic auth logic
//then register handler to chassis
func Use(middleware *BasicAuth) {
	auth = middleware
	handler.RegisterHandler("basicAuth", newBasicAuth)
}
