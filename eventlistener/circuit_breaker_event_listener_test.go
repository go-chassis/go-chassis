package eventlistener_test

import (
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/core/lager"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/eventlistener"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"

	"github.com/stretchr/testify/assert"
)

func preTest() {
	os.Setenv(fileutil.ChassisHome,
		filepath.Join(os.Getenv("GOPATH"),
			"src",
			"github.com",
			"go-chassis",
			"go-chassis",
			"examples",
			"discovery",
			"server"))
}

func TestCircuitBreakerEventListener_Event(t *testing.T) {
	preTest()
	config.Init()
	t.Log("Test circuit_breaker_event_listener.go")
	eventlistener.Init()
	eventListen := &eventlistener.CircuitBreakerEventListener{}
	t.Log("sending the events for the key cse.flowcontrol.Consumer.qps.limit.Server")
	e := &event.Event{EventType: "UPDATE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e)

	e1 := &event.Event{EventType: "CREATE", Key: "cse.flowcontrol.Provider.qps.limit.Server", Value: 100}
	eventListen.Event(e1)

	e2 := &event.Event{EventType: "DELETE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e2)

}
func TestGetNames(t *testing.T) {
	preTest()
	config.Init()
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
func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
