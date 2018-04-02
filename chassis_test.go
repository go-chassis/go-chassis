package chassis_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/util/fileutil"

	"github.com/stretchr/testify/assert"
)

const (
	Provider = "provider"
)

func TestInit(t *testing.T) {
	t.Log("Testing Chassis Init function")
	os.Setenv("CHASSIS_HOME", filepath.Join(os.Getenv("GOPATH"), "test", "chassisInit"))
	err := os.MkdirAll(fileutil.GetConfDir(), 0600)
	assert.NoError(t, err)
	globalDefFile, err := os.OpenFile(fileutil.GlobalDefinition(), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	defer globalDefFile.Close()

	// write some text line-by-line to file
	_, err = globalDefFile.WriteString(`---
#APPLICATION_ID: CSE optional

cse:
  flowcontrol:
    Consumer:
      qps:
        enabled: true
        limit:
          Server.EmployServer: 100
  loadbalance:
    strategy:
      name: RoundRobin
    retryEnabled: false
    retryOnNext: 2
    retryOnSame: 3
    backoff:
      kind: constant
      minMs: 200
      maxMs: 400
  service:
    registry:
      type: servicecenter
      scope: full
      address: http://127.0.0.1:30100
      refeshInterval : 30s
      watch: true
      register: reg
  handler:
    chain:
      consumer:
        default: bizkeeper-consumer, loadbalance, ratelimiter-consumer
  references:
    Server:
      version: 0.1
      transport: highway
    HelloServer:
      version: 0.1
      transport: highway
ssl:
  registry.consumer.cipherPlugin: default
  registry.consumer.verifyPeer: false
  registry.consumer.cipherSuits: TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
  registry.consumer.protocol: TLSv1.2
  registry.consumer.caFile:
  registry.consumer.certFile:
  registry.consumer.keyFile:
  registry.consumer.certPwdFile:
`)
	assert.NoError(t, err)
	msDefFile, err := os.OpenFile(fileutil.GetMicroserviceDesc(), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	assert.NoError(t, err)
	defer msDefFile.Close()
	_, err = msDefFile.WriteString(`---
#微服务的私有属性
service_description:
  name: nodejs2
  level: FRONT
  version: 0.1
  properties:
    allowCrossApp: true
  instance_properties:
    a: s
    p: s
`)

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	config.GlobalDefinition = &model.GlobalCfg{}

	config.Init()

	config.GlobalDefinition.Cse.Handler.Chain.Provider = map[string]string{
		"default": "bizkeeper-provider",
	}
	config.GlobalDefinition.Cse.Service.Registry.AutoRegister = "abc"

	err = chassis.Init()
	assert.NoError(t, err)

	chassis.RegisterSchema("rest", "str")
}
func TestInitError(t *testing.T) {
	t.Log("Testing chassis Init function for errors")
	p := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication/client")
	os.Setenv("CHASSIS_HOME", p)
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
}
