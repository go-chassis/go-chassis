package eventlistener_test

import (
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/eventlistener"
	"os"
	"testing"
)

func TestQpsEvent(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")

	t.Log("Test qps_event_listener.go")
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
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
