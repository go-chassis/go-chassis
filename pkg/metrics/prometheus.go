package metrics

import (
	"fmt"
	"github.com/go-chassis/openlog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"strings"
	"sync"
	"time"

	dto "github.com/prometheus/client_model/go"
)

var onceEnable sync.Once

// PrometheusExporter is a prom exporter for go chassis
type PrometheusExporter struct {
	FlushInterval time.Duration
	lc            sync.RWMutex
	lg            sync.RWMutex
	ls            sync.RWMutex
	lh            sync.RWMutex
	counters      map[string]*prometheus.CounterVec
	gauges        map[string]*prometheus.GaugeVec
	summaries     map[string]*prometheus.SummaryVec
	histograms    map[string]*prometheus.HistogramVec
}

// NewPrometheusExporter create a prometheus exporter
func NewPrometheusExporter(options Options) Registry {
	if options.EnableGoRuntimeMetrics {
		onceEnable.Do(func() {
			EnableRunTimeMetrics()
			openlog.Info("go runtime metrics is exported")
		})

	}
	return &PrometheusExporter{
		FlushInterval: options.FlushInterval,
		lc:            sync.RWMutex{},
		lg:            sync.RWMutex{},
		ls:            sync.RWMutex{},
		lh:            sync.RWMutex{},
		summaries:     make(map[string]*prometheus.SummaryVec),
		counters:      make(map[string]*prometheus.CounterVec),
		gauges:        make(map[string]*prometheus.GaugeVec),
		histograms:    make(map[string]*prometheus.HistogramVec),
	}
}

// EnableRunTimeMetrics enable runtime metrics
func EnableRunTimeMetrics() {
	GetSystemPrometheusRegistry().MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	GetSystemPrometheusRegistry().MustRegister(collectors.NewGoCollector())
}

// CreateGauge create collector
func (c *PrometheusExporter) CreateGauge(opts GaugeOpts) error {
	key := opts.Key
	ns, sub, name := Split(key)
	if len(name) == 0 {
		key, name = opts.Name, opts.Name
	}
	c.lg.RLock()
	_, ok := c.gauges[key]
	c.lg.RUnlock()
	if ok {
		return fmt.Errorf("metric [%s] is duplicated", key)
	}
	c.lg.Lock()
	defer c.lg.Unlock()
	gVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      name,
		Help:      opts.Help,
	}, opts.Labels)
	c.gauges[key] = gVec
	GetSystemPrometheusRegistry().MustRegister(gVec)
	return nil
}

// GaugeSet set value
func (c *PrometheusExporter) GaugeSet(name string, val float64, labels map[string]string) error {
	c.lg.RLock()
	gVec, ok := c.gauges[name]
	c.lg.RUnlock()
	if !ok {
		return fmt.Errorf("metrics do not exists, create it first")
	}
	gVec.With(labels).Set(val)
	return nil
}

// GaugeAdd add value, can be negative
func (c *PrometheusExporter) GaugeAdd(name string, val float64, labels map[string]string) error {
	c.lg.RLock()
	gVec, ok := c.gauges[name]
	c.lg.RUnlock()
	if !ok {
		return fmt.Errorf("metrics do not exists, create it first")
	}
	gVec.With(labels).Add(val)
	return nil
}

func (c *PrometheusExporter) GaugeValue(name string, labels map[string]string) float64 {
	return getValue(name, labels, func(m *dto.Metric) float64 {
		return m.GetGauge().GetValue()
	})
}

func (c *PrometheusExporter) CounterValue(name string, labels map[string]string) float64 {
	return getValue(name, labels, func(m *dto.Metric) float64 {
		return m.GetCounter().GetValue()
	})
}

// CreateCounter create collector
func (c *PrometheusExporter) CreateCounter(opts CounterOpts) error {
	key := opts.Key
	ns, sub, name := Split(key)
	if len(name) == 0 {
		key, name = opts.Name, opts.Name
	}
	c.lc.RLock()
	_, ok := c.counters[key]
	c.lc.RUnlock()
	if ok {
		return fmt.Errorf("metric [%s] is duplicated", key)
	}
	c.lc.Lock()
	defer c.lc.Unlock()
	v := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      name,
		Help:      opts.Help,
	}, opts.Labels)
	c.counters[key] = v
	GetSystemPrometheusRegistry().MustRegister(v)
	return nil
}

// CounterAdd increase value
func (c *PrometheusExporter) CounterAdd(name string, val float64, labels map[string]string) error {
	c.lc.RLock()
	v, ok := c.counters[name]
	c.lc.RUnlock()
	if !ok {
		return fmt.Errorf("metrics do not exists, create it first")
	}
	v.With(labels).Add(val)
	return nil
}

// CreateSummary create collector
func (c *PrometheusExporter) CreateSummary(opts SummaryOpts) error {
	key := opts.Key
	ns, sub, name := Split(key)
	if len(name) == 0 {
		key, name = opts.Name, opts.Name
	}
	c.ls.RLock()
	_, ok := c.summaries[key]
	c.ls.RUnlock()
	if ok {
		return fmt.Errorf("metric [%s] is duplicated", key)
	}
	c.ls.Lock()
	defer c.ls.Unlock()
	v := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  ns,
		Subsystem:  sub,
		Name:       name,
		Help:       opts.Help,
		Objectives: opts.Objectives,
	}, opts.Labels)
	c.summaries[key] = v
	GetSystemPrometheusRegistry().MustRegister(v)
	return nil
}

// SummaryObserve set value
func (c *PrometheusExporter) SummaryObserve(name string, val float64, labels map[string]string) error {
	c.ls.RLock()
	v, ok := c.summaries[name]
	c.ls.RUnlock()
	if !ok {
		return fmt.Errorf("metrics do not exists, create it first")
	}
	v.With(labels).Observe(val)
	return nil
}

// CreateHistogram create collector
func (c *PrometheusExporter) CreateHistogram(opts HistogramOpts) error {
	key := opts.Key
	ns, sub, name := Split(key)
	if len(name) == 0 {
		key, name = opts.Name, opts.Name
	}
	c.lh.RLock()
	_, ok := c.histograms[key]
	c.lh.RUnlock()
	if ok {
		return fmt.Errorf("metric [%s] is duplicated", key)
	}
	c.lh.Lock()
	defer c.lh.Unlock()
	v := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: ns,
		Subsystem: sub,
		Name:      name,
		Help:      opts.Help,
		Buckets:   opts.Buckets,
	}, opts.Labels)
	c.histograms[key] = v
	GetSystemPrometheusRegistry().MustRegister(v)
	return nil
}

// HistogramObserve set value
func (c *PrometheusExporter) HistogramObserve(name string, val float64, labels map[string]string) error {
	c.lh.RLock()
	v, ok := c.histograms[name]
	c.lh.RUnlock()
	if !ok {
		return fmt.Errorf("metrics do not exists, create it first")
	}
	v.With(labels).Observe(val)
	return nil
}

// Reset reset a collector metrics
func (c *PrometheusExporter) Reset(name string) error {
	c.lc.RLock()
	ct, ok := c.counters[name]
	c.lc.RUnlock()
	if ok {
		ct.Reset()
		return nil
	}
	c.lg.RLock()
	gg, ok := c.gauges[name]
	c.lg.RUnlock()
	if ok {
		gg.Reset()
		return nil
	}
	c.ls.RLock()
	sm, ok := c.summaries[name]
	c.ls.RUnlock()
	if ok {
		sm.Reset()
		return nil
	}
	c.lh.RLock()
	hg, ok := c.histograms[name]
	c.lh.RUnlock()
	if ok {
		hg.Reset()
		return nil
	}
	return nil
}

func Split(key string) (string, string, string) {
	arr := strings.Split(key, "_")
	var ns, sub string
	i, l := 0, len(arr)
	if l > 1 {
		i, ns = 1, arr[0]
	}
	if l > 2 {
		i, sub = 2, arr[1]
	}
	name := strings.Join(arr[i:], "_")
	return ns, sub, name
}

func init() {
	registries["prometheus"] = NewPrometheusExporter
}

func getValue(name string, labels map[string]string, getV func(m *dto.Metric) float64) float64 {
	f := family(name)
	if f == nil {
		return 0
	}
	matchAll := len(labels) == 0
	var sum float64
	for _, m := range f.Metric {
		if !matchAll && !matchLabels(m, labels) {
			continue
		}
		sum += getV(m)
	}
	return sum
}

func family(name string) *dto.MetricFamily {
	families, err := GetSystemPrometheusRegistry().Gather()
	if err != nil {
		return nil
	}
	for _, f := range families {
		if f.GetName() == name {
			return f
		}
	}
	return nil
}

func matchLabels(m *dto.Metric, labels map[string]string) bool {
	count := 0
	for _, label := range m.GetLabel() {
		v, ok := labels[label.GetName()]
		if ok && v != label.GetValue() {
			return false
		}
		if ok {
			count++
		}
	}
	return count == len(labels)
}
