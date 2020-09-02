package eventlistener_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"testing"

	"github.com/go-chassis/go-chassis/v2/eventlistener"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.loadbalance.strategy.name", "SessionStickiness")
	config.ReadHystrixFromArchaius()
}
func TestCircuitBreakerEventListener_Event(t *testing.T) {
	t.Log("Test circuit_breaker_event_listener.go")
	eventlistener.Init()
	eventListen := &eventlistener.CircuitBreakerEventListener{}
	t.Log("sending the events for the key servicecomb.flowcontrol.Consumer.qps.limit.Server")
	e := &event.Event{EventType: "UPDATE", Key: "servicecomb.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e)

	e1 := &event.Event{EventType: "CREATE", Key: "servicecomb.flowcontrol.Provider.qps.limit.Server", Value: 100}
	eventListen.Event(e1)

	e2 := &event.Event{EventType: "DELETE", Key: "servicecomb.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e2)

}
func TestGetNames(t *testing.T) {
	t.Log("verifying configuration keys by GetNames method")
	sourceName, serviceName := eventlistener.GetNames("servicecomb.isolation.Web.Consumer.carts.timeout.enabled")
	assert.Equal(t, "Web", sourceName)
	assert.Equal(t, "carts", serviceName)
	n := eventlistener.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Web.Consumer.carts", n)

	sourceName, serviceName = eventlistener.GetNames("servicecomb.isolation.Consumer.carts.timeout.enabled")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "carts", serviceName)
	n = eventlistener.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer.carts", n)

	sourceName, serviceName = eventlistener.GetNames("servicecomb.circuitBreaker.Consumer.forceOpen")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "", serviceName)
	n = eventlistener.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer", n)

	sourceName, serviceName = eventlistener.GetNames("servicecomb.isolation.Consumer.carts.interface.get.timeout.enabled")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "carts.interface.get", serviceName)
	n = eventlistener.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer.carts.interface.get", n)
}
