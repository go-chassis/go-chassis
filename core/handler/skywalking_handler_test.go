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

package handler_test

import (
	"context"
	"github.com/go-chassis/go-chassis/core/apm"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/stretchr/testify/assert"
	"testing"
)

//initConfig
func initConfig() {
	config.MonitorCfgDef = &model.MonitorCfg{ServiceComb: model.ServiceCombStruct{APM: model.APMStruct{Tracing: model.TracingStruct{Tracer: "skywalking", Settings: map[string]string{"URI": "127.0.0.1:11800", "enable": "true"}}}}}
	config.MicroserviceDefinition = &model.MicroserviceCfg{ServiceDescription: model.MicServiceStruct{Name: "skywalking"}}
}

//initApm
func initApm() {
	apm.Init()
}

//initInv
func initInv() *invocation.Invocation {
	i := invocation.New(context.Background())
	i.MicroServiceName = "test"
	i.Endpoint = "calculator"
	i.URLPathFormat = "/bmi"
	i.SetHeader("Sw6", "")
	return i
}

//TestProviderHandlerName
func TestProviderHandlerName(t *testing.T) {
	h := handler.SkyWalkingProviderHandler{}
	assert.Equal(t, h.Name(), handler.SkyWalkingProvider)
}

//TestNewProvier
func TestNewProvier(t *testing.T) {
	h := handler.SkyWalkingProviderHandler{}
	assert.NotEqual(t, h, nil)
	assert.Equal(t, h.Name(), handler.SkyWalkingProvider)
}

//TestProvierHandle
func TestProvierHandle(t *testing.T) {
	initConfig()
	initApm()
	c := handler.Chain{}
	c.AddHandler(&handler.SkyWalkingProviderHandler{})
	c.Next(initInv(), func(r *invocation.Response) error {
		assert.Equal(t, r.Err, nil)
		return r.Err
	})
}

//TestConsumerHandlerName
func TestConsumerHandlerName(t *testing.T) {
	c := handler.SkyWalkingConsumerHandler{}
	assert.Equal(t, c.Name(), handler.SkyWalkingConsumer)
}

//TestNewConsumer
func TestNewConsumer(t *testing.T) {
	h := handler.SkyWalkingConsumerHandler{}
	assert.NotEqual(t, h, nil)
	assert.Equal(t, h.Name(), handler.SkyWalkingConsumer)
}

//TestConsumerHandle
func TestConsumerHandle(t *testing.T) {
	initConfig()
	initApm()
	c := handler.Chain{}
	c.AddHandler(&handler.SkyWalkingConsumerHandler{})
	c.Next(initInv(), func(r *invocation.Response) error {
		assert.Equal(t, r.Err, nil)
		return r.Err
	})
}
