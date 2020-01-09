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

package tracing

import (
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/go-mesh/openlogging"

	skycom "github.com/SkyAPM/go2sky/reporter/grpc/common"
	"strconv"
)

//for skywalking use
const (
	HTTPPrefix             = "http://"
	CrossProcessProtocolV2 = "Sw6"
	DefaultTraceContext    = ""
)

//component id for skywalking which is used for topology
const (
	HTTPClientComponentID = 2
	HTTPServerComponentID = 49
)

//SkyWalkingClient for connecting and reporting to skywalking server
type SkyWalkingClient struct {
	reporter    go2sky.Reporter
	tracer      *go2sky.Tracer
	ServiceType int32
}

//CreateEntrySpan create entry span
func (s *SkyWalkingClient) CreateEntrySpan(sc *SpanContext) (go2sky.Span, error) {
	openlogging.Debug("CreateEntrySpan begin. span" + sc.OperationName)
	span, ctx, err := s.tracer.CreateEntrySpan(sc.Ctx, sc.OperationName, func() (string, error) {
		if sc.ParTraceCtx != nil {
			return sc.ParTraceCtx[CrossProcessProtocolV2], nil
		}
		return DefaultTraceContext, nil
	})
	if err != nil {
		openlogging.Error("CreateEntrySpan error:" + err.Error())
		return nil, err
	}
	span.Tag(go2sky.TagHTTPMethod, sc.Method)
	span.Tag(go2sky.TagURL, sc.URL)
	span.SetSpanLayer(skycom.SpanLayer_Http)
	span.SetComponent(s.ServiceType)
	sc.Ctx = ctx
	return span, nil
}

//CreateExitSpan create end span
func (s *SkyWalkingClient) CreateExitSpan(sc *SpanContext) (go2sky.Span, error) {
	openlogging.Debug("CreateExitSpan begin. span:" + sc.OperationName)
	span, err := s.tracer.CreateExitSpan(sc.Ctx, sc.OperationName, sc.Peer, func(header string) error {
		sc.TraceCtx[CrossProcessProtocolV2] = header
		return nil
	})
	if err != nil {
		openlogging.Error("CreateExitSpan error:" + err.Error())
		return nil, err
	}
	span.Tag(go2sky.TagHTTPMethod, sc.Method)
	span.Tag(go2sky.TagURL, sc.URL)
	span.SetSpanLayer(skycom.SpanLayer_Http)
	span.SetComponent(s.ServiceType)
	return span, nil
}

//EndSpan make span end and report to skywalking
func (s *SkyWalkingClient) EndSpan(sp go2sky.Span, statusCode int) error {
	openlogging.Debug("EndSpan status:" + strconv.Itoa(statusCode))
	span, ok := (sp).(go2sky.Span)
	if !ok || span == nil {
		openlogging.Error("EndSpan failed.")
		return nil
	}
	span.Tag(go2sky.TagStatusCode, strconv.Itoa(statusCode))
	span.End()
	return nil
}

//NewApmClient init report and tracer for connecting and sending messages to skywalking server
func NewApmClient(op TracingOptions) (TracingClient, error) {
	var (
		err    error
		client SkyWalkingClient
	)
	client.reporter, err = reporter.NewGRPCReporter(op.ServerURI)
	if err != nil {
		openlogging.Error("NewGRPCReporter error:" + err.Error())
		return nil, err
	}
	client.tracer, err = go2sky.NewTracer(op.MicServiceName, go2sky.WithReporter(client.reporter))
	//not wait for register here
	//t.WaitUntilRegister()
	if err != nil {
		openlogging.Error("NewTracer error:" + err.Error())
		return nil, err
	}
	client.ServiceType = int32(op.MicServiceType)
	openlogging.Debug("NewApmClient skywalking success:" + "uri:" + op.ServerURI)
	return &client, nil
}

//Init apm client
func Init(op TracingOptions) {
	openlogging.Info("apm Init " + op.ServerURI)
	var err error
	tc, err = NewApmClient(op)
	if err != nil {
		openlogging.Error("NewAPMClient error:" + err.Error())
		return
	}
}
