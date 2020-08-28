package eventlistener

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/openlog"
)

//RegisterKeys registers a config key to the archaius
func RegisterKeys(eventListener event.Listener, keys ...string) {
	err := archaius.RegisterListener(eventListener, keys...)
	if err != nil {
		openlog.Error(err.Error())
	}
}

//Init is a function
func Init() {
	qpsEventListener := &QPSEventListener{}
	circuitBreakerEventListener := &CircuitBreakerEventListener{}
	lbEventListener := &LoadbalancingEventListener{}

	RegisterKeys(qpsEventListener, QPSLimitKey)
	RegisterKeys(circuitBreakerEventListener, ConsumerFallbackKey, ConsumerFallbackPolicyKey, ConsumerIsolationKey, ConsumerCircuitbreakerKey)
	RegisterKeys(lbEventListener, LoadBalanceKey)
	RegisterKeys(&LagerEventListener{}, LagerLevelKey)

}
