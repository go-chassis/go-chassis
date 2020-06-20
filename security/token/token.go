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

package token

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-mesh/openlogging"
	"time"
)

//DefaultManager can be replaced
var DefaultManager Manager = &jwtTokenManager{}

//token pkg common errors
var (
	ErrInvalidExp = errors.New("expire time is illegal")
)

//jwt claims RFC 7519
//https://tools.ietf.org/html/rfc7519#section-4.1.2
const (
	JWTClaimsExp = "exp"
	JWTClaimsSub = "sub"
)

//GetToken gen token
func GetToken(claims map[string]interface{}, secret []byte, opts ...Option) (string, error) {
	return DefaultManager.GetToken(claims, secret, opts...)
}

//ParseToken return claims
func ParseToken(tokenString string, secret []byte) (map[string]interface{}, error) {
	return DefaultManager.ParseToken(tokenString, secret)
}

//Manager manages token
type Manager interface {
	GetToken(claims map[string]interface{}, secret []byte, option ...Option) (string, error)
	ParseToken(tokenString string, secret []byte) (map[string]interface{}, error)
}
type jwtTokenManager struct {
}

//GetToken gen token
func (f *jwtTokenManager) GetToken(claims map[string]interface{}, secret []byte, opts ...Option) (string, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	c := jwt.MapClaims{}
	c = claims
	if o.Expire != "" {
		d, err := time.ParseDuration(o.Expire)
		if err != nil {
			return "", ErrInvalidExp
		}
		c[JWTClaimsExp] = time.Now().Add(d).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString(secret)

}

//ParseToken return claims
func (f *jwtTokenManager) ParseToken(tokenString string, secret []byte) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
		return claims, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			openlogging.Error("not a valid jwt")
			return nil, err
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			openlogging.Error("token expired")
			return nil, err
		} else {
			openlogging.Error("parse token err:" + err.Error())
			return nil, err
		}
	} else {
		openlogging.Error("parse token err:" + err.Error())
		return nil, err
	}
}
