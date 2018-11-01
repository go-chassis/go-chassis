package metrics

// Forked from github.com/afex/hystrix-go
// Some parts of this file have been modified to make it functional in this package

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix/metric_collector"
	"github.com/rcrowley/go-metrics"
)

// CircuitCollector is a struct to keeps metric information of Http requests
type CircuitCollector struct {
	attempts          string
	errors            string
	successes         string
	failures          string
	rejects           string
	shortCircuits     string
	timeouts          string
	fallbackSuccesses string
	fallbackFailures  string
	totalDuration     string
	runDuration       string
}

// CseCollectorConfig is a struct to keep monitoring information
type CseCollectorConfig struct {
	// CseMonitorAddr is the http address of the csemonitor server
	CseMonitorAddr string
	// Headers for csemonitor server
	Header http.Header
	// TickInterval spcifies the period that this collector will send metrics to the server.
	TimeInterval time.Duration
	// Config structure to configure a TLS client for sending Metric data
	TLSConfig *tls.Config
}

// NewCseCollector creates a new Collector Object
func NewCseCollector(name string) metricCollector.MetricCollector {
	return &CircuitCollector{
		attempts:          name + ".attempts",
		errors:            name + ".errors",
		successes:         name + ".successes",
		failures:          name + ".failures",
		rejects:           name + ".rejects",
		shortCircuits:     name + ".shortCircuits",
		timeouts:          name + ".timeouts",
		fallbackSuccesses: name + ".fallbackSuccesses",
		fallbackFailures:  name + ".fallbackFailures",
		totalDuration:     name + ".totalDuration",
		runDuration:       name + ".runDuration",
	}
}

func (c *CircuitCollector) incrementCounterMetric(prefix string) {
	count, ok := metrics.GetOrRegister(prefix, metrics.NewCounter).(metrics.Counter)
	if !ok {
		return
	}
	count.Inc(1)
}

func (c *CircuitCollector) updateTimerMetric(prefix string, dur time.Duration) {
	count, ok := metrics.GetOrRegister(prefix, metrics.NewTimer).(metrics.Timer)
	if !ok {
		return
	}
	count.Update(dur)
}

func (c *CircuitCollector) cleanMetric(prefix string) {
	count, ok := metrics.GetOrRegister(prefix, metrics.NewCounter).(metrics.Counter)
	if !ok {
		return
	}
	count.Clear()
}

// IncrementAttempts function increments the number of calls to this circuit.
// This registers as a counter
func (c *CircuitCollector) IncrementAttempts() {
	c.incrementCounterMetric(c.attempts)
}

// IncrementErrors function increments the number of unsuccessful attempts.
// Attempts minus Errors will equal successes.
// Errors are result from an attempt that is not a success.
// This registers as a counter
func (c *CircuitCollector) IncrementErrors() {
	c.incrementCounterMetric(c.errors)

}

// IncrementSuccesses function increments the number of requests that succeed.
// This registers as a counter
func (c *CircuitCollector) IncrementSuccesses() {
	c.incrementCounterMetric(c.successes)

}

// IncrementFailures function increments the number of requests that fail.
// This registers as a counter
func (c *CircuitCollector) IncrementFailures() {
	c.incrementCounterMetric(c.failures)
}

// IncrementRejects function increments the number of requests that are rejected.
// This registers as a counter
func (c *CircuitCollector) IncrementRejects() {
	c.incrementCounterMetric(c.rejects)
}

// IncrementShortCircuits function increments the number of requests that short circuited due to the circuit being open.
// This registers as a counter
func (c *CircuitCollector) IncrementShortCircuits() {
	c.incrementCounterMetric(c.shortCircuits)
}

// IncrementTimeouts function increments the number of timeouts that occurred in the circuit breaker.
// This registers as a counter
func (c *CircuitCollector) IncrementTimeouts() {
	c.incrementCounterMetric(c.timeouts)
}

// IncrementFallbackSuccesses function increments the number of successes that occurred during the execution of the fallback function.
// This registers as a counter
func (c *CircuitCollector) IncrementFallbackSuccesses() {
	c.incrementCounterMetric(c.fallbackSuccesses)
}

// IncrementFallbackFailures function increments the number of failures that occurred during the execution of the fallback function.
// This registers as a counter
func (c *CircuitCollector) IncrementFallbackFailures() {
	c.incrementCounterMetric(c.fallbackFailures)
}

// UpdateTotalDuration function updates the internal counter of how long we've run for.
// This registers as a timer
func (c *CircuitCollector) UpdateTotalDuration(timeSinceStart time.Duration) {
	c.updateTimerMetric(c.totalDuration, timeSinceStart)
}

// UpdateRunDuration function updates the internal counter of how long the last run took.
// This registers as a timer
func (c *CircuitCollector) UpdateRunDuration(runDuration time.Duration) {
	c.updateTimerMetric(c.runDuration, runDuration)
}

// Reset function is a noop operation in this collector.
func (c *CircuitCollector) Reset() {
	c.cleanMetric(c.attempts)
	c.cleanMetric(c.failures)
	c.cleanMetric(c.successes)
	c.cleanMetric(c.shortCircuits)
	c.cleanMetric(c.errors)
	c.cleanMetric(c.rejects)
	c.cleanMetric(c.timeouts)
	c.cleanMetric(c.fallbackSuccesses)
	c.cleanMetric(c.fallbackFailures)
}
