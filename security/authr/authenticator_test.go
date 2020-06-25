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

package authr_test

import (
	"context"
	"errors"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/security/authr"
	"github.com/stretchr/testify/assert"
	"testing"
)

type insecureAuthenticator struct {
}

func newInsecureAuth(opts *authr.Options) (authr.Authenticator, error) {
	return &insecureAuthenticator{}, nil
}
func (a *insecureAuthenticator) Login(ctx context.Context, user string, password string) (string, error) {
	if user == archaius.GetString("username", "") && password == archaius.GetString("password", "") {
		return "token", nil
	}
	return "", errors.New("wrong credential")
}
func (a *insecureAuthenticator) Authenticate(ctx context.Context, token string) (interface{}, error) {
	return archaius.GetString("username", ""), nil
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("username", "admin")
	archaius.Set("password", "admin")
	authr.Install("default", newInsecureAuth)
	authr.Init()
	token, _ := authr.Login(ctx, "admin", "admin")
	assert.Equal(t, "token", token)
	claims, _ := authr.Authenticate(ctx, "token")
	assert.Equal(t, "admin", claims)

}
