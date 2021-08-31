package tls

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"k8s.io/apimachinery/pkg/util/sets"
)

var errSSLConfigNotExist = errors.New("No SSL config")
var useDefaultSslTag = sets.NewString(
	"registry.Consumer.",
	"configServer.Consumer.",
	"monitor.Consumer.",
	"serviceDiscovery.Consumer.",
	"registrator.Consumer.",
	"contractDiscovery.Consumer.",
	"router.Consumer",
)

func hasDefaultSslTag(tag string) bool {
	if len(tag) == 0 {
		return false
	}

	if useDefaultSslTag.Has(tag) {
		return true
	}
	return false
}

func getDefaultSslConfigMap() map[string]string {
	cipherSuits := []string{}
	for k := range TLSCipherSuiteMap {
		cipherSuits = append(cipherSuits, k)
	}

	cipherSuitesKey := strings.Join(cipherSuits, ",")
	defaultSslConfigMap := map[string]string{
		common.SslCipherPluginKey: "default",
		common.SslVerifyPeerKey:   common.FALSE,
		common.SslCipherSuitsKey:  cipherSuitesKey,
		common.SslProtocolKey:     "TLSv1.2",
		common.SslCaFileKey:       "",
		common.SslCertFileKey:     "",
		common.SslKeyFileKey:      "",
		common.SslCertPwdFilePath: "",
		common.SslServerNameKey:   "",
	}
	return defaultSslConfigMap
}

func getSSLConfigMap(tag, protocol, svcType, svcName string) map[string]string {
	sslConfigMap := config.GlobalDefinition.Ssl
	defaultSslConfigMap := getDefaultSslConfigMap()
	result := make(map[string]string)

	sslSet := false
	if tag != "" {
		tag = tag + `.`
	}

	for k, v := range defaultSslConfigMap {
		// 使用默认配置
		result[k] = v
		// 若配置了全局配置项，则覆盖默认配置
		if r, exist := sslConfigMap[k]; exist && r != "" {
			result[k] = r
			sslSet = true
		}
		// consumer如果配置了通配 生效级别为全局配置之上 指定配置之下 且不为自己代理的服务设置证书配置
		if useGeneralTLSConfig(svcType, svcName) {
			consumerKey := protocol + "." + common.Consumer + "." + k
			if c, exist := sslConfigMap[consumerKey]; exist && c != "" {
				result[k] = c
				sslSet = true
			}
		}
		// 若配置了指定交互方的配置项，则覆盖全局配置
		keyWithTag := tag + k
		if v, exist := sslConfigMap[keyWithTag]; exist && v != "" {
			result[k] = v
			sslSet = true
		}
	}
	// 未设置ssl 且不提供内部默认ss配置 返回空字典
	if !sslSet && !hasDefaultSslTag(tag) {
		return make(map[string]string)
	}

	return result
}

// use general TLSConfig
func useGeneralTLSConfig(svcType, svcName string) bool {
	return common.Consumer == svcType && svcName != config.GlobalDefinition.ServiceComb.ServiceDescription.Name
}

func parseSSLConfig(sslConfigMap map[string]string) (*SSLConfig, error) {
	sslConfig := &SSLConfig{}
	var err error

	sslConfig.CipherPlugin = sslConfigMap[common.SslCipherPluginKey]

	sslConfig.VerifyPeer, err = strconv.ParseBool(sslConfigMap[common.SslVerifyPeerKey])
	if err != nil {
		return nil, err
	}

	sslConfig.CipherSuites, err = ParseSSLCipherSuites(sslConfigMap[common.SslCipherSuitsKey])
	if err != nil {
		return nil, err
	}
	if len(sslConfig.CipherSuites) == 0 {
		return nil, fmt.Errorf("no valid cipher")
	}

	sslConfig.MinVersion, err = ParseSSLProtocol(sslConfigMap[common.SslProtocolKey])
	if err != nil {
		return nil, err
	}
	sslConfig.MaxVersion = VersionMap["TLSv1.3"]
	sslConfig.CAFile = sslConfigMap[common.SslCaFileKey]
	sslConfig.CertFile = sslConfigMap[common.SslCertFileKey]
	sslConfig.KeyFile = sslConfigMap[common.SslKeyFileKey]
	sslConfig.CertPWDFile = sslConfigMap[common.SslCertPwdFilePath]
	sslConfig.ServerName = sslConfigMap[common.SslServerNameKey]

	return sslConfig, nil
}

// GetSSLConfigByService get ssl configurations based on service
func GetSSLConfigByService(svcName, protocol, svcType string) (*SSLConfig, error) {
	tag, err := generateSSLTag(svcName, protocol, svcType)
	if err != nil {
		return nil, err
	}

	sslConfigMap := getSSLConfigMap(tag, protocol, svcType, svcName)
	if len(sslConfigMap) == 0 {
		return nil, errSSLConfigNotExist
	}

	sslConfig, err := parseSSLConfig(sslConfigMap)
	if err != nil {
		return nil, err
	}
	return sslConfig, nil
}

// GetDefaultSSLConfig get default ssl configurations
func GetDefaultSSLConfig() *SSLConfig {
	sslConfigMap := getDefaultSslConfigMap()
	sslConfig, _ := parseSSLConfig(sslConfigMap)
	return sslConfig
}

// generateSSLTag generate ssl tag
func generateSSLTag(svcName, protocol, svcType string) (string, error) {
	var tag string
	if svcName != "" {
		tag = tag + "." + svcName
	}
	if protocol != "" {
		tag = tag + "." + protocol
	}
	if tag == "" {
		return "", errors.New("Service name and protocol can't be empty both")
	}

	switch svcType {
	case common.Consumer, common.Provider:
		tag = tag + "." + svcType
	default:
		return "", fmt.Errorf("Service type not support: %s, must be: %s|%s",
			svcType, common.Provider, common.Consumer)
	}

	return tag[1:], nil
}

// GetTLSConfigByService get tls configurations based on service
func GetTLSConfigByService(svcName, protocol, svcType string) (*tls.Config, *SSLConfig, error) {
	sslConfig, err := GetSSLConfigByService(svcName, protocol, svcType)
	if err != nil {
		return nil, nil, err
	}

	var tlsConfig *tls.Config
	switch svcType {
	case common.Provider:
		tlsConfig, err = GetServerTLSConfig(sslConfig)
	case common.Consumer:
		tlsConfig, err = GetClientTLSConfig(sslConfig)
	default:
		err = fmt.Errorf("service type not support: %s, must be: %s|%s",
			svcType, common.Provider, common.Consumer)
	}
	if err != nil {
		return nil, sslConfig, err
	}

	return tlsConfig, sslConfig, nil
}

// IsSSLConfigNotExist check the status of ssl configurations
func IsSSLConfigNotExist(e error) bool {
	return errors.Is(e, errSSLConfigNotExist)
}

// GetTLSConfig returns tls config from scheme and type
func GetTLSConfig(scheme, t string) (*tls.Config, error) {
	var tlsConfig *tls.Config
	secure := scheme == common.HTTPS
	if secure {
		sslTag := t + "." + common.Consumer
		tmpTLSConfig, _, err := GetTLSConfigByService(t, "", common.Consumer)
		if err != nil {
			if IsSSLConfigNotExist(err) {
				return nil, fmt.Errorf("%s tls mode, but no ssl config", sslTag)
			}
			return nil, fmt.Errorf("Load %s TLS config failed", sslTag)
		}
		tlsConfig = tmpTLSConfig
	}
	return tlsConfig, nil
}
