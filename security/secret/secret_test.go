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

package secret_test

import (
	"github.com/go-chassis/go-chassis/security/secret"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenSecretKey(t *testing.T) {
	s, err := secret.GenRSAPrivateKey(4096)
	assert.NoError(t, err)
	t.Log(s)
}
func TestGenerateKeys(t *testing.T) {
	private, public, err := secret.GenRSAKeyPair(2048)
	assert.NoError(t, err)
	t.Log(private)
	t.Log(public)

	b, err := secret.RSAPrivate2Bytes(private)
	assert.NoError(t, err)
	t.Log(string(b))
	_, err = secret.ParseRSAPrivateKey(string(b))
	assert.NoError(t, err)
	b, err = secret.RSAPublicKey2Bytes(public)
	assert.NoError(t, err)
	t.Log(string(b))
	_, err = secret.ParseRSAPPublicKey(string(b))
	assert.NoError(t, err)
}
