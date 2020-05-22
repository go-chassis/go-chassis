package eventlistener_test

import (
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/eventlistener"
	"testing"
)

func TestQpsEvent(t *testing.T) {
	eventlistener.Init()
	eventListen := &eventlistener.QPSEventListener{}
	t.Log("sending the events for the key cse.flowcontrol.Consumer.qps.limit.Server")
	e := &event.Event{EventType: "UPDATE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e)

	e1 := &event.Event{EventType: "CREATE", Key: "cse.flowcontrol.Provider.qps.limit.Server", Value: 100}
	eventListen.Event(e1)

	e2 := &event.Event{EventType: "DELETE", Key: "cse.flowcontrol.Consumer.qps.limit.Server", Value: 199}
	eventListen.Event(e2)

}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
