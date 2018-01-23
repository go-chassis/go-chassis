package eventlistener_test

import (
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/eventlistener"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCircuitBreakerEventListener_Event(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")

	t.Log("Test circuit_breaker_event_listener.go")
	config.Init()
	lager.Initialize("", "INFO", "l", "size", true, 1, 10, 7)
	eventlistener.Init()
	eventListen := &eventlistener.CircuitBreakerEventListener{}
	t.Log("sending the events for the key cse.flowcontrol.Consumer.qps.limit.Server")
	e := &core.Event{EventType: "UPDATE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e)

	e1 := &core.Event{EventType: "CREATE", Key: "cse.flowcontrol.Provider.qps.limit.Server", Value: 100}
	eventListen.Event(e1)

	e2 := &core.Event{EventType: "DELETE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e2)

}
func TestGetNames(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	t.Log("verifying configuration keys by GetNames method")
	sourceName, serviceName := eventlistener.GetNames("cse.isolation.Web.Consumer.carts.timeout.enabled")
	assert.Equal(t, "Web", sourceName)
	assert.Equal(t, "carts", serviceName)
	n := eventlistener.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Web.Consumer.carts", n)

	sourceName, serviceName = eventlistener.GetNames("cse.isolation.Consumer.carts.timeout.enabled")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "carts", serviceName)
	n = eventlistener.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer.carts", n)

	sourceName, serviceName = eventlistener.GetNames("cse.circuitBreaker.Consumer.forceOpen")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "", serviceName)
	n = eventlistener.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer", n)

	sourceName, serviceName = eventlistener.GetNames("cse.isolation.Consumer.carts.interface.get.timeout.enabled")
	assert.Equal(t, "", sourceName)
	assert.Equal(t, "carts.interface.get", serviceName)
	n = eventlistener.GetCircuitName(sourceName, serviceName)
	assert.Equal(t, "Consumer.carts.interface.get", n)

}
