package eventlistener_test

import (
	"testing"

	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/eventlistener"
)

func TestLagerEventListener_Event(t *testing.T) {
	preTest()
	config.Init()
	t.Log("Test lager_event_listener_test.go")
	eventlistener.Init()
	eventListen := &eventlistener.LagerEventListener{}
	t.Log("sending the events for the key logger_level")

	e1 := &core.Event{EventType: "UPDATE", Key: "logger_level", Value: "INFO"}
	eventListen.Event(e1)

	e2 := &core.Event{EventType: "UPDATE", Key: "logger_level", Value: "WARN"}
	eventListen.Event(e2)

	e3 := &core.Event{EventType: "UPDATE", Key: "logger_level", Value: "ERROR"}
	eventListen.Event(e3)

	e4 := &core.Event{EventType: "UPDATE", Key: "logger_level", Value: "FATAL"}
	eventListen.Event(e4)

	e5 := &core.Event{EventType: "UPDATE", Key: "logger_level", Value: "BAD"}
	eventListen.Event(e5)

	e := &core.Event{EventType: "UPDATE", Key: "logger_level", Value: "DEBUG"}
	eventListen.Event(e)

}
