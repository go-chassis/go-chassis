package metrics

import (
	"errors"
	"sync"
	"time"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"

	"github.com/ServiceComb/cse-collector"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix/metric_collector"
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

//metricCollector use go-metrics to send metrics to cse dashboard
func reportMetricsToCSEDashboard(r metrics.Registry) error {
	metricCollector.Registry.Register(metricsink.NewCseCollector)

	monitorServerURL, err := getMonitorEndpoint()
	if err != nil {
		lager.Logger.Warn("Get Monitoring URL failed, CSE monitoring function disabled", err)
		return nil
	}

	tlsConfig, tlsError := getTLSForClient(monitorServerURL)
	if tlsError != nil {
		lager.Logger.Errorf(tlsError, "Get %s.%s TLS config failed.", Name, common.Consumer)
		return tlsError
	}

	metricsink.InitializeCseCollector(&metricsink.CseCollectorConfig{
		CseMonitorAddr: monitorServerURL,
		Header:         getAuthHeaders(),
		TimeInterval:   time.Second * 2,
		TLSConfig:      tlsConfig,
	}, r, config.GlobalDefinition.AppID, config.SelfVersion, config.SelfServiceName)
	lager.Logger.Infof("Started sending metric Data to Monitor Server : %s", monitorServerURL)
	return nil
}

//TODO ReportMetricsToOpenTSDB use go-metrics reporter to send metrics to opentsdb

func init() {
	InstallReporter("Prometheus", ReportMetricsToPrometheus)
	InstallReporter("CSE Monitoring", reportMetricsToCSEDashboard)
}
