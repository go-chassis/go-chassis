package metrics

import (
	"fmt"
	"github.com/go-chassis/go-archaius"
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

var registries = make(map[string]NewRegistry)
var prometheusRegistry = prometheus.NewRegistry()

// NewRegistry create a registry
type NewRegistry func(opts Options) Registry

// Registry holds all of metrics collectors
// name is a unique ID for different type of metrics
type Registry interface {
	CreateGauge(opts GaugeOpts) error
	CreateCounter(opts CounterOpts) error
	CreateSummary(opts SummaryOpts) error
	CreateHistogram(opts HistogramOpts) error

	GaugeSet(name string, val float64, labels map[string]string) error
	GaugeAdd(name string, val float64, labels map[string]string) error
	CounterAdd(name string, val float64, labels map[string]string) error
	SummaryObserve(name string, val float64, Labels map[string]string) error
	HistogramObserve(name string, val float64, labels map[string]string) error

	GaugeValue(name string, labels map[string]string) float64
	CounterValue(name string, labels map[string]string) float64

	Reset(name string) error
}

var defaultRegistry Registry

// CreateGauge init a new gauge type
func CreateGauge(opts GaugeOpts) error {
	return defaultRegistry.CreateGauge(opts)
}

// CreateCounter init a new counter type
func CreateCounter(opts CounterOpts) error {
	return defaultRegistry.CreateCounter(opts)
}

// CreateSummary init a new summary type
func CreateSummary(opts SummaryOpts) error {
	return defaultRegistry.CreateSummary(opts)
}

// CreateHistogram init a new summary type
func CreateHistogram(opts HistogramOpts) error {
	return defaultRegistry.CreateHistogram(opts)
}

// GaugeSet set a new value to a collector
func GaugeSet(name string, val float64, labels map[string]string) error {
	return defaultRegistry.GaugeSet(name, val, labels)
}

// GaugeAdd set a new value to a collector
func GaugeAdd(name string, val float64, labels map[string]string) error {
	return defaultRegistry.GaugeAdd(name, val, labels)
}

// CounterAdd increase value of a collector
func CounterAdd(name string, val float64, labels map[string]string) error {
	return defaultRegistry.CounterAdd(name, val, labels)
}

// SummaryObserve gives a value to summary collector
func SummaryObserve(name string, val float64, labels map[string]string) error {
	return defaultRegistry.SummaryObserve(name, val, labels)
}

// HistogramObserve gives a value to histogram collector
func HistogramObserve(name string, val float64, labels map[string]string) error {
	return defaultRegistry.HistogramObserve(name, val, labels)
}

// Reset clear collector metrics
func Reset(name string) error {
	return defaultRegistry.Reset(name)
}

func GaugeValue(name string, labels map[string]string) float64 {
	return defaultRegistry.GaugeValue(name, labels)
}

func CounterValue(name string, labels map[string]string) float64 {
	return defaultRegistry.CounterValue(name, labels)
}

// CounterOpts is options to create a counter options
type CounterOpts struct {
	// Key is the key set joining with '_', Name will be ignored when Key is not empty
	Key    string
	Name   string
	Help   string
	Labels []string
}

// GaugeOpts is options to create a gauge collector
type GaugeOpts struct {
	Key    string
	Name   string
	Help   string
	Labels []string
}

// SummaryOpts is options to create summary collector
type SummaryOpts struct {
	Key        string
	Name       string
	Help       string
	Labels     []string
	Objectives map[float64]float64
}

// HistogramOpts is options to create histogram collector
type HistogramOpts struct {
	Key     string
	Name    string
	Help    string
	Labels  []string
	Buckets []float64
}

// Options control config
type Options struct {
	FlushInterval          time.Duration
	EnableGoRuntimeMetrics bool
}

// InstallPlugin install metrics registry
func InstallPlugin(name string, f NewRegistry) {
	registries[name] = f
}

// Init load the metrics plugin and initialize it
func Init() error {
	//TODO name should be configurable
	name := "prometheus"
	f, ok := registries[name]
	if !ok {
		return fmt.Errorf("can not init metrics registry [%s]", name)
	}
	defaultRegistry = f(Options{
		FlushInterval:          10 * time.Second,
		EnableGoRuntimeMetrics: archaius.GetBool("servicecomb.metrics.enableGoRuntimeMetrics", true),
	})
	return nil
}

// GetSystemPrometheusRegistry return prometheus registry which go chassis use
func GetSystemPrometheusRegistry() *prometheus.Registry {
	return prometheusRegistry
}
