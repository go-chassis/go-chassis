package metrics

import (
	"time"

	"errors"
	"github.com/ServiceComb/cse-collector"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix/metric_collector"
	"github.com/emicklei/go-restful"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rcrowley/go-metrics"
	"sync"
)

var errMonitoringFail = errors.New("Con not report metrics to CSE monitoring service")

// constants for header parameters
const (
	//HeaderUserName is a variable of type string
	HeaderUserName   = "x-user-name"
	HeaderDomainName = "x-domain-name"
	ContentType      = "Content-Type"
	Name             = "monitor"
	defaultName      = "default_metric_registry"
	// Metrics is the constant string
	Metrics = "PrometheusMetrics"
)

var metricRegistries map[string]metrics.Registry = make(map[string]metrics.Registry)
var prometheusRegistry *prometheus.Registry = prometheus.NewRegistry()
var l sync.RWMutex

//GetSystemRegistry return metrics registry which go chassis use
func GetSystemRegistry() metrics.Registry {
	return GetOrCreateRegistry(defaultName)
}

//GetSystemPrometheusRegistry return prometheus registry which go chassis use
func GetSystemPrometheusRegistry() *prometheus.Registry {
	return prometheusRegistry
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

//ReportMetricsToPrometheus report metrics to prometheus registry, you can use GetSystemPrometheusRegistry to get prometheus registry. by default chassis will report system metrics to prometheus
func ReportMetricsToPrometheus(r metrics.Registry) {
	promConfig := GetPrometheusSinker(r, GetSystemPrometheusRegistry())
	if archaius.GetBool("cse.metrics.enableGoRuntimeMetrics", false) {
		promConfig.EnableRunTimeMetrics()
		lager.Logger.Info("Go Runtime Metrics is not enable")
	}
	go promConfig.UpdatePrometheusMetrics()
}

//metricCollector use go-metrics to send metrics to cse dashboard
func reportMetricsToCSEDashboard(r metrics.Registry) error {
	metricCollector.Registry.Register(metricsink.NewCseCollector)

	monitorServerURL := config.GlobalDefinition.Cse.Monitor.Client.ServerURI
	if monitorServerURL == "" {
		lager.Logger.Warn("empty monitor server endpoint, please provide the monitor server endpoint", nil)
		return nil
	}

	tlsConfig, tlsError := getTLSForClient()
	if tlsError != nil {
		lager.Logger.Errorf(tlsError, "Get %s.%s TLS config failed.", Name, common.Consumer)
		return tlsError
	}

	metricsink.InitializeCseCollector(&metricsink.CseCollectorConfig{
		CseMonitorAddr: monitorServerURL,
		Header:         getAuthHeaders(),
		TimeInterval:   time.Second * 2,
		TLSConfig:      tlsConfig,
	}, r)
	lager.Logger.Infof("Started sending metric Data to Monitor Server : %s", monitorServerURL)
	return nil
}

//TODO ReportMetricsToOpenTSDB use go-metrics reporter to send metrics to opentsdb

// ReportMetricsToOpenTSDB use go-metrics reporter to send metrics to opentsdb
func ReportMetricsToOpenTSDB(r metrics.Registry) {}

// MetricsHandleFunc is a restful handler which can expose metrics in http server
func MetricsHandleFunc(req *restful.Request, rep *restful.Response) {
	reg := DefaultPrometheusSinker.PromRegistry.(*prometheus.Registry)
	promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(rep.ResponseWriter, req.Request)
}

//Init prepare the metrics functions
func Init() error {
	metricRegistries[defaultName] = metrics.DefaultRegistry
	if err := reportMetricsToCSEDashboard(GetSystemRegistry()); err != nil {
		lager.Logger.Error(errMonitoringFail.Error(), err)
		return errMonitoringFail
	}
	ReportMetricsToPrometheus(GetSystemRegistry())
	return nil
}
