package metrics

import (
	"crypto/tls"
	"fmt"
	"github.com/ServiceComb/cse-collector"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix/metric_collector"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"net/url"
	"time"
)

func registerCircuitBreakerCollector(r metrics.Registry) error {
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

func getTLSForClient() (*tls.Config, error) {
	base := config.GlobalDefinition.Cse.Monitor.Client.ServerURI
	monitorServerURL, err := url.Parse(base)
	if err != nil {
		lager.Logger.Error("Error occurred while parsing Monitor Server Uri", err)
		return nil, err
	}
	scheme := monitorServerURL.Scheme
	if scheme != "https" {
		return nil, nil
	}

	sslTag := Name + "." + common.Consumer
	tlsConfig, sslConfig, err := chassisTLS.GetTLSConfigByService(Name, "", common.Consumer)
	if err != nil {
		if chassisTLS.IsSSLConfigNotExist(err) {
			return nil, fmt.Errorf("%s TLS mode, but no ssl config", sslTag)
		}
		return nil, err
	}
	lager.Logger.Warnf(nil, "%s TLS mode, verify peer: %t, cipher plugin: %s",
		sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin)

	return tlsConfig, nil
}
func getAuthHeaders() http.Header {
	userName := config.GlobalDefinition.Cse.Monitor.Client.UserName
	if userName == "" {
		userName = common.DefaultUserName
	}
	domainName := config.GlobalDefinition.Cse.Monitor.Client.DomainName
	if domainName == "" {
		domainName = common.DefaultDomainName
	}

	headers := make(http.Header)
	headers.Set(HeaderUserName, userName)
	headers.Set(HeaderDomainName, domainName)
	headers.Set(ContentType, "application/json")

	return headers
}
