package config

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/cenkalti/backoff"
	"strings"
	"sync"
	"time"
)

const (
	lbPrefix                                 = "cse.loadbalance"
	propertyStrategyName                     = "strategy.name"
	propertySessionStickinessRuleTimeout     = "SessionStickinessRule.sessionTimeoutInSeconds"
	propertySessionStickinessRuleFailedTimes = "SessionStickinessRule.successiveFailedTimes"
	propertyRetryEnabled                     = "retryEnabled"
	propertyRetryOnNext                      = "retryOnNext"
	propertyRetryOnSame                      = "retryOnSame"
	propertyBackoffKind                      = "backoff.kind"
	propertyBackoffMinMs                     = "backoff.minMs"
	propertyBackoffMaxMs                     = "backoff.maxMs"

	//DefaultStrategy is default value for strategy
	DefaultStrategy = "RoundRobin"
	//DefaultSessionTimeout is default value for timeout
	DefaultSessionTimeout = 30
	//DefaultFailedTimes is default value for failed times
	DefaultFailedTimes = 5
	backoffJittered    = "jittered"
	backoffConstant    = "constant"
	backoffZero        = "zero"
	//DefaultBackoffKind is zero
	DefaultBackoffKind = backoffZero
)

var lbMutex = sync.RWMutex{}

func genKey(s ...string) string {
	return strings.Join(s, ".")
}

func genMsKey(prefix, src, dest, property string) string {
	if src == "" {
		return genKey(prefix, dest, property)
	}
	return genKey(prefix, src, dest, property)
}

// GetServerListFilters get server list filters
func GetServerListFilters() (filters []string) {
	lbMutex.RLock()
	filters = strings.Split(GetLoadBalancing().Filters, ",")
	lbMutex.RUnlock()
	return
}

// GetStrategyName get strategy name
func GetStrategyName(source, service string) string {
	lbMutex.RLock()
	r := GetLoadBalancing().AnyService[service].Strategy["name"]
	if r == "" {
		r = GetLoadBalancing().Strategy["name"]
		if r == "" {
			r = DefaultStrategy
		}
	}
	lbMutex.RUnlock()
	return r
}

// GetSessionTimeout return session timeout
func GetSessionTimeout(source, service string) int {
	lbMutex.RLock()
	global := GetLoadBalancing().SessionStickinessRule.SessionTimeoutInSeconds
	if global == 0 {
		global = DefaultSessionTimeout
	}
	ms := archaius.GetInt(genMsKey(lbPrefix, source, service, propertySessionStickinessRuleTimeout), global)
	lbMutex.RUnlock()
	return ms
}

// StrategySuccessiveFailedTimes strategy successive failed times
func StrategySuccessiveFailedTimes(source, service string) int {
	lbMutex.RLock()
	global := GetLoadBalancing().SessionStickinessRule.SuccessiveFailedTimes
	if global == 0 {
		global = DefaultFailedTimes
	}
	ms := archaius.GetInt(genMsKey(lbPrefix, source, service, propertySessionStickinessRuleFailedTimes), global)
	lbMutex.RUnlock()
	return ms
}

// RetryEnabled retry enabled
func RetryEnabled(source, service string) bool {
	lbMutex.RLock()
	global := GetLoadBalancing().RetryEnabled
	ms := archaius.GetBool(genMsKey(lbPrefix, source, service, propertyRetryEnabled), global)
	lbMutex.RUnlock()
	return ms
}

//GetRetryOnNext return value of GetRetryOnNext
func GetRetryOnNext(source, service string) int {
	lbMutex.RLock()
	global := GetLoadBalancing().RetryOnNext
	ms := archaius.GetInt(genMsKey(lbPrefix, source, service, propertyRetryOnNext), global)
	lbMutex.RUnlock()
	return ms
}

//GetRetryOnSame return value of RetryOnSame
func GetRetryOnSame(source, service string) int {
	lbMutex.RLock()
	global := GetLoadBalancing().RetryOnSame
	ms := archaius.GetInt(genMsKey(lbPrefix, source, service, propertyRetryOnSame), global)
	lbMutex.RUnlock()
	return ms
}

func backoffKind(source, service string) string {
	r := GetLoadBalancing().AnyService[service].Backoff.Kind
	if r == "" {
		r = GetLoadBalancing().Backoff.Kind
		if r == "" {
			r = DefaultBackoffKind
		}
	}
	return r
}

func backoffMinMs(source, service string) int {
	global := GetLoadBalancing().Backoff.MinMs
	ms := archaius.GetInt(genMsKey(lbPrefix, source, service, propertyBackoffMinMs), global)
	return ms
}

func backoffMaxMs(source, service string) int {
	global := GetLoadBalancing().Backoff.MaxMs
	ms := archaius.GetInt(genMsKey(lbPrefix, source, service, propertyBackoffMaxMs), global)
	return ms
}

//GetBackOff return the the back off policy
func GetBackOff(source, service string) backoff.BackOff {
	lbMutex.RLock()
	backoffKind := backoffKind(source, service)
	backMin := backoffMinMs(source, service)
	backMax := backoffMaxMs(source, service)
	lbMutex.RUnlock()
	switch backoffKind {
	case backoffJittered:
		return &backoff.ExponentialBackOff{
			InitialInterval:     time.Duration(backMin) * time.Millisecond,
			RandomizationFactor: backoff.DefaultRandomizationFactor,
			Multiplier:          backoff.DefaultMultiplier,
			MaxInterval:         time.Duration(backMax) * time.Millisecond,
			MaxElapsedTime:      0,
			Clock:               backoff.SystemClock,
		}
	case backoffConstant:
		return backoff.NewConstantBackOff(time.Duration(backMin) * time.Millisecond)
	case backoffZero:
		return &backoff.ZeroBackOff{}
	default:
		return &backoff.ZeroBackOff{}
	}

}
