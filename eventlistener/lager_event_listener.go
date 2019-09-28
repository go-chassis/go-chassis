package eventlistener

import (
	"github.com/go-chassis/go-archaius/event"
	"strings"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/paas-lager"
	"github.com/go-chassis/paas-lager/third_party/forked/cloudfoundry/lager"
	"github.com/go-mesh/openlogging"
)

const (
	//LagerLevelKey is a variable of type string
	LagerLevelKey = "logger_level"
)

//LagerEventListener is a struct used for Event listener
type LagerEventListener struct {
	//Key []string
	Key string
}

//Event is a method for Lager event listening
func (el *LagerEventListener) Event(e *event.Event) {
	logger := openlogging.GetLogger()
	l, ok := logger.(lager.Logger)
	if !ok {
		return
	}

	openlogging.Info("Get lager e", openlogging.WithTags(openlogging.Tags{
		"key":   e.Key,
		"value": e.Value,
		"type":  e.EventType,
	}))

	v, ok := e.Value.(string)
	if !ok {
		return
	}

	var lagerLogLevel lager.LogLevel
	switch strings.ToUpper(v) {
	case log.DEBUG:
		lagerLogLevel = lager.DEBUG
	case log.INFO:
		lagerLogLevel = lager.INFO
	case log.WARN:
		lagerLogLevel = lager.WARN
	case log.ERROR:
		lagerLogLevel = lager.ERROR
	case log.FATAL:
		lagerLogLevel = lager.FATAL
	default:
		openlogging.Info("ops..., got unknown logger level")
		return
	}

	switch e.EventType {
	case common.Update:
		l.SetLogLevel(lagerLogLevel)
	}
}
