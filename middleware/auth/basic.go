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

package auth

import (
	"encoding/base64"
	"errors"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"net/http"
	"strings"
)

var (
	ErrInvalidBase64 = errors.New("invalid base64")
	ErrNoHeader      = errors.New("not authorized")
	ErrInvalidAuth   = errors.New("invalid authentication")
)

const HeaderAuth = "Authorization"

//BasicAuthHandler is is a handler
type BasicAuthHandler struct {
}

// Handle is to handle the router related things
func (ph *BasicAuthHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	if req, ok := i.Args.(*http.Request); ok {
		subject := req.Header.Get(HeaderAuth)
		if subject == "" {
			handler.WriteBackErr(ErrNoHeader, http.StatusUnauthorized, cb)
			return
		}
		u, p, err := decode(subject)
		if err != nil {
			openlogging.Error("can not decode base 64:" + err.Error())
			handler.WriteBackErr(ErrNoHeader, http.StatusUnauthorized, cb)
			return
		}
		err = auth.Authorize(u, p)
		if err != nil {
			handler.WriteBackErr(ErrNoHeader, http.StatusUnauthorized, cb)
			return
		}
		if auth.Authenticate != nil {
			err = auth.Authenticate(u, req)
			if err != nil {
				handler.WriteBackErr(ErrNoHeader, http.StatusUnauthorized, cb)
				return
			}
		}
	}
	chain.Next(i, cb)
}

func newBasicAuth() handler.Handler {
	return &BasicAuthHandler{}
}

// Name returns the router string
func (ph *BasicAuthHandler) Name() string {
	return "basicAuth"
}
func decode(subject string) (user string, pwd string, err error) {
	parts := strings.Split(subject, " ")
	if len(parts) != 2 {
		return "", "", ErrInvalidAuth

	}
	if parts[0] != "Basic" {
		return "", "", ErrInvalidAuth
	}
	s, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", "", ErrInvalidBase64
	}

	result := strings.Split(string(s), ":")
	if len(result) != 2 {
		return "", "", ErrInvalidAuth
	}

	return result[0], result[1], nil
}
