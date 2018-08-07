package archaius

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/backoff"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"strings"
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
func SaveToCBCache(raw *model.HystrixConfig) {
	lager.Logger.Debug("Loading cb config from archaius into cache")
	saveEachCB("", common.Consumer)
	saveEachCB("", common.Provider)
	// even though the input is duplicated in below loops but service name list is not certain unless we loop all namespace
	for k := range raw.IsolationProperties.Consumer.AnyService {
		saveEachCB(k, common.Consumer)
	}
	for k := range raw.IsolationProperties.Provider.AnyService {
		saveEachCB(k, common.Provider)
	}
	for k := range raw.CircuitBreakerProperties.Consumer.AnyService {
		saveEachCB(k, common.Consumer)
	}
	for k := range raw.CircuitBreakerProperties.Provider.AnyService {
		saveEachCB(k, common.Provider)
	}
	for k := range raw.FallbackPolicyProperties.Consumer.AnyService {
		saveEachCB(k, common.Consumer)
	}
	for k := range raw.FallbackPolicyProperties.Provider.AnyService {
		saveEachCB(k, common.Provider)
	}
	for k := range raw.FallbackProperties.Consumer.AnyService {
		saveEachCB(k, common.Consumer)
	}
	for k := range raw.FallbackProperties.Provider.AnyService {
		saveEachCB(k, common.Provider)
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

	CBConfigCache.Set(GetCBCacheKey(serviceName, serviceType), c, 0)
}

//GetCBCacheKey generate cache key
func GetCBCacheKey(serviceName, serviceType string) string {
	key := serviceType
	if serviceName != "" {
		key = serviceType + ":" + serviceName
	}
	return key
}
