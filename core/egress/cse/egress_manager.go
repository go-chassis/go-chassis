package cse

import (
	"errors"
	"sync"

	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-archaius/core/config-manager"
	"github.com/ServiceComb/go-archaius/core/event-system"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/egress"
)

const egressFileSourceName = "EgressFileSource"
const egressFileSourcePriority = 16

var egressRuleMgr core.ConfigMgr

type egressRuleEventListener struct{}

// update egress rule of a service
func (r *egressRuleEventListener) Event(e *core.Event) {
	if e == nil {
		lager.Logger.Warn("Event pointer is nil", nil)
		return
	}

	v := egressRuleMgr.GetConfigurationsByKey(e.Key)
	if v == nil {
		DeleteEgressRuleByKey(e.Key)
		lager.Logger.Infof("[%s] Egress rule is removed", e.Key)
		return
	}
	egressRules, ok := v.([]*model.EgressRule)
	if !ok {
		lager.Logger.Error("value is not type []*RouteRule", nil)
		return
	}

	 ok, _= egress.ValidateEgressRule(map[string][]*model.EgressRule{e.Key: egressRules})
	 if ok   {
		SetEgressRuleByKey(e.Key, egressRules)
		lager.Logger.Infof("Update [%s] route rule success", e.Key)
	}
}

// egressFileSource keeps the egress rule in egress file,
// after init, it's data does not change
type egressFileSource struct {
	once sync.Once
	d    map[string]interface{}
}

func newRouteFileSource() *egressFileSource {
	r := &egressFileSource{}
	r.once.Do(func() {
		routeRules := GetEgressRule()
		d := make(map[string]interface{}, 0)
		if routeRules == nil {
			r.d = d
			lager.Logger.Error("Can not get any router config", nil)
			return
		}
		for k, v := range routeRules {
			d[k] = v
		}
		r.d = d
	})
	return r
}

func (r *egressFileSource) GetSourceName() string {
	return egressFileSourceName
}
func (r *egressFileSource) GetConfigurations() (map[string]interface{}, error) {
	configMap := make(map[string]interface{})
	for k, v := range r.d {
		configMap[k] = v
	}
	return configMap, nil
}
func (r *egressFileSource) GetConfigurationsByDI(dimensionInfo string) (map[string]interface{}, error) {
	return nil, nil
}
func (r *egressFileSource) GetConfigurationByKey(k string) (interface{}, error) {
	v, ok := r.d[k]
	if !ok {
		return nil, errors.New("key " + k + " not exist")
	}
	return v, nil
}
func (r *egressFileSource) GetConfigurationByKeyAndDimensionInfo(key, dimensionInfo string) (interface{}, error) {
	return nil, nil
}
func (r *egressFileSource) AddDimensionInfo(dimensionInfo string) (map[string]string, error) {
	return nil, nil
}
func (r *egressFileSource) DynamicConfigHandler(core.DynamicConfigCallback) error {
	return nil
}
func (r *egressFileSource) GetPriority() int {
	return egressFileSourcePriority
}
func (r *egressFileSource) Cleanup() error { return nil }

// initialize the config mgr and add several sources
func initEgressManager() error {
	d := eventsystem.NewDispatcher()
	l := &egressRuleEventListener{}
	d.RegisterListener(l, ".*")
	egressRuleMgr = configmanager.NewConfigurationManager(d)
	if err := AddRouteRuleSource(newRouteFileSource()); err != nil {
		return err
	}
	return nil
}

// AddRouteRuleSource adds a config source to route rule manager
// Do not call this method until router init success
func AddRouteRuleSource(s core.ConfigSource) error {
	if s == nil {
		return errors.New("source nil")
	}
	if egressRuleMgr == nil {
		return errors.New("egressRuleMgr is nil, please init it firstly")
	}
	if err := egressRuleMgr.AddSource(s, s.GetPriority()); err != nil {
		return err
	}
	lager.Logger.Infof("Add [%s] source success", s.GetSourceName())
	return nil
}
