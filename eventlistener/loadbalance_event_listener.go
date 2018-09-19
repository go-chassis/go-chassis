package eventlistener

import (
	"github.com/go-chassis/go-chassis/core/lager"

	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-chassis/control/archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"regexp"
)

// constants for loadbalancer strategy name, and timeout
const (
	//LoadBalanceKey is variable of type string that matches load balancing events
	LoadBalanceKey          = "^cse\\.loadbalance\\."
	regex4normalloadbalance = "^cse\\.loadbalance\\.(strategy|SessionStickinessRule|retryEnabled|retryOnNext|retryOnSame|backoff)"
)

//LoadbalanceEventListener is a struct
type LoadbalanceEventListener struct {
	Key string
}

//Event is a method used to handle a load balancing event
func (e *LoadbalanceEventListener) Event(event *core.Event) {
	lager.Logger.Debugf("LB event, key: %s, type: %s", event.Key, event.EventType)
	if err := config.ReadLBFromArchaius(); err != nil {
		lager.Logger.Error("can not unmarshal new lb config: " + err.Error())
	}
	regNormal := regexp.MustCompile(regex4normalloadbalance)
	archaius.SaveToLBCache(config.GetLoadBalancing(), event.Key, regNormal.MatchString(event.Key))
}
