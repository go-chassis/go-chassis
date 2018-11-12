//Package metrics bootstrap metrics reporter, and supply 2 metrics registry
//native prometheus registry and rcrowley/go-metrics registry
//system registry is the place where go-chassis feed metrics data to
//you can get system registry and report them to varies monitoring system
package metrics

import (
	"sync"

	"github.com/go-chassis/go-chassis/core/lager"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-archaius"
	m "github.com/go-chassis/go-chassis/pkg/metrics"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix/metric_collector"
	"github.com/go-mesh/openlogging"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
)

// constants for header parameters
const (
	defaultName = "default_metric_registry"
	// Metrics is the constant string
	Metrics = "PrometheusMetrics"
)

var metricRegistries = make(map[string]metrics.Registry)
var l sync.RWMutex

//GetSystemRegistry return metrics registry which go chassis use
func GetSystemRegistry() metrics.Registry {
	return GetOrCreateRegistry(defaultName)
}

//GetOrCreateRegistry return a go-metrics registry which go chassis framework use to report metrics
func GetOrCreateRegistry(name string) metrics.Registry {
	r, ok := metricRegistries[name]
	if !ok {
		l.Lock()
		r = metrics.NewRegistry()
		metricRegistries[name] = r
		l.Unlock()
	}
	return r
}

// HTTPHandleFunc is a go-restful handler which can expose metrics in http server
func HTTPHandleFunc(req *restful.Request, rep *restful.Response) {
	promhttp.HandlerFor(m.GetSystemPrometheusRegistry(), promhttp.HandlerOpts{}).ServeHTTP(rep.ResponseWriter, req.Request)
}

//Init prepare the metrics registry and report metrics to other systems
func Init() error {
	metricRegistries[defaultName] = metrics.DefaultRegistry
	if archaius.GetBool("cse.metrics.enableCircuitMetrics", true) {
		metricCollector.Registry.Register(NewCseCollector)
	}

	for k, report := range reporterPlugins {
		openlogging.GetLogger().Info("report metrics to " + k)
		if err := report(GetSystemRegistry()); err != nil {
			lager.Logger.Warnf(err.Error(), err)
			return err
		}
	}
	return nil
}
