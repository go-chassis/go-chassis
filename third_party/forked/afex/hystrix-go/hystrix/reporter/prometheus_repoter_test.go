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

package reporter

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/third_party/forked/afex/hystrix-go/hystrix"
	"testing"
	"time"
)

func TestReportMetricsToPrometheus(t *testing.T) {
	hystrix.Do("cmd", func() error {
		time.Sleep(20 * time.Millisecond)
		return nil
	}, nil)
	hystrix.Do("cmd", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}, nil)
	hystrix.Do("cmd", func() error {
		time.Sleep(2 * time.Millisecond)
		return nil
	}, nil)
	time.Sleep(1 * time.Second)
	cb, _, _ := hystrix.GetCircuit("cmd")
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.metrics.flushInterval", "1s")
	time.Sleep(1 * time.Second)
	ReportMetricsToPrometheus(cb)

}
