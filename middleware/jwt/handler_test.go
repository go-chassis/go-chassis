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
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/security/token"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
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
func TestUse(t *testing.T) {
	Use(&Auth{
		MustAuth: func(req *http.Request) bool {
			if strings.Contains(req.URL.Path, "/login") {
				return false
			}
			return true
		},
		Realm:     "test-realm",
		SecretKey: []byte("my_secret"),
	})

	handler.RegisterHandler("jwt", newHandler)
	handler.RegisterHandler("fake", new)
	to, _ := token.DefaultManager.GetToken(map[string]interface{}{
		"username": "peter",
	}, []byte("my_secret"))
	t.Log(to)
	c, err := handler.CreateChain(common.Provider, "default", []string{"jwt", "fake"}...)
	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api", nil)
		req.Header.Add("Authorization", "Bearer "+to)
		inv := &invocation.Invocation{
			Args: req,
		}
		c.Next(inv, func(ir *invocation.Response) error {
			err = ir.Err
			assert.NoError(t, err)
			return err
		})
	})
	t.Run("skip auth", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login", nil)
		req.Header.Add("Authorization", "Bearer "+to)
		inv := &invocation.Invocation{
			Args: req,
		}
		c.Next(inv, func(ir *invocation.Response) error {
			err = ir.Err
			assert.NoError(t, err)
			return err
		})
	})
}
