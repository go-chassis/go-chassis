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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/status"
	"github.com/go-chassis/go-chassis/security/token"
	restfulserver "github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
)

//errors
var (
	ErrNoHeader    = errors.New("no authorization in header")
	ErrInvalidAuth = errors.New("invalid authentication")
)

//Handler is is a jwt interceptor
type Handler struct {
}

//Handle intercept unauthorized request
func (h *Handler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	var req *http.Request
	if r, ok := i.Args.(*http.Request); ok {
		req = r
	} else if r, ok := i.Args.(*restful.Request); ok {
		req = r.Request
	} else {
		openlogging.Error(fmt.Sprintf("this handler only works for http request, wrong type: %t", i.Args))
		return
	}
	if mustAuth(req) {
		v := req.Header.Get(restfulserver.HeaderAuth)
		if v == "" {
			handler.WriteBackErr(ErrNoHeader, status.Status(i.Protocol, status.Unauthorized), cb)
			return
		}
		s := strings.Split(v, " ")
		if len(s) != 2 {
			handler.WriteBackErr(ErrNoHeader, status.Status(i.Protocol, status.Unauthorized), cb)
			return
		}
		to := s[1]
		payload, err := token.DefaultManager.Verify(to, auth.SecretFunc)
		if err != nil {
			openlogging.Error("can not parse jwt:" + err.Error())
			handler.WriteBackErr(ErrNoHeader, status.Status(i.Protocol, status.Unauthorized), cb)
			return
		}
		if auth.Authorize != nil {
			err = auth.Authorize(payload, req)
			if err != nil {
				handler.WriteBackErr(ErrNoHeader, status.Status(i.Protocol, status.Unauthorized), cb)
				return
			}
		}
	} else {
		openlogging.Info("skip auth")
	}

	chain.Next(i, cb)
}
func mustAuth(req *http.Request) bool {
	if auth.MustAuth == nil {
		return true
	}
	return auth.MustAuth(req)
}
func newHandler() handler.Handler {
	return &Handler{}
}

// Name returns the router string
func (h *Handler) Name() string {
	return "jwt"
}
func init() {
	err := handler.RegisterHandler("jwt", newHandler)
	if err != nil {
		openlogging.Error(err.Error())
	}
}
