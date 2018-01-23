package eventlistener

import (
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
)

// constants for loadbalance strategy name, and timeout
const (
	//LbStrategyNameKey & LbStrategyTimeoutKey are variables of type string
	LbStrategyNameKey    = "cse.loadbalance.strategy.name"
	LbStrategyTimeoutKey = "cse.loadbalance.strategy.sessionTimeoutInSeconds"
	Update               = "update"
	Delete               = "delete"
)

//LoadbalanceEventListener is a struct
type LoadbalanceEventListener struct {
	Key string
}

//Event is a method used to handle a load balancing event
func (e *LoadbalanceEventListener) Event(event *core.Event) {

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
