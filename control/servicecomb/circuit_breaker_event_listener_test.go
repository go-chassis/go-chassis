package servicecomb_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/v2/control/servicecomb"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("cse.loadbalance.strategy.name", "SessionStickiness")
	config.ReadHystrixFromArchaius()
}
func TestCircuitBreakerEventListener_Event(t *testing.T) {
	t.Log("Test circuit_breaker_event_listener.go")
	servicecomb.Init()
	eventListen := &servicecomb.CircuitBreakerEventListener{}
	t.Log("sending the events for the key cse.flowcontrol.Consumer.qps.limit.Server")
	e := &event.Event{EventType: "UPDATE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e)

	e1 := &event.Event{EventType: "CREATE", Key: "cse.flowcontrol.Provider.qps.limit.Server", Value: 100}
	eventListen.Event(e1)

	e2 := &event.Event{EventType: "DELETE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e2)

}
func TestGetNames(t *testing.T) {
	t.Log("verifying configuration keys by GetNames method")

	sourceName, serviceName := servicecomb.GetNames("cse.isolation.Consumer.carts.timeout.enabled")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "carts", serviceName)
	n := servicecomb.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer.carts", n)

	sourceName, serviceName = servicecomb.GetNames("cse.circuitBreaker.Consumer.forceOpen")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "", serviceName)
	n = servicecomb.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer", n)

	sourceName, serviceName = servicecomb.GetNames("cse.isolation.Consumer.carts.interface.get.timeout.enabled")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "carts.interface.get", serviceName)
	n = servicecomb.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer.carts.interface.get", n)
}
