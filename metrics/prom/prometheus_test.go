/*   Copyright 2016 Csergő Bálint github.com/deathowl
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package prom

// Forked from github.com/deathowl
// Some parts of this file have been modified to make it functional in this package
import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	m "github.com/go-chassis/go-chassis/metrics"
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	makeCounterFunc      = func() interface{} { return metrics.NewCounter() }
	makeTimerFunc        = func() interface{} { return metrics.NewTimer() }
	makeGaugeFunc        = func() interface{} { return metrics.NewGauge() }
	makeGaugeFloat64Func = func() interface{} { return metrics.NewGaugeFloat64() }
	makeHistogramFunc    = func() interface{} { return metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015)) }
	makeMeterFunc        = func() interface{} { return metrics.NewMeter() }
)

func TestPrometheusSinker_UpdatePrometheusMetrics(t *testing.T) {
	config.GlobalDefinition = new(model.GlobalCfg)
	config.GlobalDefinition.Cse.Metrics.EnableGoRuntimeMetrics = false
	prometheusSinker := GetPrometheusSinker(m.GetSystemRegistry())
	t.Log("registering various metric types to go-metrics registry")
	c, _ := m.GetSystemRegistry().GetOrRegister("Consumer.server.attempts", makeCounterFunc).(metrics.Counter)
	c.Inc(1)
	c, _ = m.GetSystemRegistry().GetOrRegister("Consumer.server.successes", makeCounterFunc).(metrics.Counter)
	c.Inc(1)
	timer, _ := m.GetSystemRegistry().GetOrRegister("Consumer.server.totalDuration", makeTimerFunc).(metrics.Timer)
	timer.Update(time.Millisecond * 10)
	g, _ := m.GetSystemRegistry().GetOrRegister("Consumer.server.memory", makeGaugeFunc).(metrics.Gauge)
	g.Update(1200)
	gFloat := m.GetSystemRegistry().GetOrRegister("Consumer.server.thread", makeGaugeFloat64Func).(metrics.GaugeFloat64)
	gFloat.Update(341)
	h, _ := m.GetSystemRegistry().GetOrRegister("Consumer.server.requestDuration", makeHistogramFunc).(metrics.Histogram)
	h.Update(23)
	meter, _ := m.GetSystemRegistry().GetOrRegister("foo", makeMeterFunc).(metrics.Meter)
	meter.Mark(12)
	prometheusSinker.UpdatePrometheusMetricsOnce()
	metricsFamilies, _ := m.GetSystemPrometheusRegistry().Gather()
	t.Log(metricsFamilies)
	assert.Equal(t, 6, len(metricsFamilies))
}

func TestExtractSchemaAndOperation(t *testing.T) {
	s, sch, op, m := ExtractServiceSchemaOperationMetrics("service.schema.op.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "schema", sch)
	assert.Equal(t, "op", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = ExtractServiceSchemaOperationMetrics("service.schema.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "schema", sch)
	assert.Equal(t, "", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = ExtractServiceSchemaOperationMetrics("service.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "", sch)
	assert.Equal(t, "", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = ExtractServiceSchemaOperationMetrics("service.schema.op.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "schema", sch)
	assert.Equal(t, "op", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = ExtractServiceSchemaOperationMetrics("ErrServer.rest./sayhimessage.metrics")
	assert.Equal(t, "ErrServer", s)
	assert.Equal(t, "rest", sch)
	assert.Equal(t, "/sayhimessage", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = ExtractServiceSchemaOperationMetrics("service.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "", sch)
	assert.Equal(t, "", op)
	assert.Equal(t, "metrics", m)

	s, sch, op, m = ExtractServiceSchemaOperationMetrics("service.schema.metrics")
	assert.Equal(t, "service", s)
	assert.Equal(t, "schema", sch)
	assert.Equal(t, "", op)
	assert.Equal(t, "metrics", m)
}

func TestExtractMetricKey(t *testing.T) {
	key, target, sch, op := ExtractMetricKey("Consumer.ErrServer.rest./sayhimessage.rejects")
	assert.Equal(t, "ErrServer", target)
	assert.Equal(t, "rest", sch)
	assert.Equal(t, "/sayhimessage", op)
	assert.Equal(t, "Consumer.rejects", key)
}
