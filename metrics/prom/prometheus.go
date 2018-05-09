// Copyright 2016 Csergő Bálint github.com/deathowl
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prom

// Forked from github.com/deathowl
// Some parts of this file have been modified to make it functional in this package
import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ServiceComb/go-chassis/core/config"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/metrics"
	"github.com/prometheus/client_golang/prometheus"
	gometrics "github.com/rcrowley/go-metrics"
)

// DefaultPrometheusSinker variable for default prometheus configurations
var DefaultPrometheusSinker *PrometheusSinker
var onceInit sync.Once

// PrometheusSinker is the struct for prometheus configuration parameters
type PrometheusSinker struct {
	Registry      gometrics.Registry // Registry to be exported
	FlushInterval time.Duration      //interval to update prom metrics
	gauges        map[string]prometheus.Gauge
	gaugeVecs     map[string]*prometheus.GaugeVec
}

// GetPrometheusSinker get prometheus configurations
func GetPrometheusSinker(mr gometrics.Registry) *PrometheusSinker {
	onceInit.Do(func() {
		t, err := time.ParseDuration(config.GlobalDefinition.Cse.Metrics.FlushInterval)
		if err != nil {
			t = time.Second * 10
		}
		DefaultPrometheusSinker = NewPrometheusProvider(mr, t)
	})
	return DefaultPrometheusSinker
}

// NewPrometheusProvider returns the object of prometheus configurations
func NewPrometheusProvider(r gometrics.Registry, FlushInterval time.Duration) *PrometheusSinker {
	return &PrometheusSinker{
		Registry:      r,
		FlushInterval: FlushInterval,
		gauges:        make(map[string]prometheus.Gauge),
		gaugeVecs:     make(map[string]*prometheus.GaugeVec),
	}
}

func (c *PrometheusSinker) flattenKey(key string) string {
	key = strings.Replace(key, " ", "_", -1)
	key = strings.Replace(key, ".", "_", -1)
	key = strings.Replace(key, "-", "_", -1)
	key = strings.Replace(key, "=", "_", -1)
	return key
}

func (c *PrometheusSinker) gaugeFromNameAndValue(name string, val float64) {
	g, ok := c.gauges[name]
	if !ok {
		g = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: c.flattenKey(name),
			Help: name,
		})
		metrics.GetSystemPrometheusRegistry().MustRegister(g)
		c.gauges[name] = g
	}
	g.Set(val)
}

func (c *PrometheusSinker) gaugeVecFromNameAndValue(name string, val float64, labels prometheus.Labels) {
	var labelNames []string
	for labelName := range labels {
		labelNames = append(labelNames, labelName)
	}
	gVec, ok := c.gaugeVecs[name]
	if !ok {
		gVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: c.flattenKey(name),
			Help: name,
		}, labelNames)
		metrics.GetSystemPrometheusRegistry().MustRegister(gVec)
		c.gaugeVecs[name] = gVec
	}
	gVec.With(labels).Set(val)
}

// UpdatePrometheusMetrics update prometheus metrics
func (c *PrometheusSinker) UpdatePrometheusMetrics() {
	for range time.Tick(c.FlushInterval) {
		c.UpdatePrometheusMetricsOnce()
	}
}

// UpdatePrometheusMetricsOnce update prometheus metrics once
func (c *PrometheusSinker) UpdatePrometheusMetricsOnce() error {
	c.Registry.Each(func(name string, i interface{}) {
		metricName := extractMetricKey(name)
		operationID := extractOperationID(name)
		schemaID := extractSchemaID(name)
		hostName, _ := os.Hostname()
		promLabels := prometheus.Labels{"hostname": hostName, "service": config.SelfServiceName, "appID": config.GlobalDefinition.AppID, "version": config.SelfVersion, "schemaID": schemaID, "operationID": operationID}
		switch metric := i.(type) {
		case gometrics.Counter:
			c.gaugeVecFromNameAndValue(metricName, float64(metric.Count()), promLabels)
		case gometrics.Gauge:
			c.gaugeVecFromNameAndValue(metricName, float64(metric.Value()), promLabels)
		case gometrics.GaugeFloat64:
			c.gaugeVecFromNameAndValue(metricName, float64(metric.Value()), promLabels)
		case gometrics.Histogram:
			samples := metric.Snapshot().Sample().Values()
			if len(samples) > 0 {
				lastSample := samples[len(samples)-1]
				c.gaugeVecFromNameAndValue(metricName, float64(lastSample), promLabels)
			}
		case gometrics.Meter:
			lastSample := metric.Snapshot().Rate1()
			c.gaugeVecFromNameAndValue(metricName, float64(lastSample), promLabels)
		case gometrics.Timer:
			t := metric.Snapshot()
			ps := t.Percentiles([]float64{0.05, 0.25, 0.5, 0.75, 0.90, 0.99})
			switch getEventType(name) {
			case "runDuration":
				meanTime := t.Mean() / float64(time.Millisecond)
				key := strings.Replace(metricName, "runDuration", "request.duration.miliseconds", -1)
				c.gaugeVecFromNameAndValue(fmt.Sprintf("%s.%s", key, "mean"), meanTime, promLabels)
				c.gaugeVecFromNameAndValue(strings.Replace(metricName, "runDuration", "qps", -1), t.RateMean(), prometheus.Labels{"schemaID": schemaID, "operationID": operationID})
				promLabels["quantile"] = "0.05"
				c.gaugeVecFromNameAndValue(key, ps[0]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.25"
				c.gaugeVecFromNameAndValue(key, ps[1]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.5"
				c.gaugeVecFromNameAndValue(key, ps[2]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.75"
				c.gaugeVecFromNameAndValue(key, ps[3]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.90"
				c.gaugeVecFromNameAndValue(key, ps[4]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.99"
				c.gaugeVecFromNameAndValue(key, ps[5]/float64(time.Millisecond), promLabels)
			}
		}
	})
	return nil
}

// EnableRunTimeMetrics enable runtime metrics
func EnableRunTimeMetrics() {
	metrics.GetSystemPrometheusRegistry().MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
	metrics.GetSystemPrometheusRegistry().MustRegister(prometheus.NewGoCollector())
}

func getEventType(metricName string) string {
	tokens := strings.Split(metricName, ".")
	return tokens[len(tokens)-1]
}

func extractOperationID(key string) (operationID string) {
	var IsKeyHaveSourceName bool
	if !strings.HasPrefix(key, "Provider") && !strings.HasPrefix(key, "Consumer") {
		IsKeyHaveSourceName = true
	}
	tokens := strings.Split(key, ".")
	switch len(tokens) {
	case 3:
		return
	case 4:
		if IsKeyHaveSourceName {
			return
		}
		operationID = tokens[2]
		return
	case 5:
		if IsKeyHaveSourceName {
			operationID = tokens[3]
			return
		}
		operationID = tokens[3]
		return
	case 6:
		operationID = tokens[4]
		return
	}
	return
}

func extractSchemaID(key string) (schemaID string) {
	var IsKeyHaveSourceName bool
	if !strings.HasPrefix(key, "Provider") && !strings.HasPrefix(key, "Consumer") {
		IsKeyHaveSourceName = true
	}
	tokens := strings.Split(key, ".")
	switch len(tokens) {
	case 3:
		return
	case 4:
		if IsKeyHaveSourceName {
			return
		}
		schemaID = tokens[2]
		return
	case 5:
		if IsKeyHaveSourceName {
			schemaID = tokens[3]
			return
		}
		schemaID = tokens[2]
		return
	case 6:
		schemaID = tokens[3]
		return
	}
	return
}

func extractMetricKey(key string) string {
	key = strings.Replace(key, config.SelfServiceName, "service", 1)
	opID := extractOperationID(key)
	scID := extractSchemaID(key)
	if opID == "" && scID == "" {
		return key
	}

	key = strings.Replace(key, opID+".", "", 1)
	key = strings.Replace(key, scID+".", "", 1)

	return key
}

var onceEnable sync.Once

//ReportMetricsToPrometheus report metrics to prometheus registry, you can use GetSystemPrometheusRegistry to get prometheus registry. by default chassis will report system metrics to prometheus
func ReportMetricsToPrometheus(r gometrics.Registry) error {
	promConfig := GetPrometheusSinker(r)
	if archaius.GetBool("cse.metrics.enableGoRuntimeMetrics", true) {
		onceEnable.Do(func() {
			EnableRunTimeMetrics()
			lager.Logger.Info("Go Runtime Metrics is enabled")
		})

	}
	go promConfig.UpdatePrometheusMetrics()
	return nil
}
func init() {
	metrics.InstallReporter("Prometheus", ReportMetricsToPrometheus)
}
