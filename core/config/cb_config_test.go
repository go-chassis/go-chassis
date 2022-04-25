package config_test

import (
	"os"
	"testing"

	_ "github.com/go-chassis/go-chassis/v2/initiator"

	"io"
	"path/filepath"
	"time"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/pkg/util/fileutil"
	"github.com/stretchr/testify/assert"
)

func TestCBInit(t *testing.T) {
	b := []byte(`
---
cse:
  isolation:
    Consumer:
      timeoutInMilliseconds: 10
      maxConcurrentRequests: 100
      Server:
        timeoutInMilliseconds: 1
        maxConcurrentRequests: 10
  circuitBreaker:
    Consumer:
      enabled: true
      forceOpen: false
      forceClosed: false
      sleepWindowInMilliseconds: 10000
      requestVolumeThreshold: 30
      errorThresholdPercentage: 30
      Server:
        enabled: true
        forceOpen: false
        forceClosed: false
        sleepWindowInMilliseconds: 1000
        requestVolumeThreshold: 3
        errorThresholdPercentage: 3
  #容错处理函数，目前暂时按照开源的方式来不进行区分处理，统一调用fallback函数
  fallback:
    Consumer:
      enabled: true
      force: true
      Server:
        force: false
  fallbackpolicy:
    Consumer:
      policy: throwexception
      Server:
        policy: nil
`)
	d, _ := os.Getwd()
	filename1 := filepath.Join(d, "circuit_breaker.yaml")
	f1, err := os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	assert.NoError(t, err)
	_, err = f1.Write(b)
	assert.NoError(t, err)

	b = []byte(`
servicecomb:
  registry:
      type: servicecenter           #optional:可选zookeeper/servicecenter，zookeeper供中软使用，不配置的情况下默认为servicecenter
      scope: full                   #optional:scope不为full时，只允许在本app间访问，不允许跨app访问；为full就是注册时允许跨app，并且发现本租户全部微服务
      address: http://127.0.0.1:30100
      refreshInterval : 30s
      watch: true
  protocols:
    rest:
      listenAddress: 127.0.0.1:8081
      advertiseAddress: 127.0.0.1:8081
  handler:
    chain:
      Consumer:
        default: bizkeeper-consumer,router,loadbalance,tracing-consumer,ratelimiter-consumer,transport
  transport:
    failure:
      rest: http_500,http_502
    maxIdleCon:
      rest: 1024
`)
	d, _ = os.Getwd()
	filename1 = filepath.Join(d, "chassis.yaml")
	f1, err = os.OpenFile(filename1, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	assert.NoError(t, err)
	_, err = f1.Write(b)
	assert.NoError(t, err)

	b = []byte(`
servicecomb:
  service:
    name: Client

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
	defer os.Unsetenv(fileutil.ChassisConfDir)
	err = config.Init()
	assert.NoError(t, err)

	t.Run("TestGetFallbackEnabled", func(t *testing.T) {
		check := config.GetFallbackEnabled("Consumer.test", common.Consumer)
		assert.Equal(t, true, check)
	})
	t.Run("TestGetCircuitBreakerEnabled", func(t *testing.T) {
		check := config.GetCircuitBreakerEnabled("Consumer.test", common.Consumer)
		assert.Equal(t, true, check)
		check = config.GetCircuitBreakerEnabled(common.Consumer, common.Consumer)
		assert.Equal(t, true, check)
	})
	t.Run("TestGetForceOpen", func(t *testing.T) {
		check := config.GetForceOpen("test", common.Consumer)
		assert.Equal(t, false, check)
		check = config.GetForceOpen("Server", common.Consumer)
		assert.Equal(t, false, check)
	})
	t.Run("TestGetForceClose", func(t *testing.T) {
		config.HystrixConfig.HystrixConfig.CircuitBreakerProperties.Consumer.ForceClose = true
		check := config.GetForceClose("test", common.Consumer)
		assert.Equal(t, true, check)
		check = config.GetForceClose("Server", common.Consumer)
		assert.Equal(t, false, check)
	})
	t.Run("TestTimeout", func(t *testing.T) {
		check := config.GetTimeout("Consumer.test", common.Consumer)
		assert.Equal(t, 10, check)

		check = config.GetTimeout("Consumer.Server", common.Consumer)
		assert.Equal(t, 1, check)

		config.GetHystrixConfig().IsolationProperties.Consumer.TimeoutInMilliseconds = 0
		check = config.GetTimeout("Consumer.some", common.Consumer)
		assert.Equal(t, config.DefaultTimeout, check)

		d := config.GetTimeoutDuration("Consumer.some", common.Consumer)
		assert.Equal(t, config.DefaultTimeout*time.Millisecond, d)

		check = config.GetTimeout("Provider.some", common.Provider)
		assert.Equal(t, config.DefaultTimeout, check)

	})
	t.Run("TestTimeout in archaius", func(t *testing.T) {
		check := config.GetTimeoutDurationFromArchaius("Consumer.test", common.Consumer)
		assert.Equal(t, 10*time.Millisecond, check)

		check = config.GetTimeoutDurationFromArchaius("Consumer.Server", common.Consumer)
		assert.Equal(t, 1*time.Millisecond, check)

		check = config.GetTimeoutDurationFromArchaius("Consumer.some", common.Consumer)
		assert.Equal(t, 10*time.Millisecond, check)

		check = config.GetTimeoutDurationFromArchaius("Provider.some", common.Provider)
		assert.Equal(t, config.DefaultTimeout*time.Millisecond, check)

	})
	t.Run("TestGetMaxConcurrentRequests", func(t *testing.T) {
		check := config.GetMaxConcurrentRequests("Consumer.test", common.Consumer)
		assert.Equal(t, 100, check)
		check = config.GetMaxConcurrentRequests("Consumer.Server", common.Consumer)
		assert.Equal(t, 10, check)

		config.GetHystrixConfig().IsolationProperties.Consumer.MaxConcurrentRequests = 0
		check = config.GetMaxConcurrentRequests("Consumer.some", common.Consumer)
		assert.Equal(t, config.DefaultMaxConcurrent, check)
	})
	t.Run("TestGetSleepWindow",
		func(t *testing.T) {
			check := config.GetSleepWindow("Consumer.test", common.Consumer)
			assert.Equal(t, 10000, check)
			check = config.GetSleepWindow("Consumer.Server", common.Consumer)
			assert.Equal(t, 1000, check)

			config.GetHystrixConfig().CircuitBreakerProperties.Consumer.SleepWindowInMilliseconds = 0
			check = config.GetSleepWindow("Consumer.some", common.Consumer)
			assert.Equal(t, config.DefaultSleepWindow, check)
		})
	t.Run("TestGetRequestVolumeThreshold", func(t *testing.T) {
		check := config.GetRequestVolumeThreshold("Consumer.test", common.Consumer)
		assert.Equal(t, 30, check)

		k := config.GetRequestVolumeThresholdKey("Consumer.Server")
		t.Log(k)

		check = config.GetRequestVolumeThreshold("Consumer.Server", common.Consumer)
		assert.Equal(t, 3, check)

		config.GetHystrixConfig().CircuitBreakerProperties.Consumer.RequestVolumeThreshold = 0
		check = config.GetRequestVolumeThreshold("Consumer.test", common.Consumer)
		assert.Equal(t, config.DefaultRequestVolumeThreshold, check)
	})

	t.Run("TestGetErrorPercentThreshold",
		func(t *testing.T) {
			check := config.GetErrorPercentThreshold("Consumer.test", common.Consumer)
			assert.Equal(t, 30, check)
			check = config.GetErrorPercentThreshold("Consumer.Server", common.Consumer)
			assert.Equal(t, 3, check)

			config.GetHystrixConfig().CircuitBreakerProperties.Consumer.ErrorThresholdPercentage = 0
			check = config.GetErrorPercentThreshold("Consumer.some", common.Consumer)
			assert.Equal(t, 50, check)

			check = config.GetErrorPercentThreshold("Provider.some", common.Provider)
			assert.Equal(t, 50, check)
		})
	t.Run("TestGetPolicy", func(t *testing.T) {
		check := config.GetPolicy("test", common.Consumer)
		assert.Equal(t, "throwexception", check)
		check = config.GetPolicy("Server", common.Consumer)
		assert.Equal(t, "nil", check)

		check = config.GetPolicy("Server", common.Provider)
		assert.Equal(t, "throwexception", check)
	})
	t.Run("TestGetForceFallback",
		func(t *testing.T) {
			check := config.GetForceFallback("test", common.Consumer)
			assert.True(t, check)
			check = config.GetForceFallback("Server", common.Consumer)
			assert.False(t, check)

			check = config.GetForceFallback("Server", common.Provider)
			assert.False(t, check)

		})
}
