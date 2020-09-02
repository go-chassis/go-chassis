package eventlistener_test

import (
	"github.com/go-chassis/go-archaius/event"
	"testing"

	"github.com/go-chassis/go-chassis/v2/eventlistener"
)

func TestLagerEventListener_Event(t *testing.T) {
	t.Log("Test lager_event_listener_test.go")
	eventlistener.Init()
	eventListen := &eventlistener.LagerEventListener{}
	t.Log("sending the events for the key logger_level")

	e1 := &event.Event{EventType: "UPDATE", Key: "logger_level", Value: "INFO"}
	eventListen.Event(e1)

	e2 := &event.Event{EventType: "UPDATE", Key: "logger_level", Value: "WARN"}
	eventListen.Event(e2)

	e3 := &event.Event{EventType: "UPDATE", Key: "logger_level", Value: "ERROR"}
	eventListen.Event(e3)

	e4 := &event.Event{EventType: "UPDATE", Key: "logger_level", Value: "FATAL"}
	eventListen.Event(e4)

	e5 := &event.Event{EventType: "UPDATE", Key: "logger_level", Value: "BAD"}
	eventListen.Event(e5)

	e := &event.Event{EventType: "UPDATE", Key: "logger_level", Value: "DEBUG"}
	eventListen.Event(e)

}
