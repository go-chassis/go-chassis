package cse

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"sync"
	"github.com/ServiceComb/go-chassis/core/egress"
)

//Egress is cse Egress service
type Egress struct {
}

//FetchEgressRule return all rules
func (r *Egress) FetchEgressRule() map[string][]*model.EgressRule {
	return GetEgressRule()
}

//SetEgressRule set rules
func (r *Egress) SetEgressRule(rr map[string][]*model.EgressRule) {
	SetEgressRule(rr)
}

//FetchEgressRuleByName get rules by name
func (r *Egress) FetchEgressRuleByName(name string) []*model.EgressRule {
	return GetEgressRuleByKey(name)
}

//Init init router config
func (r *Egress) Init() error {
	// the manager use dests to init, so must init after dests
	if err := initEgressManager(); err != nil {
		return err
	}
	return refresh()
}


func newEgress() (egress.Egress, error) {
	return &Egress{}, nil
}

// refresh all the router config
func refresh() error {
	configs := egressRuleMgr.GetConfigurations()
	d := make(map[string][]*model.EgressRule)
	for k, v := range configs {
		rules, ok := v.([]*model.EgressRule)
		if !ok {
			err := fmt.Errorf("Egress rule type assertion fail, key: %s", k)
			return err
		}
		d[k] = rules
	}

	ok, _:= egress.ValidateEgressRule(d)
	if ok {
		dests = d
	}
	return nil
}

var dests = make(map[string][]*model.EgressRule)
var lock sync.RWMutex

// SetEgressRuleByKey set route rule by key
func SetEgressRuleByKey(k string, r []*model.EgressRule) {
	lock.Lock()
	dests[k] = r
	lock.Unlock()
}

// DeleteEgressRuleByKey set route rule by key
func DeleteEgressRuleByKey(k string) {
	lock.Lock()
	delete(dests, k)
	lock.Unlock()
}

// GetEgressRuleByKey get route rule by key
func GetEgressRuleByKey(k string) []*model.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests[k]
}

// GetEgressRule get route rule
func GetEgressRule() map[string][]*model.EgressRule {
	lock.RLock()
	defer lock.RUnlock()
	return dests
}

// SetEgressRule set route rule
func SetEgressRule(rule map[string][]*model.EgressRule) {
	lock.RLock()
	defer lock.RUnlock()
	dests = rule
}
func init() {
	egress.InstallEgressService("cse", newEgress)
}

