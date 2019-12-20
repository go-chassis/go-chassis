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

package reporter

// Forked from github.com/deathowl
// Some parts of this file have been modified to make it functional in this package
import (
	"github.com/go-chassis/go-archaius"
	"strings"
	"sync"
	"time"

	"fmt"
	"github.com/go-chassis/go-chassis/pkg/circuit"
	m "github.com/go-chassis/go-chassis/pkg/metrics"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/prometheus/client_golang/prometheus"
)

var onceInit sync.Once

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

var FlushInterval time.Duration //interval to update prom metrics
var gauges map[string]prometheus.Gauge
var gaugeVecs map[string]*prometheus.GaugeVec

// GetPrometheusSinker get prometheus configurations
func GetPrometheusSinker() {
	onceInit.Do(func() {
		t, err := time.ParseDuration(archaius.GetString("cse.metrics.flushInterval", "10s"))
		if err != nil {
			t = time.Second * 10
		}
		FlushInterval = t
		gauges = make(map[string]prometheus.Gauge)
		gaugeVecs = make(map[string]*prometheus.GaugeVec)
	})
}

func flattenKey(key string) string {
	key = strings.Replace(key, " ", "_", -1)
	key = strings.Replace(key, ".", "_", -1)
	key = strings.Replace(key, "-", "_", -1)
	key = strings.Replace(key, "=", "_", -1)
	return key
}

func gaugeVecFromNameAndValue(name string, val float64, labels prometheus.Labels) {
	var labelNames []string
	for labelName := range labels {
		labelNames = append(labelNames, labelName)
	}
	gVec, ok := gaugeVecs[name]
	if !ok {
		gVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: flattenKey(name),
			Help: GetDesc(name),
		}, labelNames)
		m.GetSystemPrometheusRegistry().MustRegister(gVec)
		gaugeVecs[name] = gVec
	}
	gVec.With(labels).Set(val)
}

//ReportMetricsToPrometheus report metrics to prometheus registry, you can use GetSystemPrometheusRegistry to get prometheus registry. by default chassis will report system metrics to prometheus
func ReportMetricsToPrometheus(cb *hystrix.CircuitBreaker) error {
	GetPrometheusSinker()

	now := time.Now()
	attemptsName := cb.Name + ".attempts"
	errorsName := cb.Name + ".errors"
	successesName := cb.Name + ".successes"
	failuresName := cb.Name + ".failures"
	rejectsName := cb.Name + ".rejects"
	shortCircuitsName := cb.Name + ".shortCircuits"
	timeoutsName := cb.Name + ".timeouts"
	fallbackSuccessesName := cb.Name + ".fallbackSuccesses"
	fallbackFailuresName := cb.Name + ".fallbackFailures"
	totalDurationName := cb.Name + ".totalDuration"
	runDurationName := cb.Name + ".runDuration"
	_, sn, operationID, schemaID := circuit.ParseCircuitCMD(errorsName)
	promLabels := prometheus.Labels{"hostname": runtime.HostName, "self": runtime.ServiceName,
		"target": sn, "appID": runtime.App, "version": runtime.Version,
		"schemaID": schemaID, "operationID": operationID}
	for k, v := range runtime.MD {
		promLabels[k] = v
	}
	metricName := circuit.GetMetricsName(errorsName)
	errCount := cb.Metrics.DefaultCollector().Errors().Sum(now)
	gaugeVecFromNameAndValue(metricName, errCount, promLabels)

	attemptsCount := cb.Metrics.DefaultCollector().NumRequests().Sum(now)
	metricName = circuit.GetMetricsName(attemptsName)
	gaugeVecFromNameAndValue(metricName, attemptsCount, promLabels)

	successesCount := cb.Metrics.DefaultCollector().Successes().Sum(now)
	metricName = circuit.GetMetricsName(successesName)
	gaugeVecFromNameAndValue(metricName, successesCount, promLabels)

	failureCount := cb.Metrics.DefaultCollector().Failures().Sum(now)
	metricName = circuit.GetMetricsName(failuresName)
	gaugeVecFromNameAndValue(metricName, failureCount, promLabels)

	rejectCount := cb.Metrics.DefaultCollector().Rejects().Sum(now)
	metricName = circuit.GetMetricsName(rejectsName)
	gaugeVecFromNameAndValue(metricName, rejectCount, promLabels)

	scCount := cb.Metrics.DefaultCollector().ShortCircuits().Sum(now)
	metricName = circuit.GetMetricsName(shortCircuitsName)
	gaugeVecFromNameAndValue(metricName, scCount, promLabels)

	timeoutCount := cb.Metrics.DefaultCollector().Timeouts().Sum(now)
	metricName = circuit.GetMetricsName(timeoutsName)
	gaugeVecFromNameAndValue(metricName, timeoutCount, promLabels)

	fbsCount := cb.Metrics.DefaultCollector().FallbackSuccesses().Sum(now)
	metricName = circuit.GetMetricsName(fallbackSuccessesName)
	gaugeVecFromNameAndValue(metricName, fbsCount, promLabels)

	fbfCount := cb.Metrics.DefaultCollector().FallbackFailures().Sum(now)
	metricName = circuit.GetMetricsName(fallbackFailuresName)
	gaugeVecFromNameAndValue(metricName, fbfCount, promLabels)

	latencyTotalMean := cb.Metrics.DefaultCollector().TotalDuration().Mean()
	metricName = circuit.GetMetricsName(totalDurationName)
	gaugeVecFromNameAndValue(fmt.Sprintf("%s.%s", metricName, "mean"),
		float64(latencyTotalMean), promLabels)

	runDuration := cb.Metrics.DefaultCollector().RunDuration()
	metricName = circuit.GetMetricsName(runDurationName)
	gaugeVecFromNameAndValue(fmt.Sprintf("%s.%s", metricName, "mean"),
		float64(runDuration.Mean()), promLabels)
	promLabels["quantile"] = "0.05"
	gaugeVecFromNameAndValue(metricName, float64(runDuration.Percentile(5)), promLabels)
	promLabels["quantile"] = "0.25"
	gaugeVecFromNameAndValue(metricName, float64(runDuration.Percentile(25)), promLabels)
	promLabels["quantile"] = "0.5"
	gaugeVecFromNameAndValue(metricName, float64(runDuration.Percentile(5)), promLabels)
	promLabels["quantile"] = "0.75"
	gaugeVecFromNameAndValue(metricName, float64(runDuration.Percentile(75)), promLabels)
	promLabels["quantile"] = "0.90"
	gaugeVecFromNameAndValue(metricName, float64(runDuration.Percentile(90)), promLabels)
	promLabels["quantile"] = "0.99"
	gaugeVecFromNameAndValue(metricName, float64(runDuration.Percentile(99)), promLabels)
	return nil
}
func init() {
	hystrix.InstallReporter("Prometheus", ReportMetricsToPrometheus)
}
