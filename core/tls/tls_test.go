package tls_test

import (
	"github.com/go-chassis/go-archaius"

	"crypto/tls"
	"os"
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	chassisTLS "github.com/go-chassis/go-chassis/v2/core/tls"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	os.Setenv("CHASSIS_HOME", "/tmp")
	defer os.Unsetenv("CHASSIS_HOME")

	archaius.Init(archaius.WithMemorySource())
	archaius.Set("ssl.test.Consumer.certFile", "test.cer")
	archaius.Set("ssl.test.Consumer.keyFile", "test.key")
	config.ReadGlobalConfigFromArchaius()

	testConsumerSslConfig, err := chassisTLS.GetSSLConfigByService("test", "", common.Consumer)
	assert.NoError(t, err)
	assert.Nil(t, err)
	assert.Equal(t, "default", testConsumerSslConfig.CipherPlugin)
	assert.Equal(t, uint16(tls.VersionTLS12), testConsumerSslConfig.MinVersion)
	assert.Equal(t, "test.cer", testConsumerSslConfig.CertFile)
	assert.Equal(t, "test.key", testConsumerSslConfig.KeyFile)

	_, err = chassisTLS.GetSSLConfigByService("none", "", common.Consumer)
	assert.True(t, chassisTLS.IsSSLConfigNotExist(err))

	_, _, err = chassisTLS.GetTLSConfigByService("svcname", "protocol", "svctype")
	assert.Error(t, err)
	defaultCfg := chassisTLS.GetDefaultSSLConfig()
	assert.NotEmpty(t, defaultCfg)
	assert.Equal(t, uint16(tls.VersionTLS13), defaultCfg.MaxVersion)
}
