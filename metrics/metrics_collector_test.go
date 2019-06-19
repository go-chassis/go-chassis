package metrics_test

import (
	chassisMetrics "github.com/go-chassis/go-chassis/metrics"
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var cltr = chassisMetrics.NewCseCollector("server")

func TestIncrementAttempts(t *testing.T) {
	cltr.IncrementAttempts()
	metric := metrics.Get("server.attempts").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementSuccesses(t *testing.T) {
	cltr.IncrementSuccesses()
	metric := metrics.Get("server.successes").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementFailures(t *testing.T) {
	cltr.IncrementFailures()
	metric := metrics.Get("server.failures").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementRejects(t *testing.T) {
	cltr.IncrementRejects()
	metric := metrics.Get("server.rejects").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementShortCircuits(t *testing.T) {
	cltr.IncrementShortCircuits()
	metric := metrics.Get("server.shortCircuits").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementTimeouts(t *testing.T) {
	cltr.IncrementTimeouts()
	metric := metrics.Get("server.timeouts").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementFallbackSuccesses(t *testing.T) {
	cltr.IncrementFallbackSuccesses()
	metric := metrics.Get("server.fallbackSuccesses").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementFallbackFailures(t *testing.T) {
	cltr.IncrementFallbackFailures()
	metric := metrics.Get("server.fallbackFailures").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestUpdateTotalDuration(t *testing.T) {
	cltr.UpdateTotalDuration(time.Second)
	metric := metrics.Get("server.totalDuration").(metrics.Timer)
	assert.Equal(t, metric.Count(), int64(1))
}
func TestUpdateRunDuration(t *testing.T) {
	cltr.UpdateRunDuration(time.Second)
	metric := metrics.Get("server.runDuration").(metrics.Timer)
	assert.Equal(t, metric.Count(), int64(1))
}

func TestIncrementErrors(t *testing.T) {
	cltr.IncrementErrors()
	metric := metrics.Get("server.errors").(metrics.Counter)
	assert.Equal(t, metric.Count(), int64(1))

	t.Run("reset", func(t *testing.T) {
		cltr.Reset()
		metric := metrics.Get("server.errors").(metrics.Counter)
		assert.Equal(t, int64(0), metric.Count())
	})
}
