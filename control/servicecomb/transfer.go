package servicecomb

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/backoff"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
)

//SaveToLBCache save configs
func SaveToLBCache(raw *model.LoadBalancing) {
	openlogging.Debug("Loading lb config from archaius into cache")
	oldKeys := LBConfigCache.Items()
	newKeys := make(map[string]bool)
	// if there is no config, none key will be updated
	if raw != nil {
		newKeys = reloadLBCache(raw)
	}
	// remove outdated keys
	for old := range oldKeys {
		if _, ok := newKeys[old]; !ok {
			LBConfigCache.Delete(old)
		}
	}

}
func saveDefaultLB(raw *model.LoadBalancing) string { // return updated key
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
	return ""
}
func saveEachLB(k string, raw model.LoadBalancingSpec) string { // return updated key
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
	return k
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
	openlogging.Debug("Loading cb config from archaius into cache")
	oldKeys := CBConfigCache.Items()
	newKeys := make(map[string]bool)
	// if there is no config, none key will be updated
	if raw != nil {
		client.SetTimeoutToClientCache(raw.IsolationProperties)
		newKeys = reloadCBCache(raw)
	}
	// remove outdated keys
	for old := range oldKeys {
		if _, ok := newKeys[old]; !ok {
			CBConfigCache.Delete(old)
		}
	}
}

func saveEachCB(serviceName, serviceType string) string { //return updated key
	command := serviceType
	if serviceName != "" {
		command = strings.Join([]string{serviceType, serviceName}, ".")
	}
	c := hystrix.CommandConfig{
		ForceFallback:          config.GetForceFallback(serviceName, serviceType),
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
		openlogging.GetLogger().Infof(formatString, c, serviceName)
		CBConfigCache.Set(cbcCacheKey, c, 0)
		return cbcCacheKey
	}
	commandConfig, ok := cbcCacheValue.(hystrix.CommandConfig)
	if !ok {
		openlogging.GetLogger().Infof(formatString, c, serviceName)
		CBConfigCache.Set(cbcCacheKey, c, 0)
		return cbcCacheKey
	}
	if c == commandConfig {
		return cbcCacheKey
	}
	openlogging.GetLogger().Infof(formatString, c, serviceName)
	CBConfigCache.Set(cbcCacheKey, c, 0)
	return cbcCacheKey
}

//GetCBCacheKey generate cache key
func GetCBCacheKey(serviceName, serviceType string) string {
	key := serviceType
	if serviceName != "" {
		key = serviceType + ":" + serviceName
	}
	return key
}

func reloadLBCache(src *model.LoadBalancing) map[string]bool { //return updated keys
	keys := make(map[string]bool)
	k := saveDefaultLB(src)
	keys[k] = true
	if src.AnyService == nil {
		return keys
	}
	for name, conf := range src.AnyService {
		k = saveEachLB(name, conf)
		keys[k] = true
	}
	return keys
}

func reloadCBCache(src *model.HystrixConfig) map[string]bool { //return updated keys
	keys := make(map[string]bool)
	// global level config
	k := saveEachCB("", common.Consumer)
	keys[k] = true
	k = saveEachCB("", common.Provider)
	keys[k] = true
	// get all services who have configs
	consumers := make([]string, 0)
	providers := make([]string, 0)
	consumerMap := map[string]bool{}
	providerMap := map[string]bool{}

	// if a service has configurations of IsolationProperties|
	// CircuitBreakerProperties|FallbackPolicyProperties|FallbackProperties,
	// it's configuration should be added to cache when framework starts
	for _, p := range []interface{}{
		src.IsolationProperties,
		src.CircuitBreakerProperties,
		src.FallbackProperties,
		config.GetHystrixConfig().FallbackPolicyProperties} {
		if services, err := getServiceNamesByServiceTypeAndAnyService(p, common.Consumer); err != nil {
			openlogging.GetLogger().Errorf("Parse services from config failed: %v", err.Error())
		} else {
			consumers = append(consumers, services...)
		}
		if services, err := getServiceNamesByServiceTypeAndAnyService(p, common.Provider); err != nil {
			openlogging.GetLogger().Errorf("Parse services from config failed: %v", err.Error())
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
		k = saveEachCB(name, common.Consumer)
		keys[k] = true
	}
	for name := range providerMap {
		k = saveEachCB(name, common.Provider)
		keys[k] = true
	}
	return keys
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
