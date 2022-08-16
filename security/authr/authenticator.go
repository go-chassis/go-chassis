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

// Package authr defines a standard interface to decouple with specify auth solution.
// it also decouple user name and password from authentication action
package authr

import (
	"context"
	"errors"
	"fmt"
)

var defaultAuthenticator Authenticator

// errors
var (
	ErrNoImpl = errors.New("no implementation")
)

type newFunc func(opts *Options) (Authenticator, error)

var plugins = make(map[string]newFunc)

// Install install a Plugin
func Install(name string, f newFunc) {
	plugins[name] = f
}

// Authenticator can sign a token and authenticate that token
type Authenticator interface {
	Login(ctx context.Context, user string, password string, opts ...LoginOption) (string, error)
	Authenticate(ctx context.Context, token string) (interface{}, error)
}

// Login verify a user info and return a token
func Login(ctx context.Context, user string, password string, opts ...LoginOption) (string, error) {
	return defaultAuthenticator.Login(ctx, user, password, opts...)
}

// Authenticate parse a token and return the claims in that token
func Authenticate(ctx context.Context, token string) (interface{}, error) {
	return defaultAuthenticator.Authenticate(ctx, token)
}

// Init initiate this module
func Init(opts ...Option) error {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	if o.Plugin == "" {
		o.Plugin = "default"
	}
	f, ok := plugins[o.Plugin]
	if !ok {
		return fmt.Errorf("plugin is no installed: %s", o.Plugin)
	}
	var err error
	defaultAuthenticator, err = f(o)
	return err
}
