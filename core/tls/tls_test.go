package tls_test

import (
	"crypto/tls"
	"io"
	"os"
	"testing"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
	"github.com/stretchr/testify/assert"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestInit(t *testing.T) {
	os.Setenv("CHASSIS_HOME", "/tmp")

	yamlContent := "a:\n  b:\n    c: valueC\n    d: valueD\n  \ryamlkeytest1: test1"
	chassisyamlContent := "APPLICATION_ID: CSE\n  \ncse:\n  service:\n    registry:\n      type: servicecenter\n  protocols:\n       highway:\n         listenAddress: 127.0.0.1:8080\n  \nssl:\n  test.Consumer.certFile: test.cer\n  test.Consumer.keyFile: test.key\n"
	os.Args = append(os.Args, "--argument=cmdtest")

	confdir := "/tmp/conf"
	filename1 := "/tmp/conf/chassis.yaml"
	filename2 := "/tmp/conf/circuit_breaker.yaml"
	filename3 := "/tmp/conf/lager.yaml"
	filename4 := "/tmp/conf/chassis.yaml"
	filename5 := "/tmp/conf/microservice.yaml"

	os.Remove(filename1)
	os.Remove(filename2)
	os.Remove(filename3)
	os.Remove(filename4)
	os.Remove(filename5)
	err := os.MkdirAll(confdir, 0777)
	check(err)

	f1, err1 := os.Create(filename1)
	check(err1)
	f2, err2 := os.Create(filename2)
	check(err2)
	f3, err3 := os.Create(filename3)
	check(err3)
	_, err4 := os.Create(filename4)
	check(err4)
	_, err5 := os.Create(filename5)
	check(err5)
	_, err1 = io.WriteString(f1, chassisyamlContent)
	_, err1 = io.WriteString(f2, yamlContent)
	_, err1 = io.WriteString(f3, yamlContent)

	config.Init()

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
	defaultCnfg := chassisTLS.GetDefaultSSLConfig()
	assert.NotEmpty(t, defaultCnfg)
	f1.Close()
	f2.Close()
	f3.Close()
	os.Remove(filename1)
	os.Remove(filename2)
	os.Remove(filename3)
	os.Remove(confdir)
}
