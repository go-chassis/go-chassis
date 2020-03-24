package servicecomb

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-mesh/openlogging"
	"sync"
)

var cseRouter *Router

//Router is cse router service
type Router struct {
	routeRule map[string][]*config.RouteRule
	lock      sync.RWMutex
}

//SetRouteRule set rules
func (r *Router) SetRouteRule(rr map[string][]*config.RouteRule) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.routeRule = rr
}

//FetchRouteRuleByServiceName get rules for service
func (r *Router) FetchRouteRuleByServiceName(service string) []*config.RouteRule {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.routeRule[service]
}

//ListRouteRule get rules for all service
func (r *Router) ListRouteRule() map[string][]*config.RouteRule {
	r.lock.RLock()
	defer r.lock.RUnlock()
	rr := make(map[string][]*config.RouteRule, len(r.routeRule))
	for k, v := range r.routeRule {
		rr[k] = v
	}
	return rr
}

//Init init router config
func (r *Router) Init(o router.Options) error {
	archaius.RegisterListener(&routeRuleEventListener{}, DarkLaunchKey, DarkLaunchKeyV2)
	return r.LoadRules()
}

func newRouter() (router.Router, error) {
	cseRouter = &Router{
		routeRule: make(map[string][]*config.RouteRule, 0),
		lock:      sync.RWMutex{},
	}
	return cseRouter, nil
}

// LoadRules load all the router config
func (r *Router) LoadRules() error {
	configs, err := MergeLocalAndRemoteConfig()
	if err != nil {
		openlogging.Error("init route rule failed", openlogging.WithTags(openlogging.Tags{
			"err": err.Error(),
		}))
	}

	if router.ValidateRule(configs) {
		r.routeRule = configs
		openlogging.Info("load route rule", openlogging.WithTags(openlogging.Tags{
			"rule": r.routeRule,
		}))
	}
	return nil
}

// SetRouteRuleByKey set route rule by key
func (r *Router) SetRouteRuleByKey(k string, rr []*config.RouteRule) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.routeRule[k] = rr
	openlogging.Info("update route rule success", openlogging.WithTags(
		openlogging.Tags{
			"service": k,
			"rule":    rr,
		}))
}

// DeleteRouteRuleByKey set route rule by key
func (r *Router) DeleteRouteRuleByKey(k string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.routeRule, k)
	openlogging.Info("route rule is removed", openlogging.WithTags(
		openlogging.Tags{
			"service": k,
		}))
}

func init() {
	router.InstallRouterService("cse", newRouter)
}
