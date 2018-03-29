package eventlistener

import (
	"github.com/ServiceComb/go-chassis/core/lager"

	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/config"
)

// constants for loadbalancer strategy name, and timeout
const (
	//LoadBalanceKey is variable of type string that matches load balancing events
	LoadBalanceKey = "^cse\\.loadbalance\\."
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
}
