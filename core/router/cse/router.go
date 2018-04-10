package cse

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/router"
	"sync"
)

//Router is cse router service
type Router struct {
}

//FetchRouteRule return all rules
func (r *Router) FetchRouteRule() map[string][]*model.RouteRule {
	return GetRouteRule()
}

//SetRouteRule set rules
func (r *Router) SetRouteRule(rr map[string][]*model.RouteRule) {
	SetRouteRule(rr)
}

//FetchRouteRuleByServiceName get rules for service
func (r *Router) FetchRouteRuleByServiceName(service string) []*model.RouteRule {
	return GetRouteRuleByKey(service)
}
func newRouter() (router.Router, error) {
	// the manager use dests to init, so must init after dests
	if err := initRouterManager(); err != nil {
		return nil, err
	}

	if err := refresh(); err != nil {
		return nil, err
	}
	return &Router{}, nil
}

// refresh all the router config
func refresh() error {
	configs := routeRuleMgr.GetConfigurations()
	d := make(map[string][]*model.RouteRule)
	for k, v := range configs {
		rules, ok := v.([]*model.RouteRule)
		if !ok {
			err := fmt.Errorf("route rule type assertion fail, key: %s", k)
			return err
		}
		d[k] = rules
	}

	if router.ValidateRule(d) {
		dests = d
	}
	return nil
}

var dests = make(map[string][]*model.RouteRule)
var lock sync.RWMutex

// SetRouteRuleByKey set route rule by key
func SetRouteRuleByKey(k string, r []*model.RouteRule) {
	lock.Lock()
	dests[k] = r
	lock.Unlock()
}

// DeleteRouteRuleByKey set route rule by key
func DeleteRouteRuleByKey(k string) {
	lock.Lock()
	delete(dests, k)
	lock.Unlock()
}

// GetRouteRuleByKey get route rule by key
func GetRouteRuleByKey(k string) []*model.RouteRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests[k]
}

// GetRouteRule get route rule
func GetRouteRule() map[string][]*model.RouteRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests
}

// SetRouteRule set route rule
func SetRouteRule(rule map[string][]*model.RouteRule) {
	lock.RLock()
	defer lock.RUnlock()
	dests = rule
}
func init() {
	router.InstallRouterService("cse", newRouter)
}
