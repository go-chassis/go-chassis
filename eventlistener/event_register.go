package eventlistener

import (
	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-chassis/core/archaius"
)

//RegisterKeys registers a config key to the archaius
func RegisterKeys(eventListener core.EventListener, keys ...string) {

	archaius.RegisterListener(eventListener, keys...)
}

//Init is a function
func Init() {
	qpsEventListener := &QPSEventListener{}
	circuitBreakerEventListener := &CircuitBreakerEventListener{}
	lbEventListener := &LoadbalanceEventListener{}

	RegisterKeys(qpsEventListener, QPSLimitKey)
	RegisterKeys(circuitBreakerEventListener, ConsumerFallbackKey, ConsumerFallbackPolicyKey, ConsumerIsolationKey, ConsumerCircuitbreakerKey)
	RegisterKeys(lbEventListener, LoadBalanceKey)
	RegisterKeys(&DarkLaunchEventListener{}, DarkLaunchKey)

}
