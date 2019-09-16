package eventlistener_test

import (
	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/eventlistener"
	"os"
	"testing"
)

func TestQpsEvent(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")

	t.Log("Test qps_event_listener.go")
	config.Init()
	eventlistener.Init()
	eventListen := &eventlistener.QPSEventListener{}
	t.Log("sending the events for the key cse.flowcontrol.Consumer.qps.limit.Server")
	e := &core.Event{EventType: "UPDATE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e)

	e1 := &core.Event{EventType: "CREATE", Key: "cse.flowcontrol.Provider.qps.limit.Server", Value: 100}
	eventListen.Event(e1)

	e2 := &core.Event{EventType: "DELETE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e2)

}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
