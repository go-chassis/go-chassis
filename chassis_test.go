package chassis_test

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

const (
	Provider = "provider"
)

func TestInit(t *testing.T) {
	t.Log("Testing Chassis Init function")

	path := "root/conf/chassis.yaml"
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		err := os.MkdirAll("root", 777)
		assert.NoError(t, err)
		err = os.MkdirAll("root/conf", 777)
		assert.NoError(t, err)
		file, err := os.Create("root/conf/chassis.yaml")
		assert.NoError(t, err)
		defer file.Close()
	}
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	assert.NoError(t, err)
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(`---
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
      sessionTimeoutInSeconds: 30
    retryEnabled: false
    retryOnNext: 2
    retryOnSame: 3
    backoff:
      kind: constant
      MinMs: 200
      MaxMs: 400
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
	path = filepath.Join("root", "conf", "microservice.yaml")
	_, err = os.Stat(path)
	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(filepath.Join("root", "conf", "microservice.yaml"))
		assert.NoError(t, err)
		defer file.Close()
	}
	file, err = os.OpenFile(path, os.O_RDWR, 0644)
	assert.NoError(t, err)
	defer file.Close()
	_, err = file.WriteString(`---
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
	assert.NoError(t, err)
	// save changes
	err = file.Sync()
	assert.NoError(t, err)

	os.Setenv("CHASSIS_HOME", filepath.Join("root"))
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

	err = os.Remove(path)
	assert.NoError(t, err)

	err = os.RemoveAll("root")
	assert.NoError(t, err)
}
func TestInitError(t *testing.T) {
	t.Log("Testing chassis Init function for errors")
	p := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication/client")
	os.Setenv("CHASSIS_HOME", p)
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

}
