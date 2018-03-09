package archaius_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/stretchr/testify/assert"
)

func TestGetForceFallback(t *testing.T) {
	check := archaius.GetForceFallback("test", common.Consumer)
	assert.Equal(t, check, false)
}

func TestGetTimeoutEnabled(t *testing.T) {
	check := archaius.GetTimeoutEnabled("test", common.Consumer)
	assert.Equal(t, check, true)
}

func TestGetTimeout(t *testing.T) {
	check := archaius.GetTimeout("test", common.Consumer)
	assert.NotEqual(t, check, 0)
}

func TestGetMaxConcurrentRequests(t *testing.T) {
	check := archaius.GetMaxConcurrentRequests("test", common.Consumer)
	assert.NotEqual(t, check, 0)
}

func TestGetErrorPercentThresholdk(t *testing.T) {
	check := archaius.GetErrorPercentThreshold("test", common.Consumer)
	assert.NotEqual(t, check, 0)
}

func TestGetRequestVolumeThreshold(t *testing.T) {
	check := archaius.GetRequestVolumeThreshold("test", common.Consumer)
	assert.NotEqual(t, check, 0)
}

func TestGetSleepWindow(t *testing.T) {
	check := archaius.GetSleepWindow("test", common.Consumer)
	assert.NotEqual(t, check, 0)
}

func TestGetForceClose(t *testing.T) {
	check := archaius.GetForceClose("test", common.Consumer)
	assert.Equal(t, check, false)
}

func TestGetForceOpen(t *testing.T) {
	check := archaius.GetForceOpen("test", common.Consumer)
	assert.Equal(t, check, false)
}

func TestGetCircuitBreakerEnabled(t *testing.T) {
	check := archaius.GetCircuitBreakerEnabled("test", common.Consumer)
	assert.Equal(t, check, true)
}

func TestGetFallbackEnabled(t *testing.T) {
	check := archaius.GetFallbackEnabled("test", common.Consumer)
	assert.Equal(t, check, false)
}

func TestGetPolicy(t *testing.T) {
	check := archaius.GetPolicy("test", common.Consumer)
	assert.NotEqual(t, check, "")
}
