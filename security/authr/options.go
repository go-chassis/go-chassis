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

package authr

// Options is a struct to stores options
type Options struct {
	Plugin string
}

// Option is option
type Option func(options *Options)

// WithPlugin specify plugin name
func WithPlugin(p string) Option {
	return func(options *Options) {
		options.Plugin = p
	}
}

// Options is a struct to stores options
type LoginOptions struct {
	ExpireAfter string
}

// Option is option
type LoginOption func(options *LoginOptions)

// ExpireAfter specify time duration, for example: 3d, 3m, 1s, 3h
func ExpireAfter(p string) LoginOption {
	return func(options *LoginOptions) {
		options.ExpireAfter = p
	}
}
