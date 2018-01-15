package metrics

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rcrowley/go-metrics"
	"os"
	"strings"
	"sync"
	"time"
)

// DefaultPrometheusSinker variable for default prometheus configurations
var DefaultPrometheusSinker *PrometheusSinker
var once sync.Once

// PrometheusSinker is the struct for prometheus configuration parameters
type PrometheusSinker struct {
	Registry      metrics.Registry      // Registry to be exported
	PromRegistry  prometheus.Registerer //Prometheus registry
	FlushInterval time.Duration         //interval to update prom metrics
	gauges        map[string]prometheus.Gauge
	gaugeVecs     map[string]*prometheus.GaugeVec
}

// GetPrometheusSinker get prometheus configurations
func GetPrometheusSinker(mr metrics.Registry, pr *prometheus.Registry) *PrometheusSinker {
	once.Do(func() {
		DefaultPrometheusSinker = NewPrometheusProvider(mr, pr, time.Second)
	})
	return DefaultPrometheusSinker
}

// NewPrometheusProvider returns the object of prometheus configurations
func NewPrometheusProvider(r metrics.Registry, promRegistry prometheus.Registerer, FlushInterval time.Duration) *PrometheusSinker {
	return &PrometheusSinker{
		Registry:      r,
		PromRegistry:  promRegistry,
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
		c.PromRegistry.MustRegister(g)
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
		c.PromRegistry.MustRegister(gVec)
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
		promLabels := prometheus.Labels{"hostname": hostName, "servicename": config.SelfServiceName, "appID": config.GlobalDefinition.AppID, "version": config.SelfVersion, "schemaID": schemaID, "operationID": operationID}
		switch metric := i.(type) {
		case metrics.Counter:
			c.gaugeVecFromNameAndValue(metricName, float64(metric.Count()), promLabels)
		case metrics.Gauge:
			c.gaugeVecFromNameAndValue(metricName, float64(metric.Value()), promLabels)
		case metrics.GaugeFloat64:
			c.gaugeVecFromNameAndValue(metricName, float64(metric.Value()), promLabels)
		case metrics.Histogram:
			samples := metric.Snapshot().Sample().Values()
			if len(samples) > 0 {
				lastSample := samples[len(samples)-1]
				c.gaugeVecFromNameAndValue(metricName, float64(lastSample), promLabels)
			}
		case metrics.Meter:
			lastSample := metric.Snapshot().Rate1()
			c.gaugeVecFromNameAndValue(metricName, float64(lastSample), promLabels)
		case metrics.Timer:
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
func (c *PrometheusSinker) EnableRunTimeMetrics() {
	c.PromRegistry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
	c.PromRegistry.MustRegister(prometheus.NewGoCollector())
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
