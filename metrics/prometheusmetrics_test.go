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
package metrics_test

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	m "github.com/ServiceComb/go-chassis/metrics"
	"github.com/prometheus/client_golang/prometheus"
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

func TestPrometheusConfig_UpdatePrometheusMetricsOnce(t *testing.T) {
	config.GlobalDefinition = new(model.GlobalCfg)
	config.GlobalDefinition.Cse.Metrics.EnableGoRuntimeMetrics = false
	prometheusConfig := m.GetPrometheusSinker(m.GetSystemRegistry(), m.GetSystemPrometheusRegistry())
	t.Log("registering various metric types to go-metrics registry")
	c, _ := m.GetSystemRegistry().GetOrRegister("server.attempts", makeCounterFunc).(metrics.Counter)
	c.Inc(1)
	c, _ = m.GetSystemRegistry().GetOrRegister("server.successes", makeCounterFunc).(metrics.Counter)
	c.Inc(1)
	timer, _ := m.GetSystemRegistry().GetOrRegister("server.totalDuration", makeTimerFunc).(metrics.Timer)
	timer.Update(time.Millisecond * 10)
	g, _ := m.GetSystemRegistry().GetOrRegister("server.memory", makeGaugeFunc).(metrics.Gauge)
	g.Update(1200)
	gFloat := m.GetSystemRegistry().GetOrRegister("server.thread", makeGaugeFloat64Func).(metrics.GaugeFloat64)
	gFloat.Update(341)
	h, _ := m.GetSystemRegistry().GetOrRegister("server.requestDuration", makeHistogramFunc).(metrics.Histogram)
	h.Update(23)
	meter, _ := m.GetSystemRegistry().GetOrRegister("foo", makeMeterFunc).(metrics.Meter)
	meter.Mark(12)
	prometheusConfig.UpdatePrometheusMetricsOnce()
	metricsGatherer := prometheusConfig.PromRegistry.(prometheus.Gatherer)
	metricsFamilies, _ := metricsGatherer.Gather()
	assert.Equal(t, len(metricsFamilies), 6)
}
