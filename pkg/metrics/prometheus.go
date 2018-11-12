package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"sync"
	"time"
)

//PrometheusExporter is a prom exporter for go chassis
type PrometheusExporter struct {
	FlushInterval time.Duration
	lc            sync.RWMutex
	lg            sync.RWMutex
	ls            sync.RWMutex
	counters      map[string]*prometheus.CounterVec
	gauges        map[string]*prometheus.GaugeVec
	summaries     map[string]*prometheus.SummaryVec
}

//NewPrometheusExporter create a prometheus exporter
func NewPrometheusExporter(options Options) Registry {
	return &PrometheusExporter{
		FlushInterval: options.FlushInterval,
		lc:            sync.RWMutex{},
		lg:            sync.RWMutex{},
		ls:            sync.RWMutex{},
		summaries:     make(map[string]*prometheus.SummaryVec),
		counters:      make(map[string]*prometheus.CounterVec),
		gauges:        make(map[string]*prometheus.GaugeVec),
	}
}

//CreateGauge create collector
func (c *PrometheusExporter) CreateGauge(opts GaugeOpts) error {
	c.lg.RLock()
	_, ok := c.gauges[opts.Name]
	c.lg.RUnlock()
	if ok {
		return fmt.Errorf("metric [%s] is duplicated", opts.Name)
	}
	c.lg.Lock()
	defer c.lg.Unlock()
	gVec := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: opts.Name,
		Help: opts.Help,
	}, opts.Labels)
	c.gauges[opts.Name] = gVec
	return nil
}

//GaugeSet set value
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

//CreateCounter create collector
func (c *PrometheusExporter) CreateCounter(opts CounterOpts) error {
	c.lc.RLock()
	_, ok := c.counters[opts.Name]
	c.lc.RUnlock()
	if ok {
		return fmt.Errorf("metric [%s] is duplicated", opts.Name)
	}
	c.lc.Lock()
	defer c.lc.Unlock()
	v := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: opts.Name,
		Help: opts.Help,
	}, opts.Labels)
	c.counters[opts.Name] = v
	return nil
}

//CounterAdd increase value
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

//CreateSummary create collector
func (c *PrometheusExporter) CreateSummary(opts SummaryOpts) error {
	c.ls.RLock()
	_, ok := c.summaries[opts.Name]
	c.ls.RUnlock()
	if ok {
		return fmt.Errorf("metric [%s] is duplicated", opts.Name)
	}
	c.ls.Lock()
	defer c.ls.Unlock()
	v := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       opts.Name,
		Help:       opts.Help,
		Objectives: opts.Objectives,
	}, opts.Labels)
	c.summaries[opts.Name] = v
	return nil
}

//SummaryObserve set value
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
func init() {
	registries["prometheus"] = NewPrometheusExporter
}
