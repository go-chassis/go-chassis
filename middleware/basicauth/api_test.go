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

package basicauth

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"testing"
)

type fakeHandler struct{}

func (h *fakeHandler) Name() string {

	return "fake"
}

func (h *fakeHandler) Handle(*handler.Chain, *invocation.Invocation, invocation.ResponseCallBack) {
	log.Println("authorized")
	return
}

func new() handler.Handler {
	return &fakeHandler{}
}
func TestUseBasicAuth(t *testing.T) {
	Use(&BasicAuth{
		Realm: "test-realm",
		Authenticate: func(u, p string) error {
			return nil
		},
	})

	handler.RegisterHandler("basicAuth", newBasicAuth)
	handler.RegisterHandler("fake", new)

	c, err := handler.CreateChain(common.Provider, "default", []string{"basicAuth", "fake"}...)
	t.Run("Invalid", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Add("Authorization", "QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
		inv := &invocation.Invocation{
			Args: req,
		}
		c.Next(inv, func(ir *invocation.Response) {
			err = ir.Err
			assert.Error(t, err)
		})
	})

	t.Run("normal", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Add("Authorization", "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==")
		inv := &invocation.Invocation{
			Args: req,
		}
		c.Next(inv, func(ir *invocation.Response) {
			err = ir.Err
			assert.NoError(t, err)
		})
	})
}
