package configserver

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-chassis/go-archaius/source/remote"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/endpoint"
	chassisTLS "github.com/go-chassis/go-chassis/v2/core/tls"
	"net/url"
	"strings"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/openlog"
)

const (
	//configServerName is a variable of type string of config server
	configServerName = "configServer"
)

// ErrRefreshMode means config is mis used
var (
	ErrRefreshMode      = errors.New("refreshMode must be 0 or 1")
	ErrRegistryDisabled = errors.New("discovery is disabled")
)

// Init initialize config server
func Init() error {
	configServerURL, err := GetConfigServerEndpoint()
	if err != nil {
		openlog.Warn("can not get config server endpoint: " + err.Error())
		return err
	}

	var enableSSL bool
	tlsConfig, tlsError := getTLSForClient(configServerURL)
	if tlsError != nil {
		openlog.Error(fmt.Sprintf("Get %s.%s TLS config failed, err:[%s]",
			configServerName, common.Consumer, tlsError.Error()))
		return tlsError
	}

	/*This condition added because member discovery can have multiple ip's with IsHTTPS
	having both true and false value.*/
	if tlsConfig != nil {
		enableSSL = true
	}

	interval := config.GetConfigServerConf().RefreshInterval
	if interval == 0 {
		interval = 30
	}

	err = initConfigServer(configServerURL, enableSSL, tlsConfig, interval)
	if err != nil {
		openlog.Error("failed to init config server: " + err.Error())
		return err
	}

	openlog.Warn("config server init success")
	return nil
}

// GetConfigServerEndpoint will read local config server uri first, if there is not,
// it will try to discover config server from registry
func GetConfigServerEndpoint() (string, error) {
	configServerURL := config.GetConfigServerConf().ServerURI
	if configServerURL == "" {
		if registry.DefaultServiceDiscoveryService != nil {
			openlog.Debug("find config server in registry")
			ccURL, err := endpoint.GetEndpoint("default", "CseConfigCenter", "latest")
			if err != nil {
				openlog.Warn("failed to find config server endpoints, err: " + err.Error())
				return "", err
			}
			configServerURL = ccURL
		} else {
			return "", ErrRegistryDisabled
		}
	}

	return configServerURL, nil
}

func getTLSForClient(configServerURL string) (*tls.Config, error) {
	if !strings.Contains(configServerURL, "://") {
		return nil, nil
	}
	ccURL, err := url.Parse(configServerURL)
	if err != nil {
		openlog.Error("Error occurred while parsing config Server Uri" + err.Error())
		return nil, err
	}
	if ccURL.Scheme == common.HTTP {
		return nil, nil
	}

	sslTag := configServerName + "." + common.Consumer
	tlsConfig, sslConfig, err := chassisTLS.GetTLSConfigByService(configServerName, "", common.Consumer)
	if err != nil {
		if chassisTLS.IsSSLConfigNotExist(err) {
			return nil, fmt.Errorf("%s TLS mode, but no ssl config", sslTag)
		}
		return nil, err
	}
	openlog.Warn(fmt.Sprintf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
		sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin))

	return tlsConfig, nil
}

func initConfigServer(endpoint string, enableSSL bool, tlsConfig *tls.Config, interval int) error {

	refreshMode := archaius.GetInt("servicecomb.config.client.refreshMode", common.DefaultRefreshMode)
	if refreshMode != remote.ModeWatch && refreshMode != remote.ModeInterval {
		openlog.Error(ErrRefreshMode.Error())
		return ErrRefreshMode
	}

	remoteSourceType := archaius.GetString("servicecomb.config.client.type", archaius.KieSource)

	var ri = &archaius.RemoteInfo{
		DefaultDimension: map[string]string{
			remote.LabelApp:         runtime.App,
			remote.LabelService:     runtime.ServiceName,
			remote.LabelVersion:     runtime.Version,
			remote.LabelEnvironment: runtime.Environment,
		},
		URL:             endpoint,
		EnableSSL:       enableSSL,
		TLSConfig:       tlsConfig,
		RefreshMode:     refreshMode,
		RefreshInterval: interval,
		AutoDiscovery:   config.GetConfigServerConf().AutoDiscovery,
		APIVersion:      config.GetConfigServerConf().APIVersion.Version,
		RefreshPort:     config.GetConfigServerConf().RefreshPort,
	}

	err := archaius.EnableRemoteSource(remoteSourceType, ri)

	if err != nil {
		return err
	}

	if err := refreshGlobalConfig(); err != nil {
		openlog.Error("failed to refresh global config for lb and cb:" + err.Error())
		return err
	}
	return nil
}

func refreshGlobalConfig() error {
	err := config.ReadHystrixFromArchaius()
	if err != nil {
		return err
	}
	return config.ReadLBFromArchaius()
}
