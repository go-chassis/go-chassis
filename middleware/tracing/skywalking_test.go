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

package tracing_test

import (
	"context"
	"github.com/go-chassis/go-chassis/middleware/tracing"

	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	op        tracing.TracingOptions
	apmClient tracing.TracingClient
	sc        *tracing.SpanContext
)

func InitOption() {
	op = tracing.TracingOptions{
		ServerURI:      "192.168.88.64:8080",
		MicServiceName: "mesher",
		MicServiceType: 1}
}

func IniSpanContext() {
	sc = &tracing.SpanContext{
		Ctx:           context.Background(),
		OperationName: "test",
		ParTraceCtx:   map[string]string{},
		TraceCtx:      map[string]string{},
		Peer:          "test",
		Method:        "get",
		URL:           "/etc/url",
		ComponentID:   "1",
		SpanLayerID:   "11",
		ServiceName:   "mesher"}
}

func TestInit(t *testing.T) {
	InitOption()
	tracing.Init(op)
	op = tracing.TracingOptions{}
	tracing.Init(op)
}

func TestNewApmClient(t *testing.T) {
	InitOption()
	var err error
	apmClient, err = tracing.NewApmClient(op)
	assert.Equal(t, err, nil)
}

func TestCreateEntrySpan(t *testing.T) {
	InitOption()
	IniSpanContext()
	var err error
	apmClient, err = tracing.NewApmClient(op)
	assert.Equal(t, err, nil)
	span, err := apmClient.CreateEntrySpan(sc)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)

	op = tracing.TracingOptions{}
	apmClient, err = tracing.NewApmClient(op)
	assert.NotEqual(t, err, nil)
}

func TestCreateExitSpan(t *testing.T) {
	InitOption()
	IniSpanContext()
	var err error
	apmClient, err = tracing.NewApmClient(op)
	assert.Equal(t, err, nil)
	span, err := apmClient.CreateExitSpan(sc)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)

	op = tracing.TracingOptions{}
	apmClient, err = tracing.NewApmClient(op)
	assert.NotEqual(t, err, nil)
}

func TestEndSpan(t *testing.T) {
	InitOption()
	IniSpanContext()
	var err error
	apmClient, err = tracing.NewApmClient(op)
	assert.Equal(t, err, nil)
	span, err := apmClient.CreateEntrySpan(sc)
	assert.Equal(t, err, nil)
	assert.NotEqual(t, span, nil)
	err = apmClient.EndSpan(span, 1)
	assert.Equal(t, err, nil)

	err = apmClient.EndSpan(nil, 1)
	assert.Equal(t, err, nil)
}
