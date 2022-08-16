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

package jwt

import (
	"github.com/go-chassis/go-chassis/v2/security/token"
	"github.com/go-chassis/openlog"
	"net/http"
	"time"
)

var auth *Auth

// Auth should implement auth logic
// it is singleton
type Auth struct {
	SecretFunc token.SecretFunc //required
	Expire     time.Duration
	Realm      string //required

	//optional. Authorize check whether this request could access some resource or API based on json claims.
	//Typically, this method should communicate with a RBAC, ABAC system
	Authorize func(payload map[string]interface{}, req *http.Request) error

	//optional.
	// this function control whether a request should be validate or not
	// if this func is nil, validate all requests.
	MustAuth func(req *http.Request) bool
}

// Use put a custom auth logic
// then register handler to chassis
func Use(middleware *Auth) {
	auth = middleware
	if auth.Expire == 0 {
		openlog.Warn("token issued by service will not expire")
	}
	if auth.MustAuth == nil {
		openlog.Info("auth all requests")
	} else {
		openlog.Warn("under some condition, no auth")
	}
}

// SetExpire reset the expire time
func SetExpire(duration time.Duration) {
	auth.Expire = duration
}
