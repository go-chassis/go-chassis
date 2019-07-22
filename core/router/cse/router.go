package cse

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-mesh/openlogging"
	"sync"
)

var cseRouter *Router

//Router is cse router service
type Router struct {
	dests map[string][]*model.RouteRule
	lock  sync.RWMutex
}

//SetRouteRule set rules
func (r *Router) SetRouteRule(rr map[string][]*model.RouteRule) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.dests = rr
}

//FetchRouteRuleByServiceName get rules for service
func (r *Router) FetchRouteRuleByServiceName(service string) []*model.RouteRule {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.dests[service]
}

//Init init router config
func (r *Router) Init(o router.Options) error {
	archaius.RegisterListener(&routeRuleEventListener{}, DarkLaunchKey)
	return r.LoadRules()
}

func newRouter() (router.Router, error) {
	cseRouter = &Router{
		dests: make(map[string][]*model.RouteRule, 0),
		lock:  sync.RWMutex{},
	}
	return cseRouter, nil
}

// LoadRules load all the router config
func (r *Router) LoadRules() error {
	configs, err := GetRouterRuleFromArchaius()
	if err != nil {
		openlogging.Error("init route rule failed", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
	}

	if router.ValidateRule(configs) {
		r.dests = configs
	}
	return nil
}

// SetRouteRuleByKey set route rule by key
func (r *Router) SetRouteRuleByKey(k string, rr []*model.RouteRule) {
	r.lock.Lock()
	r.dests[k] = rr
	r.lock.Unlock()
}

// DeleteRouteRuleByKey set route rule by key
func (r *Router) DeleteRouteRuleByKey(k string) {
	r.lock.Lock()
	delete(r.dests, k)
	r.lock.Unlock()
}

func init() {
	router.InstallRouterService("cse", newRouter)
}
