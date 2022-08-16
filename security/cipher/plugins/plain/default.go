/*Copyright 2017 Huawei Technologies Co., Ltd
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *    http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package plain created on 2017/6/22
package plain

import (
	"github.com/go-chassis/cari/security"
	"github.com/go-chassis/go-chassis/v2/security/cipher"
)

// DefaultCipher is a struct
type DefaultCipher struct {
}

func init() {
	cipher.InstallCipherPlugin("default", new)
}
func new() security.Cipher {

	return &DefaultCipher{}
}

// Encrypt is method used for encryption
func (c *DefaultCipher) Encrypt(src string) (string, error) {
	return src, nil
}

// Decrypt is method used for decryption
func (c *DefaultCipher) Decrypt(src string) (string, error) {
	return src, nil
}
