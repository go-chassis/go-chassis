package archaius

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-chassis/go-archaius"
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
func SaveToLBCache(raw *model.LoadBalancing, key string, isAnyService bool) {
	lager.Logger.Debug("Loading lb config from archaius into cache")
	saveDefaultLB(raw)
	for k, v := range raw.AnyService {
		saveEachLB(k, v)
	}
	if !isAnyService {
		stringSlice := strings.Split(key, ".")
		if strings.Contains(key, "strategy.name") {
			value := archaius.Get(key)
			if value != nil {
				saveEachLB(stringSlice[2], raw.AnyService[stringSlice[2]])
			} else {
				LBConfigCache.Delete(stringSlice[2])
			}
		}

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
	formatString := "save circuit breaker config [%#v] for [%s] "
	if !b || cbcCacheValue == nil {
		lager.Logger.Infof(formatString, c, serviceName)
		CBConfigCache.Set(cbcCacheKey, c, 0)
		return
	}
	commandConfig, ok := cbcCacheValue.(hystrix.CommandConfig)
	if !ok {
		lager.Logger.Infof(formatString, c, serviceName)
		CBConfigCache.Set(cbcCacheKey, c, 0)
		return
	}
	if c == commandConfig {
		return
	}
	lager.Logger.Infof(formatString, c, serviceName)
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

func initLBCache() {
	src := config.GetLoadBalancing()
	saveDefaultLB(src)
	if src.AnyService == nil {
		return
	}
	for name, conf := range src.AnyService {
		saveEachLB(name, conf)
	}
}

func initCBCache() {
	// global level config
	saveEachCB("", common.Consumer)
	saveEachCB("", common.Provider)
	// get all services who have configs
	consumers := make([]string, 0)
	providers := make([]string, 0)
	consumerMap := map[string]bool{}
	providerMap := map[string]bool{}

	// if a service has configurations of IsolationProperties|
	// CircuitBreakerProperties|FallbackPolicyProperties|FallbackProperties,
	// it's configuration should be added to cache when framework starts
	for _, p := range []interface{}{
		config.GetHystrixConfig().IsolationProperties,
		config.GetHystrixConfig().CircuitBreakerProperties,
		config.GetHystrixConfig().FallbackProperties,
		config.GetHystrixConfig().FallbackPolicyProperties} {
		if services, err := getServiceNamesByServiceTypeAndAnyService(p, common.Consumer); err != nil {
			lager.Logger.Errorf("Parse services from config failed: %v", err.Error())
		} else {
			consumers = append(consumers, services...)
		}
		if services, err := getServiceNamesByServiceTypeAndAnyService(p, common.Provider); err != nil {
			lager.Logger.Errorf("Parse services from config failed: %v", err.Error())
		} else {
			providers = append(providers, services...)
		}
	}
	// remove duplicate service names
	for _, name := range consumers {
		consumerMap[name] = true
	}
	for _, name := range providers {
		providerMap[name] = true
	}
	// service level config
	for name := range consumerMap {
		saveEachCB(name, common.Consumer)
	}
	for name := range providerMap {
		saveEachCB(name, common.Provider)
	}
}

func getServiceNamesByServiceTypeAndAnyService(i interface{}, serviceType string) (services []string, err error) {
	// check type
	tmpType := reflect.TypeOf(i)
	if tmpType.Kind() != reflect.Ptr {
		return nil, errors.New("input must be an ptr")
	}
	// check value
	tmpValue := reflect.ValueOf(i)
	if !tmpValue.IsValid() {
		return []string{}, nil
	}

	inType := tmpType.Elem()
	propertyName := inType.Name()

	formatFieldNotExist := "field %s not exist"
	formatFieldNotExpected := "field %s is not type %s"
	// check type
	tmpFieldType, ok := inType.FieldByName(serviceType)
	if !ok {
		return nil, fmt.Errorf(formatFieldNotExist, propertyName+"."+serviceType)
	}
	if tmpFieldType.Type.Kind() != reflect.Ptr {
		return nil, fmt.Errorf(formatFieldNotExpected, propertyName+"."+serviceType, reflect.Ptr)
	}
	// check value
	inValue := reflect.Indirect(tmpValue)
	tmpFieldValue := inValue.FieldByName(serviceType)
	if !tmpFieldValue.IsValid() {
		return []string{}, nil
	}

	anyServiceFieldName := "AnyService"
	//check type
	fieldType := tmpFieldType.Type.Elem()
	tmpAnyServiceFieldType, ok := fieldType.FieldByName(anyServiceFieldName)
	if !ok {
		return nil, fmt.Errorf(formatFieldNotExist, propertyName+"."+serviceType+"."+anyServiceFieldName)
	}
	if tmpAnyServiceFieldType.Type.Kind() != reflect.Map {
		return nil, fmt.Errorf(formatFieldNotExpected, propertyName+"."+serviceType+"."+anyServiceFieldName, reflect.Map)
	}
	// check value
	fieldValue := reflect.Indirect(tmpFieldValue)
	anyServiceFieldValue := fieldValue.FieldByName(anyServiceFieldName)
	if !anyServiceFieldValue.IsValid() {
		return []string{}, nil
	}

	// get service names
	names := anyServiceFieldValue.MapKeys()
	services = make([]string, 0)
	for _, name := range names {
		if name.Kind() != reflect.String {
			return nil, fmt.Errorf(formatFieldNotExpected, "key of "+propertyName+"."+serviceType+"."+anyServiceFieldName, reflect.String)
		}
		services = append(services, name.String())
	}
	return services, nil
}
