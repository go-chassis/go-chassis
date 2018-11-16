package metrics_test

import (
	"github.com/go-chassis/go-chassis/pkg/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	err := metrics.Init()
	assert.NoError(t, err)

}
func TestCounterAdd(t *testing.T) {
	err := metrics.CounterAdd("total", 1, map[string]string{
		"service": "s",
	})
	assert.Error(t, err)

	err = metrics.CreateCounter(metrics.CounterOpts{
		Name:   "total",
		Help:   "1",
		Labels: []string{"service"},
	})
	assert.NoError(t, err)
	err = metrics.CreateCounter(metrics.CounterOpts{
		Name:   "total",
		Help:   "1",
		Labels: []string{"service"},
	})
	assert.Error(t, err)

	err = metrics.CounterAdd("total", 1, map[string]string{
		"service": "s",
	})
	assert.NoError(t, err)
}

func TestGaugeSet(t *testing.T) {
	err := metrics.GaugeSet("cpu", 1, map[string]string{
		"service": "s",
	})
	assert.Error(t, err)

	err = metrics.CreateGauge(metrics.GaugeOpts{
		Name:   "cpu",
		Help:   "1",
		Labels: []string{"service"},
	})
	assert.NoError(t, err)
	err = metrics.CreateGauge(metrics.GaugeOpts{
		Name:   "cpu",
		Help:   "1",
		Labels: []string{"service"},
	})
	assert.Error(t, err)

	err = metrics.GaugeSet("cpu", 1, map[string]string{
		"service": "s",
	})
	assert.NoError(t, err)
}
func TestSummaryObserve(t *testing.T) {
	err := metrics.SummaryObserve("latency", 1, map[string]string{
		"service": "s",
	})
	assert.Error(t, err)

	err = metrics.CreateSummary(metrics.SummaryOpts{
		Name:   "latency",
		Help:   "1",
		Labels: []string{"service"},
	})
	assert.NoError(t, err)
	err = metrics.CreateSummary(metrics.SummaryOpts{
		Name:   "latency",
		Help:   "1",
		Labels: []string{"service"},
	})
	assert.Error(t, err)

	err = metrics.SummaryObserve("latency", 1, map[string]string{
		"service": "s",
	})
	assert.NoError(t, err)
}
