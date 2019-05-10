package config_test

import (
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"path/filepath"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func TestInit1(t *testing.T) {
	b := []byte(`
cse:
  service:
    registry:
      #disabled: false           optional:禁用注册发现选项，默认开始注册发现
      type: servicecenter           #optional:可选zookeeper/servicecenter，zookeeper供中软使用，不配置的情况下默认为servicecenter
      scope: full                   #optional:scope不为full时，只允许在本app间访问，不允许跨app访问；为full就是注册时允许跨app，并且发现本租户全部微服务
      address: http://127.0.0.1:30100
      #register: manual          optional：register不配置时默认为自动注册，可选参数有自动注册auto和手动注册manual
      refeshInterval : 30s
      watch: true
  transport:
    failure:
      rest: http_500,http_502
    maxIdleCon:
      rest: 1024
  protocols:
    rest:
      listenAddress: 127.0.0.1:8080
      advertiseAddress: 127.0.0.1:8080
region:
  name: test
  region: cn
  availableZone: 1
  
`)
	d, _ := os.Getwd()
	filename1 := filepath.Join(d, "chassis.yaml")
	f1, err := os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	assert.NoError(t, err)
	_, err = f1.Write(b)
	assert.NoError(t, err)

	b = []byte(`
---
#微服务的私有属性
#APPLICATION_ID: CSE #optional
service_description:
  name: Client
  #version: 0.1 #optional

`)
	d, _ = os.Getwd()
	filename1 = filepath.Join(d, "microservice.yaml")
	os.Remove(filename1)
	f1, err = os.Create(filename1)
	assert.NoError(t, err)
	defer f1.Close()
	_, err = io.WriteString(f1, string(b))
	assert.NoError(t, err)

	os.Setenv(fileutil.ChassisConfDir, d)
	err = config.Init()
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	config.ReadGlobalConfigFile()

	c := config.GetConfigCenterConf()
	assert.Equal(t, "", c.ServerURI)

	dc := config.GetDataCenter()
	assert.Equal(t, "test", dc.Name)
	tc := config.GetTransportConf()
	assert.Equal(t, 1, len(tc.MaxIdlCons))
}
func TestInit2(t *testing.T) {
	t.Log("testing config initialization")

	assert.Equal(t, "servicecenter", config.GlobalDefinition.Cse.Service.Registry.Type)
	assert.Equal(t, "127.0.0.1:8080", config.GlobalDefinition.Cse.Protocols["rest"].Listen)

}

func TestInit3(t *testing.T) {
	file := []byte(`
cse:
  isolation:
    Consumer:
      timeoutInMilliseconds: 10
      maxConcurrentRequests: 100
      Server:
        timeoutInMilliseconds: 1000
        maxConcurrentRequests: 100
    Provider:
      Server:
        timeoutInMilliseconds: 10
        maxConcurrentRequests: 100
  circuitBreaker:
    scope: service
    Consumer:
      enabled: true
      forceOpen: false
      forceClosed: true
      sleepWindowInMilliseconds: 10000
      requestVolumeThreshold: 20
      errorThresholdPercentage: 50
      Server:
        enabled: true
        forceOpen: false
        forceClosed: true
        sleepWindowInMilliseconds: 10000
        requestVolumeThreshold: 20
        errorThresholdPercentage: 50
    Provider:
      Server:
        enabled: true
        forceOpen: false
        forceClosed: true
        sleepWindowInMilliseconds: 10000
        requestVolumeThreshold: 20
        errorThresholdPercentage: 50
  fallback:
    Consumer:
      enabled: false
      maxConcurrentRequests: 20
  fallbackpolicy:
    Consumer:
      policy: throwexception
`)
	c := &model.HystrixConfigWrapper{}
	err := yaml.Unmarshal(file, c)
	assert.NoError(t, err)
	assert.Equal(t, "service", c.HystrixConfig.CircuitBreakerProperties.Scope)
	assert.Equal(t, 20, c.HystrixConfig.FallbackProperties.Consumer.MaxConcurrentRequests)
	assert.Equal(t, "throwexception", c.HystrixConfig.FallbackPolicyProperties.Consumer.Policy)
	assert.Equal(t, 50, c.HystrixConfig.CircuitBreakerProperties.Consumer.AnyService["Server"].ErrorThresholdPercentage)
	assert.NotEqual(t, nil, config.GetHystrixConfig())
}

func TestGetLoadBalancing(t *testing.T) {
	lbBytes := []byte(`
cse: 
  loadbalance: 
    TargetService: 
      backoff: 
        maxMs: 400
        minMs: 200
        kind: constant
      retryEnabled: false
      retryOnNext: 2
      retryOnSame: 3
      serverListFilters: zoneaware
      strategy: 
        name: WeightedResponse
    backoff: 
      maxMs: 400
      minMs: 200
      kind: constant
    retryEnabled: false
    retryOnNext: 2
    retryOnSame: 3
    serverListFilters: zoneaware
    strategy: 
      name: WeightedResponse

`)
	lbConfig := &model.LBWrapper{}
	err := yaml.Unmarshal(lbBytes, lbConfig)
	assert.NoError(t, err)
	assert.Equal(t, "WeightedResponse", lbConfig.Prefix.LBConfig.Strategy["name"])
	assert.Equal(t, loadbalancer.ZoneAware, lbConfig.Prefix.LBConfig.Filters)
	t.Log(lbConfig.Prefix.LBConfig.AnyService)
	assert.Equal(t, "WeightedResponse", lbConfig.Prefix.LBConfig.AnyService["TargetService"].Strategy["name"])

	assert.Equal(t, "WeightedResponse", lbConfig.Prefix.LBConfig.Strategy["name"])
	assert.NotEqual(t, nil, config.GetLoadBalancing())

}

func TestInitErrorWithBlankEnv(t *testing.T) {
	os.Setenv("CHASSIS_HOME", "")
	os.Setenv("CHASSIS_CONF_DIR", "")
	err := config.Init()
	t.Log(err)
	assert.Error(t, err)
}
