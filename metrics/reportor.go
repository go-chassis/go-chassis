package metrics

import (
	"errors"
	"sync"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/rcrowley/go-metrics"
)

//Reporter receive a go-metrics registry and sink it to monitoring system
type Reporter func(metrics.Registry) error

//ErrDuplicated means you can not install reporter with same name
var ErrDuplicated = errors.New("duplicated reporter")
var reporterPlugins = make(map[string]Reporter)

//InstallReporter install reporter implementation
func InstallReporter(name string, reporter Reporter) error {
	_, ok := reporterPlugins[name]
	if ok {
		return ErrDuplicated
	}
	reporterPlugins[name] = reporter
	return nil
}

var onceEnable sync.Once

//ReportMetricsToPrometheus report metrics to prometheus registry, you can use GetSystemPrometheusRegistry to get prometheus registry. by default chassis will report system metrics to prometheus
func ReportMetricsToPrometheus(r metrics.Registry) error {
	promConfig := GetPrometheusSinker(r, GetSystemPrometheusRegistry())
	if archaius.GetBool("cse.metrics.enableGoRuntimeMetrics", true) {
		onceEnable.Do(func() {
			promConfig.EnableRunTimeMetrics()
			lager.Logger.Info("Go Runtime Metrics is enabled")
		})

	}
	go promConfig.UpdatePrometheusMetrics()
	return nil
}

//TODO ReportMetricsToOpenTSDB use go-metrics reporter to send metrics to opentsdb

func init() {
	InstallReporter("Prometheus", ReportMetricsToPrometheus)
}
