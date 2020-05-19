package chassis_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"

	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/stretchr/testify/assert"
	"syscall"
)

const (
	Provider = "provider"
)

func TestInit(t *testing.T) {
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	t.Log("Testing Chassis Init function")
	os.Setenv("CHASSIS_HOME", filepath.Join(os.Getenv("GOPATH"), "test", "chassisInit"))
	err := os.MkdirAll(fileutil.GetConfDir(), 0700)
	assert.NoError(t, err)
	globalDefFile, err := os.OpenFile(fileutil.GlobalConfigPath(), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0700)
	defer globalDefFile.Close()

	// write some text line-by-line to file
	_, err = globalDefFile.WriteString(`---
#APPLICATION_ID: CSE optional
controlPanel:
  infra: istio
  settings:
    Address: xxx
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
  protocols:
    rest:
      listenAddress: 127.0.0.1:5001
  handler:
    chain:
      Consumer:
        rest: bizkeeper-consumer, loadbalance
      Provider:
        rest: bizkeeper-provider
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
	msDefFile, err := os.OpenFile(fileutil.MicroServiceConfigPath(), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0700)
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

	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})

	config.GlobalDefinition = &model.GlobalCfg{}

	config.Init()

	config.GlobalDefinition.Cse.Service.Registry.AutoRegister = "abc"

	chassis.SetDefaultConsumerChains(nil)
	chassis.SetDefaultProviderChains(nil)

	sigs := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGILL, syscall.SIGTRAP, syscall.SIGABRT}

	chassis.HackSignal(sigs...)

	chassis.InstalPreShutdown("pre_test", func(){
		t.Log("pre_shutdown_test")
	})

	chassis.InstalPostShutdown("pre_test", func(){
		t.Log("post_shutdown_test")
	})

	chassis.HackGracefulShutdown(chassis.GracefulShutdown)

	err = chassis.Init()
	assert.NoError(t, err)

	chassis.RegisterSchema("rest", "str")

	restServer, err := server.GetServer("rest")
	assert.NotNil(t, restServer)
	assert.NoError(t, err)

	v := reflect.ValueOf(restServer)
	opts := reflect.Indirect(v).FieldByName("opts")
	chainName := opts.FieldByName("ChainName")
	assert.Equal(t, "rest", chainName.String())

}

func TestInitError(t *testing.T) {
	t.Log("Testing chassis Init function for errors")
	p := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis", "examples", "communication/client")
	os.Setenv("CHASSIS_HOME", p)

	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
