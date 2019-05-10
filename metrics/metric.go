//Package metrics bootstrap metrics reporter, and supply a metrics registries
//powered by rcrowley/go-metrics
//there is a default registry "default_metric_registry"
//which saves go chassis runtime metrics
//you can report them to varies monitoring system by install Reporter plugin
package metrics

import (
	"sync"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix/metric_collector"
	"github.com/go-mesh/openlogging"
	"github.com/rcrowley/go-metrics"
)

// constants for header parameters
const (
	defaultName = "default_metric_registry"
)

var metricRegistries = make(map[string]metrics.Registry)
var l sync.RWMutex

//GetSystemRegistry return metrics registry which go chassis use
//it saves go chassis runtime metrics, aka. system registry
func GetSystemRegistry() metrics.Registry {
	return GetOrCreateRegistry(defaultName)
}

//GetOrCreateRegistry return a go-metrics registry
//you can use it freely
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

//Init prepare the metrics registry and
//report system registry to other systems
func Init() error {
	metricRegistries[defaultName] = metrics.DefaultRegistry
	if archaius.GetBool("cse.metrics.enableCircuitMetrics", true) {
		metricCollector.Registry.Register(NewCseCollector)
	}

	for k, report := range reporterPlugins {
		openlogging.GetLogger().Info("report metrics to " + k)
		if err := report(GetSystemRegistry()); err != nil {
			openlogging.Warn("can not report: " + err.Error())
			return err
		}
	}
	return nil
}
