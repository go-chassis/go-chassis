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

	"github.com/go-chassis/go-chassis/core/config"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/metrics"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/prometheus/client_golang/prometheus"
	gometrics "github.com/rcrowley/go-metrics"
	"regexp"
)

// DefaultPrometheusSinker variable for default prometheus configurations
var DefaultPrometheusSinker *PrometheusSinker
var onceInit sync.Once

const (
	regex       = "(Provider|Consumer)\\.(.*)"
	regexSource = "\\.(.+)\\.(Provider|Consumer)\\.(.*)"
)

var desc = map[string]string{
	"Consumer.attempts":          "all requests count to target provider service",
	"Consumer.errors":            "if a request is failed because of error or timeout or circuit open, it will increase",
	"Consumer.successes":         "if request is not timeout or failed, it will increase",
	"Consumer.failures":          "if request return error, it will increase",
	"Consumer.rejects":           "if circuit open, then all request will be reject immediately, it will increas",
	"Consumer.shortCircuits":     "after circuit open, it will increase",
	"Consumer.timeouts":          "after timeout, it will increase ",
	"Consumer.fallbackSuccesses": "if fallback is executed and no error returns, it will increase",
	"Consumer.fallbackFailures":  "if fallback is executed and error returns, it will increase",
	"Consumer.totalDuration":     "how long all requests consumed totally",
	"Consumer.runDuration":       "how long a request consumed",

	"Provider.attempts":          "all requests count to target provider service",
	"Provider.errors":            "if a request is failed because of error or timeout or circuit open, it will increase",
	"Provider.successes":         "if request is not timeout or failed, it will increase",
	"Provider.failures":          "if request return error, it will increase",
	"Provider.rejects":           "if circuit open, then all request will be reject immediately, it will increas",
	"Provider.shortCircuits":     "after circuit open, it will increase",
	"Provider.timeouts":          "after timeout, it will increase ",
	"Provider.fallbackSuccesses": "if fallback is executed and no error returns, it will increase",
	"Provider.fallbackFailures":  "if fallback is executed and error returns, it will increase",
	"Provider.totalDuration":     "how long all requests consumed totally",
	"Provider.runDuration":       "how long a request consumed",
}

//GetDesc retrieve metric doc
func GetDesc(name string) string {
	h := desc[name]
	if h != "" {
		return h
	}
	return name
}

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
			Help: desc[name],
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
			Help: GetDesc(name),
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
		metricName, sn, operationID, schemaID := ExtractMetricKey(name)
		promLabels := prometheus.Labels{"hostname": runtime.HostName, "self": runtime.ServiceName, "target": sn, "appID": runtime.App, "version": runtime.Version, "schemaID": schemaID, "operationID": operationID}
		for k, v := range runtime.MD {
			promLabels[k] = v
		}
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
				c.gaugeVecFromNameAndValue(fmt.Sprintf("%s.%s", metricName, "mean"), meanTime, promLabels)
				c.gaugeVecFromNameAndValue(strings.Replace(metricName, "runDuration", "qps", -1), t.RateMean(), promLabels)
				promLabels["quantile"] = "0.05"
				c.gaugeVecFromNameAndValue(metricName, ps[0]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.25"
				c.gaugeVecFromNameAndValue(metricName, ps[1]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.5"
				c.gaugeVecFromNameAndValue(metricName, ps[2]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.75"
				c.gaugeVecFromNameAndValue(metricName, ps[3]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.90"
				c.gaugeVecFromNameAndValue(metricName, ps[4]/float64(time.Millisecond), promLabels)
				promLabels["quantile"] = "0.99"
				c.gaugeVecFromNameAndValue(metricName, ps[5]/float64(time.Millisecond), promLabels)
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

//ExtractServiceSchemaOperationMetrics parse service,schema and operation
//key example Microservice.SchemaID.OperationId.metrics
func ExtractServiceSchemaOperationMetrics(raw string) (target, schemaID, operation, metrics string) {
	metrics = getEventType(raw)
	tokens := strings.Split(raw, ".")
	switch len(tokens) {
	case 2:
		target = tokens[0]
	case 3:
		target = tokens[0]
		schemaID = tokens[1]
	case 4:
		target = tokens[0]
		schemaID = tokens[1]
		operation = tokens[2]
	}
	return
}

//ExtractMetricKey return metrics related infors
//example Consumer.ErrServer.rest./sayhimessage.rejects
//the first and last string consist of metrics name
//second is ErrServer
//3th and 4th is schema and operation
func ExtractMetricKey(key string) (source string, target string, schema string, op string) {
	regNormal := regexp.MustCompile(regex)
	regSource := regexp.MustCompile(regexSource)
	var raw, role string
	if regNormal.MatchString(key) {
		s := regNormal.FindStringSubmatch(key)
		role = s[1]
		raw = s[2]
	}
	if regSource.MatchString(key) {
		s := regNormal.FindStringSubmatch(key)
		source = s[1]
		role = s[2]
		raw = s[3]
	}
	sn, scID, opID, m := ExtractServiceSchemaOperationMetrics(raw)

	return role + "." + m, sn, scID, opID
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
