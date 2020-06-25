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

package token_test

import (
	"github.com/go-chassis/go-chassis/security/secret"
	"github.com/go-chassis/go-chassis/security/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJWTTokenManager_GetToken(t *testing.T) {
	privateKey, public, err := secret.GenRSAKeyPair(4096)
	assert.NoError(t, err)
	to, err := token.Sign(map[string]interface{}{
		"username": "peter",
	}, privateKey, token.WithSigningMethod(token.RS512))
	assert.NoError(t, err)
	t.Log(to)
	m, err := token.Verify(to, func(claims interface{}, method token.SigningMethod) (interface{}, error) {
		return public, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "peter", m["username"])
	t.Run("with exp", func(t *testing.T) {
		to, err := token.Sign(map[string]interface{}{
			"username": "peter",
		}, []byte("my secret"), token.WithExpTime("1s"))
		assert.NoError(t, err)
		t.Log(to)
	})
}
