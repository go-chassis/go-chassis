package servicecomb

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/openlog"
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
	err := archaius.RegisterListener(&routeRuleEventListener{}, DarkLaunchKey, DarkLaunchKeyV2)
	if err != nil {
		openlog.Error(err.Error())
	}
	return r.LoadRules()
}

func newRouter() (router.Router, error) {
	cseRouter = &Router{
		routeRule: make(map[string][]*config.RouteRule),
		lock:      sync.RWMutex{},
	}
	return cseRouter, nil
}

// LoadRules load all the router config
func (r *Router) LoadRules() error {
	configs, err := MergeLocalAndRemoteConfig()
	if err != nil {
		openlog.Error("init route rule failed", openlog.WithTags(openlog.Tags{
			"err": err.Error(),
		}))
	}

	if router.ValidateRule(configs) {
		r.routeRule = configs
		openlog.Info("load route rule", openlog.WithTags(openlog.Tags{
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
	openlog.Info("update route rule success", openlog.WithTags(
		openlog.Tags{
			"service": k,
			"rule":    rr,
		}))
}

// DeleteRouteRuleByKey set route rule by key
func (r *Router) DeleteRouteRuleByKey(k string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.routeRule, k)
	openlog.Info("route rule is removed", openlog.WithTags(
		openlog.Tags{
			"service": k,
		}))
}

func init() {
	router.InstallRouterPlugin("cse", newRouter)
}
