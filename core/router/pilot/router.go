package pilot

import (
	"fmt"
	"sync"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/router"
)

func init() { router.InstallRouterService("pilot", newPilotRouter) }

func newPilotRouter() (router.Router, error) { return &PilotRouter{}, nil }

//PilotRouter is pilot router service
type PilotRouter struct{}

//FetchRouteRule return all rules
func (r *PilotRouter) FetchRouteRule() map[string][]*model.RouteRule {
	return GetRouteRule()
}

//SetRouteRule set rules
func (r *PilotRouter) SetRouteRule(rr map[string][]*model.RouteRule) {
	SetRouteRule(rr)
}

//FetchRouteRuleByServiceName get rules for service
func (r *PilotRouter) FetchRouteRuleByServiceName(service string) []*model.RouteRule {
	return GetRouteRuleByKey(service)
}

//Init init router config
func (r *PilotRouter) Init(o router.Options) error {
	// SetDestinations router destinations
	for target := range config.GetRouterReference() {
		SetRouteRuleByKey(target, nil)
	}
	// the manager use dests to init, so must init after dests
	if err := InitPilotFetcher(o); err != nil {
		return err
	}
	return refresh()
}

// InitRouteRuleByKey init route rule by service key
func (r *PilotRouter) InitRouteRuleByKey(k string) {
	lock.RLock()
	_, ok1 := dests[k]
	lock.RUnlock()

	if !ok1 {
		lock.Lock()
		if _, ok2 := dests[k]; !ok2 {
			dests[k] = nil
			setChanForPilot(k)
		}
		lock.Unlock()
	}
}

// refresh all the router config
func refresh() error {
	configs := pilotfetcher.GetConfigurations()
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
	lock.Lock()
	defer lock.Unlock()
	dests = rule
}
