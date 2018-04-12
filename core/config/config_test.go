package config_test

import (
	"os"
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func TestInit(t *testing.T) {
	t.Log("testing config initialization")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	//config.Init()

	err := config.Init()
	assert.NoError(t, err)

	assert.Equal(t, "default", config.GlobalDefinition.AppID)
	assert.Equal(t, "servicecenter", config.GlobalDefinition.Cse.Service.Registry.Type)
	assert.Equal(t, "127.0.0.1:8082", config.GlobalDefinition.Cse.Protocols["highway"].Listen)

}

func TestInit2(t *testing.T) {
	file := []byte(`
cse:
  isolation:
    Consumer:
      timeout:
        enabled: true
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
	s, _ := c.String()
	t.Log(string(s))
	assert.Equal(t, 20, c.HystrixConfig.FallbackProperties.Consumer.MaxConcurrentRequests)
	assert.Equal(t, true, c.HystrixConfig.IsolationProperties.Consumer.TimeoutEnable.Enabled)
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
