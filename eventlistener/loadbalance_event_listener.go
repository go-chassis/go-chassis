package eventlistener

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"

	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/config"
)

// constants for loadbalance strategy name, and timeout
const (
	//LoadBalanceKey is variable of type string that matches load balancing events
	LoadBalanceKey = "^cse\\.loadbalance\\."
	Update         = "update"
	Delete         = "delete"
)

//LoadbalanceEventListener is a struct
type LoadbalanceEventListener struct {
	Key string
}

//Event is a method used to handle a load balancing event
func (e *LoadbalanceEventListener) Event(event *core.Event) {
	lager.Logger.Debugf("LB event, key: %s, type: %s", event.Key, event.EventType)
	if err := config.ReadLBFromArchaius(); err != nil {
		lager.Logger.Error("can not unmarshal new lb config", err)
	}
	if event.Key == "cse.loadbalance.strategy.name" {
		switch event.EventType {
		case Update:
			strategyName := event.Value
			strategy, err := loadbalance.GetStrategyPlugin(strategyName.(string))
			if err != nil {
				lager.Logger.Errorf(err, "Get strategy ["+strategyName.(string)+"] failed")
			} else {

				o := loadbalance.DefaultSelector.Options()
				o.Strategy = strategy
			}
		case Delete:
			strategyName := "RoundRobin"
			strategy, err := loadbalance.GetStrategyPlugin(strategyName)
			if err != nil {
				lager.Logger.Errorf(err, "Get strategy ["+strategyName+"] failed")
			} else {

				o := loadbalance.DefaultSelector.Options()
				o.Strategy = strategy
			}
		}
	}
}
