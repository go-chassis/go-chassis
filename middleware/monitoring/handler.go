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

package monitoring

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/status"
	"github.com/go-chassis/go-chassis/v2/pkg/metrics"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/openlog"
	"net/http"
	"time"
)

// errors
const (
	MetricsLatency = "scb_request_process_duration"
	MetricsRequest = "scb_request_count"
	MetricsErrors  = "scb_error_response_count"
	Name           = "monitoring"
)

var labels = []string{"service", "instance", "version", "app", "env"}
var labels4Resp = []string{"service", "instance", "version", "app", "env", "code"}
var labelMap map[string]string

// Handler monitor server side metrics, the key metrics is latency, QPS, Errors, do not use it in consumer chain
type Handler struct {
}

// Handle record metrics
func (ph *Handler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	start := time.Now()
	switch i.Args.(type) {
	case *http.Request:
		err := metrics.CounterAdd(MetricsRequest, 1, labelMap)
		if err != nil {
			openlog.Error("can not monitor:" + err.Error())
		}
	case *restful.Request:
		err := metrics.CounterAdd(MetricsRequest, 1, labelMap)
		if err != nil {
			openlog.Error("can not monitor:" + err.Error())
		}
	default:
		//skip monitoring
		chain.Next(i, cb)
		return
	}
	chain.Next(i, func(resp *invocation.Response) {
		if resp.Status >= status.Status(i.Protocol, status.InternalServerError) {
			m := map[string]string{
				"service":  runtime.ServiceName,
				"instance": runtime.InstanceID,
				"version":  runtime.Version,
				"app":      runtime.App,
				"env":      runtime.Environment,
				"code":     fmt.Sprintf("%d", resp.Status),
			}
			err := metrics.CounterAdd(MetricsErrors, 1, m)
			if err != nil {
				openlog.Error(err.Error())
			}
		}
		duration := time.Since(start)
		err := metrics.SummaryObserve(MetricsLatency, float64(duration.Milliseconds()), labelMap)
		if err != nil {
			openlog.Error(err.Error())
		}
		cb(resp)
	})

}
func newHandler() handler.Handler {
	err := metrics.CreateCounter(metrics.CounterOpts{
		Name:   MetricsRequest,
		Labels: labels,
	})
	if err != nil {
		openlog.Fatal(err.Error())
	}
	err = metrics.CreateSummary(metrics.SummaryOpts{
		Name:       MetricsLatency,
		Labels:     labels,
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
	if err != nil {
		openlog.Fatal(err.Error())
	}
	err = metrics.CreateCounter(metrics.CounterOpts{
		Name:   MetricsErrors,
		Labels: labels4Resp,
	})
	if err != nil {
		openlog.Fatal(err.Error())
	}
	labelMap = map[string]string{
		"service":  runtime.ServiceName,
		"instance": runtime.InstanceID,
		"version":  runtime.Version,
		"app":      runtime.App,
		"env":      runtime.Environment,
	}
	return &Handler{}
}

// Name returns the router string
func (ph *Handler) Name() string {
	return Name
}
func init() {
	err := handler.RegisterHandler(Name, newHandler)
	if err != nil {
		openlog.Error(err.Error())
	}
}
