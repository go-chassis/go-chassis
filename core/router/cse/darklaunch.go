package cse

import (
	"encoding/json"
	"strings"

	"errors"
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
)

// RouteDarkLaunchGovernSourceName is source name of dark launch configuration
const RouteDarkLaunchGovernSourceName = "RouteDarkLaunchGovernSource"

// RouteDarkLaunchGovernSourcePriority is priority of dark launch configuration
const RouteDarkLaunchGovernSourcePriority = 8

// constant for route rule keys
const (
	DarkLaunchPrefix   = "cse.darklaunch.policy."
	DarkLaunchTypeRule = "RULE"
	DarkLaunchTypeRate = "RATE"
)

var source *RouteDarkLaunchGovernSource

// RouteDarkLaunchGovernSource gets dark launch configs from global archaius,
// it cannot work until global archaius inits success,
// it keeps no data
type RouteDarkLaunchGovernSource struct {
	d             core.DynamicConfigCallback
	callbackCheck chan bool
	chanStatus    bool
}

//GetSourceName returns name of dark launch configuration
func (r *RouteDarkLaunchGovernSource) GetSourceName() string {
	return RouteDarkLaunchGovernSourceName
}

//GetConfigurations gets all dark launch configurations
func (r *RouteDarkLaunchGovernSource) GetConfigurations() (map[string]interface{}, error) {
	routerConfigs, err := getRouterConfigFromDarkLaunch()
	if err != nil {
		lager.Logger.Error("Get router config from dark launch failed", err)
		return nil, err
	}
	d := make(map[string]interface{}, 0)
	for k, v := range routerConfigs.Destinations {
		d[k] = v
	}
	return d, nil
}

// GetConfigurationsByDI implements ConfigSource.GetConfigurationsByDI
func (r *RouteDarkLaunchGovernSource) GetConfigurationsByDI(dimensionInfo string) (map[string]interface{}, error) {
	return nil, nil
}

//GetConfigurationByKey gets required dark launch configuration for a particular key
func (r *RouteDarkLaunchGovernSource) GetConfigurationByKey(k string) (interface{}, error) {
	s := archaius.GetString(DarkLaunchPrefix+k, "")
	rule := &model.DarkLaunchRule{}
	if err := json.Unmarshal([]byte(s), rule); err != nil {
		return nil, err
	}
	routeRules := DarkLaunchRule2RouteRule(rule)
	return routeRules, nil
}

//GetConfigurationByKeyAndDimensionInfo implements ConfigSource.GetConfigurationByKeyAndDimensionInfo
func (r *RouteDarkLaunchGovernSource) GetConfigurationByKeyAndDimensionInfo(key, dimensionInfo string) (interface{}, error) {
	return nil, nil
}

// AddDimensionInfo implements ConfigSource.AddDimensionInfo
func (r *RouteDarkLaunchGovernSource) AddDimensionInfo(dimensionInfo string) (map[string]string, error) {
	return nil, nil
}

//DynamicConfigHandler dynamically handles a dark launch configuration
func (r *RouteDarkLaunchGovernSource) DynamicConfigHandler(d core.DynamicConfigCallback) error {
	r.d = d
	r.callbackCheck <- true
	return nil
}

//GetPriority returns priority of the dark launch configuration
func (r *RouteDarkLaunchGovernSource) GetPriority() int {
	return RouteDarkLaunchGovernSourcePriority
}

//Cleanup implements ConfigSource.Cleanup
func (r *RouteDarkLaunchGovernSource) Cleanup() error { return nil }

// Callback callbacks when receive an event
// only operates after dynamicCallback is initialized
func (r *RouteDarkLaunchGovernSource) Callback(e *core.Event) {
	if !r.chanStatus {
		<-r.callbackCheck
		r.chanStatus = true
	}
	lager.Logger.Debugf("Get event, key: %s, type: %s", e.Key, e.EventType)
	if r.d == nil {
		lager.Logger.Warn("Dynamic config handler is nil", nil)
		return
	}
	r.d.OnEvent(e)
}

// get router config from dark launch, including file and governance
func getRouterConfigFromDarkLaunch() (*model.RouterConfig, error) {
	routeRules := &model.RouterConfig{
		Destinations: map[string][]*model.RouteRule{},
	}

	configMap := archaius.GetConfigs()
	//filter out key:value pairs which are not route rules
	for k := range configMap {
		if !strings.HasPrefix(k, DarkLaunchPrefix) {
			delete(configMap, k)
		}
	}

	//put route rules in configMap into routeRules
	rule := &model.DarkLaunchRule{}
	for k, v := range configMap {
		// todo bug fix
		value, ok := v.(string)
		if !ok {
			return routeRules, errors.New("route rule is not a json string format please check the configuration in config center")
		}
		if err := json.Unmarshal([]byte(value), rule); err != nil {
			return routeRules, err
		}
		key := strings.Replace(k, DarkLaunchPrefix, "", 1)
		routeRules.Destinations[key] = DarkLaunchRule2RouteRule(rule)
	}
	return routeRules, nil
}

// NewRouteDarkLaunchGovernSource returns default RouteDarkLaunchGovernSource
func NewRouteDarkLaunchGovernSource() *RouteDarkLaunchGovernSource {
	if source == nil {
		source = &RouteDarkLaunchGovernSource{
			callbackCheck: make(chan bool),
		}
	}
	return source
}
