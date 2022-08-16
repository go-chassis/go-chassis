package tls

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chassis/cari/security"
	"github.com/go-chassis/foundation/tlsutil"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/security/cipher"

	//this import used for plain cipher
	_ "github.com/go-chassis/go-chassis/v2/security/cipher/plugins/plain"
)

// SSLConfig struct stores the necessary info for SSL configuration
type SSLConfig struct {
	CipherPlugin string   `yaml:"cipher_plugin" json:"cipherPlugin"`
	VerifyPeer   bool     `yaml:"verify_peer" json:"verifyPeer"`
	CipherSuites []uint16 `yaml:"cipher_suites" json:"cipherSuits"`
	MinVersion   uint16   `yaml:"min_version" json:"minVersion"`
	MaxVersion   uint16   `yaml:"max_version" json:"maxVersion"`
	CAFile       string   `yaml:"ca_file" json:"caFile"`
	CertFile     string   `yaml:"cert_file" json:"certFile"`
	KeyFile      string   `yaml:"key_file" json:"keyFile"`
	CertPWDFile  string   `yaml:"cert_pwd_file" json:"certPwdFile"`
	ServerName   string   `yaml:"server_name" json:"serverName"`
}

// TLSCipherSuiteMap is a map with key of type string and value of type unsigned integer
var TLSCipherSuiteMap = map[string]uint16{
	"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256": tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384": tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
}

// VersionMap is a map with key of type string and value of type unsigned integer
var VersionMap = map[string]uint16{
	"TLSv1.0": tls.VersionTLS10,
	"TLSv1.1": tls.VersionTLS11,
	"TLSv1.2": tls.VersionTLS12,
	"TLSv1.3": tls.VersionTLS13,
}

// GetX509CACertPool read a certificate file and gets the certificate configuration
func GetX509CACertPool(caCertFile string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	caCert, err := os.ReadFile(filepath.Clean(caCertFile))
	if err != nil {
		return nil, fmt.Errorf("read ca cert file %s failed", caCert)
	}

	pool.AppendCertsFromPEM(caCert)
	return pool, nil
}

func getTLSConfig(sslConfig *SSLConfig, role string) (tlsConfig *tls.Config, err error) {
	clientAuthMode := tls.NoClientCert
	var pool *x509.CertPool
	// ca file is needed when veryPeer is true
	if sslConfig.VerifyPeer {
		pool, err = GetX509CACertPool(sslConfig.CAFile)
		if err != nil {
			return nil, err
		}

		clientAuthMode = tls.RequireAndVerifyClientCert
	}

	// if cert pwd file is set, get the pwd
	var keyPassphase []byte
	if sslConfig.CertPWDFile != "" {
		keyPassphase, err = os.ReadFile(sslConfig.CertPWDFile)
		if err != nil {
			return nil, fmt.Errorf("read cert pwd %s failed: %w", sslConfig.CertPWDFile, err)
		}
	}

	// certificate is necessary for server, optional for client
	var certs []tls.Certificate
	if !(role == common.Client && sslConfig.KeyFile == "" && sslConfig.CertFile == "") {
		var cipherPlugin security.Cipher
		if cipherPlugin, err = cipher.NewCipher(sslConfig.CipherPlugin); err != nil {
			return nil, fmt.Errorf("get cipher plugin [%s] failed, %w", sslConfig.CipherPlugin, err)
		} else if cipherPlugin == nil {
			return nil, errors.New("invalid cipher plugin")
		}
		certs, err = tlsutil.LoadTLSCertificate(sslConfig.CertFile, sslConfig.KeyFile, strings.TrimSpace(string(keyPassphase)), func(src string) (string, error) {
			return cipherPlugin.Decrypt(src)
		})
		if err != nil {
			return nil, err
		}
	}

	switch role {
	case "server":
		tlsConfig = &tls.Config{
			ClientCAs:                pool,
			Certificates:             certs,
			CipherSuites:             sslConfig.CipherSuites,
			PreferServerCipherSuites: true,
			ClientAuth:               clientAuthMode,
			MinVersion:               sslConfig.MinVersion,
			MaxVersion:               sslConfig.MaxVersion,
		}
	case common.Client:
		tlsConfig = &tls.Config{
			RootCAs:            pool,
			Certificates:       certs,
			CipherSuites:       sslConfig.CipherSuites,
			InsecureSkipVerify: !sslConfig.VerifyPeer,
			MinVersion:         sslConfig.MinVersion,
			MaxVersion:         sslConfig.MaxVersion,
			ServerName:         sslConfig.ServerName,
		}
	}

	return tlsConfig, nil
}

// GetClientTLSConfig function gets client side TLS config
func GetClientTLSConfig(sslConfig *SSLConfig) (*tls.Config, error) {
	return getTLSConfig(sslConfig, "client")
}

// GetServerTLSConfig function gets server side TLD config
func GetServerTLSConfig(sslConfig *SSLConfig) (*tls.Config, error) {
	return getTLSConfig(sslConfig, "server")
}

// ParseSSLCipherSuites function parses cipher suites in to a list
func ParseSSLCipherSuites(ciphers string) ([]uint16, error) {
	cipherSuiteList := make([]uint16, 0)
	cipherSuiteNameList := strings.Split(ciphers, ",")
	for _, cipherSuiteName := range cipherSuiteNameList {
		cipherSuiteName = strings.TrimSpace(cipherSuiteName)
		if len(cipherSuiteName) == 0 {
			continue
		}

		if cipherSuite, ok := TLSCipherSuiteMap[cipherSuiteName]; ok {
			cipherSuiteList = append(cipherSuiteList, cipherSuite)
		} else {
			// 配置算法不存在
			return nil, fmt.Errorf("cipher %s not exist", cipherSuiteName)
		}
	}

	return cipherSuiteList, nil
}

// ParseSSLProtocol function parses SSL protocols
func ParseSSLProtocol(sprotocol string) (uint16, error) {
	var result uint16 = tls.VersionTLS12
	if protocol, ok := VersionMap[sprotocol]; ok {
		result = protocol
	} else {
		return result, fmt.Errorf("invalid ssl minimal version invalid(%s)", sprotocol)
	}

	return result, nil
}
