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

package ratelimiter

import (
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/resilience/rate"
	"github.com/go-mesh/openlogging"
	"math"
)

func init() {
	err := handler.RegisterHandler(Name, newRateLimiterHandler)
	if err != nil {
		openlogging.Error(err.Error())
	}
}

//Handler can only be used in server(provider) side
type Handler struct{}

// Handle limit request rate according to marker
func (h *Handler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	if inv.GetMark() == "" { //if some user do not use invocation marker feature, then should skip rate limiter
		chain.Next(inv, cb)
		return
	}
	if rate.GetRateLimiters().TryAccept(inv.GetMark(), math.MaxInt32, 0) {
		chain.Next(inv, cb)
		return
	}
	r := newErrResponse(inv)
	cb(r)
}

// Name returns name
func (h *Handler) Name() string {
	return Name
}
func newRateLimiterHandler() handler.Handler {
	return &Handler{}
}
