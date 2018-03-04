package metrics

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/endpoint-discovery"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
)

func getTLSForClient(monitorURL string) (*tls.Config, error) {
	monitorServerURL, err := url.Parse(monitorURL)
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

func getMonitorEndpoint() (string, error) {
	monitorEndpoint := config.GlobalDefinition.Cse.Monitor.Client.ServerURI
	if monitorEndpoint == "" {
		monitorURL, err := endpoint.GetEndpointFromServiceCenter("default", "CseMonitoring", "latest")
		if err != nil {
			lager.Logger.Warn("empty monitor server endpoint, please provide the monitor server endpoint", err)
			return "", err
		}

		monitorEndpoint = monitorURL
	}

	return monitorEndpoint, nil
}
