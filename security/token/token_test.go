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
	"github.com/go-chassis/go-chassis/security/token"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJWTTokenManager_GetToken(t *testing.T) {
	s := []byte(`-----BEGIN RSA PRIVATE KEY-----
	MIIEowIBAAKCAQEAsDAuUCgHK78p0y8bkBPtSEBe5TdyhBqFL+DQskxu79z8n1o2
	W1n8cF60UYrh5JdOGpEAJXJIu3UJy0EcIyBsp8zbCCfLj6vQMLLDPe0R5AZSyCTH
	janPp1TSvBMEh7G7bC2HOzB+KsBgZfH+v3FO0ZyGX9oNbBFs+ShlVMxENsgzqFtF
	z/dUqeWZwwWGLzNedkYo3zeya0idtPQ83MrGPYe3IvT66IMFBiRtxWBHipepNLXl
	6/T5USDpWay/oXXmsbKf/2IIxmE+BjMU0gOW6JylYLoIxoF1IwXp2gyna+4RNNc6
	FlNeEFNfdcd3+hyeMaobFgr+WVY3nXYIYg0TkQIDAQABAoIBAAyuWxb/2oxGhQ8j
	K9ux43k40Nu0ovRpKD7q8npyz+VJxZD+oDzw/B9mYZog4eNfFIsK9rS7RgrgAKV1
	eT35/ngRYY5zts4Pcruekjjp0EjWP60SIJ7MoxqLG2PYBpJxs2i02i/jbKFNGWMd
	CNXkpOSnXHCXtDGcC3jfdHOnBB3hKlDL2pzdLZJHInDhT2xKXIbr187kB8dyOwVN
	zUzSFD8nUHYWWpC/x1FBh9tFKQrCdr8uqRnvqzfbCVRKS91ZywF60szehqA2r63p
	xJ+KfEg568+pGD4+FINsgYezVn9pGw1qKqrF89NULXf7Mm5FqkRpoJasbmsbZmVt
	ZNZCC6UCgYEA5usTw3FPQbptJABwLczKiMOwU4XrGa8OHWcyCaI2AYRFQK3Cy3xj
	YUMslk0tW4C5yZOJwzmOyOW28L8xf7D4uk8QXw5OwRhlvLVylrerIm1FQ/8yMNI1
	iKT5dsc8gqLgPqWdGfuGejHq7xL7fzbhNPTVEf8jYdji+8G7O6ysV4MCgYEAw1NG
	/6TzKZim0sf4RCmnJg8LWLN0EfPbSHu8hee95x6N+cYSLwlwVpXSfyA5M1sjbCiO
	8vWnhiP6fKEDrj0jqOyaGQOsPssOIEyaXVw5vyzoVpqFaBTzXZFCgbc1nwTiDeeR
	0vzc7AxuJgZj2yi3ETRfGlThzJySoT9UbcB0qFsCgYBBY4PfLjDhTecl8LHTZlBb
	1f4SSLPAPB/lF5nFvJdKaqgpnoqwkHKb0ifID+auKI9zk0HJdH0ISnQ5TAq6O+TS
	7RyXrjeC2mPEwiTGpQ/i2cppbNRLmtrp7L1vcw+hdnnFg6Qu/VihNY1vUZLB/Upc
	co/7XqIoTQBJhhx803Kh/QKBgEPQCMk+mlFptxlc5bu8flR/SqAsBXMqJ4p9sxEG
	SO8Rs5bxBmUgMMlO0LrkFBfZX23wktiVIuk2WoOkXyPCBDxkkId4t/dBBhF+puUc
	3MubqrpOgVyGUYu9n8prMgmYZ2cOa5lFwumM0z0OYOK4uv4VIaOBrrcb8OhclVJZ
	S+cbAoGBAJSDKtIypXlKs1qWKBjSkirrGAgpAJu4DrY3hKnZI2a0nEQE+ATNRnq+
		0XnzS6XkkiP4UyevJXVxCRlF+DHiBxd6f54yo7EIKUo+h5QnBqoGXqot2EWS7Kt2
	coVIagbb2GU1llZANN2vSF7twF/N0EUoSaIocgROhDVfQrhzGyW2
	-----END RSA PRIVATE KEY-----`)
	to, err := token.DefaultManager.GetToken(map[string]interface{}{
		"username": "peter",
	}, s)
	assert.NoError(t, err)
	t.Log(to)
	m, err := token.DefaultManager.ParseToken(to, s)
	assert.NoError(t, err)
	assert.Equal(t, "peter", m["username"])
	t.Run("with exp", func(t *testing.T) {
		to, err := token.DefaultManager.GetToken(map[string]interface{}{
			"username": "peter",
		}, []byte("my secret"), token.WithExpTime("1s"))
		assert.NoError(t, err)
		t.Log(to)
	})
}
