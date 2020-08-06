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

package handler

import (
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/marker"
)

//TrafficMarker
const (
	TrafficMarker = "traffic-marker"
)

//MarkHandler compares the match rule with invocation and mark this invocation
type MarkHandler struct {
}

//Name return the handler name
func (m *MarkHandler) Name() string {
	return TrafficMarker
}

//Handle to handle the mart invocation
func (m *MarkHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	marker.Mark(i)
	chain.Next(i, cb)
}

func newMarkHandler() Handler {
	return &MarkHandler{}
}
