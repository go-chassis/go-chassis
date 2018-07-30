package metrics

import (
	"errors"

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

//TODO ReportMetricsToOpenTSDB use go-metrics reporter to send metrics to opentsdb
