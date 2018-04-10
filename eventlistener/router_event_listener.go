package eventlistener

import (
	"strings"

	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/router/cse"
)

// constants for dark launch key and prefix
const (
	//DarkLaunchKey & DarkLaunchPrefix is a variable of type string
	DarkLaunchKey    = "^cse\\.darklaunch\\.policy\\."
	DarkLaunchPrefix = "cse.darklaunch.policy."
)

//DarkLaunchEventListener is a struct
type DarkLaunchEventListener struct{}

//Event is method used for dark launch event listening
func (d *DarkLaunchEventListener) Event(event *core.Event) {
	lager.Logger.Debugf("Get darkLaunch event, key: %s, type: %s", event.Key, event.EventType)
	s := strings.Replace(event.Key, DarkLaunchPrefix, "", 1)
	newEvent := &core.Event{
		EventSource: cse.RouteDarkLaunchGovernSourceName,
		EventType:   event.EventType,
		Key:         s,
		Value:       event.Value,
	}
	cse.NewRouteDarkLaunchGovernSource().Callback(newEvent)
}
