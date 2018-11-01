package metrics_test

import (
	chassisMetrics "github.com/go-chassis/go-chassis/metrics"
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var cltr = chassisMetrics.NewCseCollector("server.server.server.server.server")

func TestIncrementAttempts(t *testing.T) {
	cltr.IncrementAttempts()
	metric := metrics.Get("server.server.server.server.server.attempts").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementSuccesses(t *testing.T) {
	cltr.IncrementSuccesses()
	metric := metrics.Get("server.server.server.server.server.successes").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementFailures(t *testing.T) {
	cltr.IncrementFailures()
	metric := metrics.Get("server.server.server.server.server.failures").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementRejects(t *testing.T) {
	cltr.IncrementRejects()
	metric := metrics.Get("server.server.server.server.server.rejects").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementShortCircuits(t *testing.T) {
	cltr.IncrementShortCircuits()
	metric := metrics.Get("server.server.server.server.server.shortCircuits").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementTimeouts(t *testing.T) {
	cltr.IncrementTimeouts()
	metric := metrics.Get("server.server.server.server.server.timeouts").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementFallbackSuccesses(t *testing.T) {
	cltr.IncrementFallbackSuccesses()
	metric := metrics.Get("server.server.server.server.server.fallbackSuccesses").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementFallbackFailures(t *testing.T) {
	cltr.IncrementFallbackFailures()
	metric := metrics.Get("server.server.server.server.server.fallbackFailures").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestUpdateTotalDuration(t *testing.T) {
	cltr.UpdateTotalDuration(time.Second)
	metric := metrics.Get("server.server.server.server.server.totalDuration").(metrics.Timer)
	assert.Equal(t, metric.Count(), int64(1))
}
func TestUpdateRunDuration(t *testing.T) {
	cltr.UpdateRunDuration(time.Second)
	metric := metrics.Get("server.server.server.server.server.runDuration").(metrics.Timer)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementErrors(t *testing.T) {
	cltr.IncrementErrors()
	metric := metrics.Get("server.server.server.server.server.errors").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}
