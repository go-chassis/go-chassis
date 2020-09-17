package servicecomb_test

import (
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/v2/control/servicecomb"
	"testing"
)

func TestLagerEventListener_Event(t *testing.T) {
	t.Log("Test lager_event_listener_test.go")
	servicecomb.Init()
	eventListen := &servicecomb.LagerEventListener{}
	t.Log("sending the events for the key logLevel")

	e1 := &event.Event{EventType: "UPDATE", Key: "logLevel", Value: "INFO"}
	eventListen.Event(e1)

	e2 := &event.Event{EventType: "UPDATE", Key: "logLevel", Value: "WARN"}
	eventListen.Event(e2)

	e3 := &event.Event{EventType: "UPDATE", Key: "logLevel", Value: "ERROR"}
	eventListen.Event(e3)

	e4 := &event.Event{EventType: "UPDATE", Key: "logLevel", Value: "FATAL"}
	eventListen.Event(e4)

	e5 := &event.Event{EventType: "UPDATE", Key: "logLevel", Value: "BAD"}
	eventListen.Event(e5)

	e := &event.Event{EventType: "UPDATE", Key: "logLevel", Value: "DEBUG"}
	eventListen.Event(e)

}
