package archaius

// constant for hystrix parameters
const (
	DefaultForceFallback             = false
	DefaultTimeoutEnabled            = false
	DefaultCircuitBreakerEnabled     = true
	DefaultCircuitBreakerForceOpen   = false
	DefaultCircuitBreakerForceClosed = false
	DefaultFallbackEnable            = true
	DefaultMaxConcurrent             = 1000
	DefaultSleepWindow               = 15000
	DefaultTimeout                   = 30000
	DefaultErrorPercentThreshold     = 50
	DefaultRequestVolumeThreshold    = 20
	PolicyNull                       = "returnnull"
	PolicyThrowException             = "throwexception"
)

// GetForceFallback get force fallback
func GetForceFallback(command, t string) bool {
	return GetBool(GetForceFallbackKey(command), GetBool(GetDefaultForceFallbackKey(t), DefaultForceFallback))
}

//GetTimeoutEnabled get timeout enabled
func GetTimeoutEnabled(command, t string) bool {
	return GetBool(GetTimeEnabledKey(command), GetBool(GetDefaultTimeEnabledKey(t), DefaultTimeoutEnabled))
}

// GetTimeout get is to get timeout period
func GetTimeout(command, t string) int {
	return GetInt(GetTimeoutKey(command), GetInt(GetDefaultTimeoutKey(t), DefaultTimeout))
}

// GetMaxConcurrentRequests is to get maximum concurrent requests
func GetMaxConcurrentRequests(command, t string) int {
	return GetInt(GetMaxConcurrentKey(command), GetInt(GetDefaultMaxConcurrentKey(t), DefaultMaxConcurrent))
}

// GetErrorPercentThreshold is to get error percentage threshold
func GetErrorPercentThreshold(command, t string) int {
	return GetInt(GetErrorPercentThresholdKey(command), GetInt(GetDefaultErrorPercentThreshold(t), DefaultErrorPercentThreshold))
}

// GetRequestVolumeThreshold is to get request volume threshold
func GetRequestVolumeThreshold(command, t string) int {
	return GetInt(GetRequestVolumeThresholdKey(command), GetInt(GetDefaultRequestVolumeThresholdKey(t), DefaultRequestVolumeThreshold))
}

// GetSleepWindow get sleep window
func GetSleepWindow(command, t string) int {
	return GetInt(GetSleepWindowKey(command), GetInt(GetDefaultSleepWindowKey(t), DefaultSleepWindow))
}

// GetForceClose get force close
func GetForceClose(command, t string) bool {
	return GetBool(GetForceCloseKey(command), GetBool(GetDefaultForceCloseKey(t), DefaultCircuitBreakerForceClosed))
}

// GetForceOpen get force open
func GetForceOpen(command, t string) bool {
	return GetBool(GetForceOpenKey(command), GetBool(GetDefaultForceOpenKey(t), DefaultCircuitBreakerForceOpen))
}

// GetCircuitBreakerEnabled get circuit breaker enabled
func GetCircuitBreakerEnabled(command, t string) bool {
	return GetBool(GetCircuitBreakerEnabledKey(command), GetBool(GetDefaultCircuitBreakerEnabledKey(t), DefaultCircuitBreakerEnabled))
}

// GetFallbackEnabled get fallback enabled
func GetFallbackEnabled(command, t string) bool {
	return GetBool(GetFallbackEnabledKey(command), GetBool(GetDefaultGetFallbackEnabledKey(t), DefaultFallbackEnable))
}

// GetPolicy get policy
func GetPolicy(command, t string) string {
	return GetString(GetFallbackPolicyKey(command), GetString(GetDefaultFallbackPolicyKey(t), PolicyThrowException))
}
