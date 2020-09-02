package metrics_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/pkg/metrics"
	"github.com/go-chassis/openlog"
	"github.com/prometheus/common/expfmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	t.Run("install a plugin", func(t *testing.T) {
		metrics.InstallPlugin("test", metrics.NewPrometheusExporter)
	})
	err := archaius.Init(archaius.WithENVSource())
	assert.NoError(t, err)
	err = metrics.Init()
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
func TestCreateHistogram(t *testing.T) {
	err := metrics.HistogramObserve("hlatency", 1, map[string]string{
		"service": "s",
	})
	assert.Error(t, err)

	err = metrics.CreateHistogram(metrics.HistogramOpts{
		Name:   "hlatency",
		Help:   "1",
		Labels: []string{"service"},
	})
	assert.NoError(t, err)
	err = metrics.CreateHistogram(metrics.HistogramOpts{
		Name:   "hlatency",
		Help:   "1",
		Labels: []string{"service"},
	})
	assert.Error(t, err)

	err = metrics.HistogramObserve("hlatency", 1, map[string]string{
		"service": "s",
	})
	assert.NoError(t, err)
}

type writer struct {
}

func (w *writer) Write(b []byte) (n int, err error) {
	openlog.Error(string(b))
	return len(b), nil
}
func TestNewPrometheusExporter(t *testing.T) {
	mfs, err := metrics.GetSystemPrometheusRegistry().Gather()
	assert.NoError(t, err)
	w := &writer{}
	enc := expfmt.NewEncoder(w, "text/plain; version=0.0.4; charset=utf-8")
	for _, mf := range mfs {
		err := enc.Encode(mf)
		assert.NoError(t, err)
	}
}
