package archaius

import (
	"strings"

	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/backoff"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
)

//SaveToLBCache save configs
func SaveToLBCache(raw *model.LoadBalancing) {
	lager.Logger.Debug("Loading lb config from archaius into cache")
	saveDefaultLB(raw)
	for k, v := range raw.AnyService {
		saveEachLB(k, v)
	}

}
func saveDefaultLB(raw *model.LoadBalancing) {
	c := control.LoadBalancingConfig{
		Strategy:                raw.Strategy["name"],
		RetryEnabled:            raw.RetryEnabled,
		RetryOnSame:             raw.RetryOnSame,
		RetryOnNext:             raw.RetryOnNext,
		BackOffKind:             raw.Backoff.Kind,
		BackOffMin:              raw.Backoff.MinMs,
		BackOffMax:              raw.Backoff.MaxMs,
		SessionTimeoutInSeconds: raw.SessionStickinessRule.SessionTimeoutInSeconds,
		SuccessiveFailedTimes:   raw.SessionStickinessRule.SuccessiveFailedTimes,
	}

	setDefaultLBValue(&c)
	LBConfigCache.Set("", c, 0)

}
func saveEachLB(k string, raw model.LoadBalancingSpec) {
	c := control.LoadBalancingConfig{
		Strategy:                raw.Strategy["name"],
		RetryEnabled:            raw.RetryEnabled,
		RetryOnSame:             raw.RetryOnSame,
		RetryOnNext:             raw.RetryOnNext,
		BackOffKind:             raw.Backoff.Kind,
		BackOffMin:              raw.Backoff.MinMs,
		BackOffMax:              raw.Backoff.MaxMs,
		SessionTimeoutInSeconds: raw.SessionStickinessRule.SessionTimeoutInSeconds,
		SuccessiveFailedTimes:   raw.SessionStickinessRule.SuccessiveFailedTimes,
	}
	setDefaultLBValue(&c)
	LBConfigCache.Set(k, c, 0)

}

func setDefaultLBValue(c *control.LoadBalancingConfig) {
	if c.Strategy == "" {
		c.Strategy = loadbalancer.StrategyRoundRobin
	}
	if c.BackOffKind == "" {
		c.BackOffKind = backoff.DefaultBackOffKind
	}
}

//SaveToCBCache save configs
func SaveToCBCache(raw *model.HystrixConfig, key string, isAnyService bool) {
	lager.Logger.Debug("Loading cb config from archaius into cache")
	saveEachCB("", common.Consumer)
	saveEachCB("", common.Provider)
	if !isAnyService {
		stringSlice := strings.Split(key, ".")
		saveEachCB(stringSlice[3], stringSlice[2])
	}
}

func saveEachCB(serviceName, serviceType string) {
	command := serviceType
	if serviceName != "" {
		command = strings.Join([]string{serviceType, serviceName}, ".")
	}
	c := hystrix.CommandConfig{
		ForceFallback:          config.GetForceFallback(serviceName, serviceType),
		TimeoutEnabled:         config.GetTimeoutEnabled(serviceName, serviceType),
		Timeout:                config.GetTimeout(command, serviceType),
		MaxConcurrentRequests:  config.GetMaxConcurrentRequests(command, serviceType),
		ErrorPercentThreshold:  config.GetErrorPercentThreshold(command, serviceType),
		RequestVolumeThreshold: config.GetRequestVolumeThreshold(command, serviceType),
		SleepWindow:            config.GetSleepWindow(command, serviceType),
		ForceClose:             config.GetForceClose(serviceName, serviceType),
		ForceOpen:              config.GetForceOpen(serviceName, serviceType),
		CircuitBreakerEnabled:  config.GetCircuitBreakerEnabled(command, serviceType),
	}
	cbcCacheKey := GetCBCacheKey(serviceName, serviceType)
	cbcCacheValue, b := CBConfigCache.Get(cbcCacheKey)
	if !b || cbcCacheValue == nil {
		lager.Logger.Infof("save circuit breaker config [%#v] for [%s] ", c, serviceName)
		CBConfigCache.Set(cbcCacheKey, c, 0)
		return
	}
	commandConfig, ok := cbcCacheValue.(hystrix.CommandConfig)
	if !ok {
		lager.Logger.Infof("save circuit breaker config [%#v] for [%s] ", c, serviceName)
		CBConfigCache.Set(cbcCacheKey, c, 0)
		return
	}
	if c == commandConfig {
		return
	}
	lager.Logger.Infof("save circuit breaker config [%#v] for [%s] ", c, serviceName)
	CBConfigCache.Set(cbcCacheKey, c, 0)
}

//GetCBCacheKey generate cache key
func GetCBCacheKey(serviceName, serviceType string) string {
	key := serviceType
	if serviceName != "" {
		key = serviceType + ":" + serviceName
	}
	return key
}
