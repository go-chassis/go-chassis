package cse

import (
	"errors"
	"sync"

	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-archaius/core/config-manager"
	"github.com/go-chassis/go-archaius/core/event-system"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-mesh/openlogging"
)

const routeFileSourceName = "RouteFileSource"
const routeFileSourcePriority = 16

var routeRuleMgr core.ConfigMgr

// routeFileSource keeps the route rule in router file,
// after init, it's data does not change
type routeFileSource struct {
	once sync.Once
	d    map[string]interface{}
}

func newRouteFileSource() *routeFileSource {
	r := &routeFileSource{}
	r.once.Do(func() {
		routeRules := dests
		d := make(map[string]interface{}, 0)
		if routeRules == nil {
			r.d = d
			openlogging.Warn("can not get router config in local")
			return
		}
		for k, v := range routeRules {
			d[k] = v
		}
		r.d = d
	})
	return r
}

func (r *routeFileSource) GetSourceName() string {
	return routeFileSourceName
}
func (r *routeFileSource) GetConfigurations() (map[string]interface{}, error) {
	configMap := make(map[string]interface{})
	for k, v := range r.d {
		configMap[k] = v
	}
	return configMap, nil
}
func (r *routeFileSource) GetConfigurationsByDI(dimensionInfo string) (map[string]interface{}, error) {
	return nil, nil
}
func (r *routeFileSource) GetConfigurationByKey(k string) (interface{}, error) {
	v, ok := r.d[k]
	if !ok {
		return nil, errors.New("key " + k + " not exist")
	}
	return v, nil
}
func (r *routeFileSource) GetConfigurationByKeyAndDimensionInfo(key, dimensionInfo string) (interface{}, error) {
	return nil, nil
}
func (r *routeFileSource) AddDimensionInfo(dimensionInfo string) (map[string]string, error) {
	return nil, nil
}
func (r *routeFileSource) DynamicConfigHandler(core.DynamicConfigCallback) error {
	return nil
}
func (r *routeFileSource) GetPriority() int {
	return routeFileSourcePriority
}
func (r *routeFileSource) SetPriority(int) {
}
func (r *routeFileSource) Cleanup() error { return nil }

// initialize the config mgr and add several sources
func initRouterManager() error {
	d := eventsystem.NewDispatcher()
	l := &routeRuleEventListener{}
	d.RegisterListener(l, ".*")
	routeRuleMgr = configmanager.NewConfigurationManager(d)
	if err := AddRouteRuleSource(newRouteFileSource()); err != nil {
		return err
	}
	return AddRouteRuleSource(NewRouteDarkLaunchGovernSource())
}

// AddRouteRuleSource adds a config source to route rule manager
// Do not call this method until router init success
func AddRouteRuleSource(s core.ConfigSource) error {
	if s == nil {
		return errors.New("source nil")
	}
	if routeRuleMgr == nil {
		return errors.New("routeRuleMgr is nil, please init it firstly")
	}
	if err := routeRuleMgr.AddSource(s, s.GetPriority()); err != nil {
		return err
	}

	lager.Logger.Infof("add [%s] source success", s.GetSourceName())
	return nil
}
