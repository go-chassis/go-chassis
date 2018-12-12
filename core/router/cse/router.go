package cse

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/router"
	"sync"
)

//Router is cse router service
type Router struct {
}

//SetRouteRule set rules
func (r *Router) SetRouteRule(rr map[string][]*model.RouteRule) {
	lock.Lock()
	defer lock.Unlock()
	dests = rr
}

//FetchRouteRuleByServiceName get rules for service
func (r *Router) FetchRouteRuleByServiceName(service string) []*model.RouteRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests[service]
}

//Init init router config
func (r *Router) Init(o router.Options) error {
	// the manager use dests to init, so must init after dests
	if err := initRouterManager(); err != nil {
		return err
	}
	return refresh()
}

// InitRouteRuleByKey init route rule by service key
func (r *Router) InitRouteRuleByKey(k string) {}

func newRouter() (router.Router, error) {
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

func init() {
	router.InstallRouterService("cse", newRouter)
}
