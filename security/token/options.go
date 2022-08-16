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

type SigningMethod int

// const
const (
	RS256 SigningMethod = 1
	RS512 SigningMethod = 2
	HS256 SigningMethod = 3
)

// Options is options
type Options struct {
	Expire        string
	SigningMethod SigningMethod
}

// Option is option
type Option func(options *Options)

// WithExpTime generate a token which expire after a duration
// for example 5s,1m,24h
func WithExpTime(exp string) Option {
	return func(options *Options) {
		options.Expire = exp
	}
}

// WithSigningMethod specify the sign method
func WithSigningMethod(m SigningMethod) Option {
	return func(options *Options) {
		options.SigningMethod = m
	}
}
