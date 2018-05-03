package config_test

import (
	"os"
	"testing"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/stretchr/testify/assert"
)

func TestCBInit(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	archaius.Init()
}

func TestGetFallbackEnabled(t *testing.T) {
	check := config.GetFallbackEnabled("test", common.Consumer)
	assert.Equal(t, false, check)
}

func TestGetCircuitBreakerEnabled(t *testing.T) {
	check := config.GetCircuitBreakerEnabled("test", common.Consumer)
	assert.Equal(t, true, check)
}

func TestGetTimeoutEnabled(t *testing.T) {
	check := config.GetTimeoutEnabled("test", common.Consumer)
	assert.Equal(t, true, check)
	check = config.GetTimeoutEnabled("Server", common.Consumer)
	assert.Equal(t, false, check)
	check = config.GetTimeoutEnabled("test", common.Provider)
	assert.Equal(t, false, check)
}

func TestGetForceOpen(t *testing.T) {
	check := config.GetForceOpen("test", common.Consumer)
	assert.Equal(t, false, check)
	check = config.GetForceOpen("Server", common.Consumer)
	assert.Equal(t, true, check)
}

func TestGetForceClose(t *testing.T) {
	config.HystrixConfig.HystrixConfig.CircuitBreakerProperties.Consumer.ForceClose = true
	check := config.GetForceClose("test", common.Consumer)
	assert.Equal(t, true, check)
	check = config.GetForceClose("Server", common.Consumer)
	assert.Equal(t, false, check)
}

func TestTimeout(t *testing.T) {
	check := config.GetTimeout("test", common.Consumer)
	assert.Equal(t, 10, check)
	check = config.GetTimeout("Server", common.Consumer)
	assert.Equal(t, 10, check)
}

func TestGetMaxConcurrentRequests(t *testing.T) {
	check := config.GetMaxConcurrentRequests("test", common.Consumer)
	assert.Equal(t, 100, check)
	check = config.GetMaxConcurrentRequests("Server", common.Consumer)
	assert.Equal(t, 100, check)
}

func TestGetSleepWindow(t *testing.T) {
	check := config.GetSleepWindow("test", common.Consumer)
	assert.Equal(t, 10000, check)
	check = config.GetSleepWindow("Server", common.Consumer)
	assert.Equal(t, 10000, check)
}

func TestGetRequestVolumeThreshold(t *testing.T) {
	check := config.GetRequestVolumeThreshold("test", common.Consumer)
	assert.Equal(t, 20, check)
	check = config.GetRequestVolumeThreshold("Server", common.Consumer)
	assert.Equal(t, 20, check)
}

func TestGetErrorPercentThresholdk(t *testing.T) {
	check := config.GetErrorPercentThreshold("test", common.Consumer)
	assert.Equal(t, 50, check)
	check = config.GetErrorPercentThreshold("Server", common.Consumer)
	assert.Equal(t, 50, check)
}

func TestGetPolicy(t *testing.T) {
	check := config.GetPolicy("test", common.Consumer)
	assert.Equal(t, "throwexception", check)
	check = config.GetPolicy("Server", common.Consumer)
	assert.Equal(t, "throwexception", check)
}

func TestGetForceFallback(t *testing.T) {
	check := config.GetForceFallback("test", common.Consumer)
	assert.False(t, check)
	check = config.GetForceFallback("Server", common.Consumer)
	assert.False(t, check)
}
