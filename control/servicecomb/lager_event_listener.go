package servicecomb

import (
	"github.com/go-chassis/go-archaius/event"
	"strings"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/openlog"
	"github.com/go-chassis/seclog"
	"github.com/go-chassis/seclog/third_party/forked/cloudfoundry/lager"
)

const (
	//LagerLevelKey is a variable of type string
	LagerLevelKey = "logLevel"
)

//LagerEventListener is a struct used for Event listener
type LagerEventListener struct {
	//Key []string
	Key string
}

//Event is a method for Lager event listening
func (el *LagerEventListener) Event(e *event.Event) {
	logger := openlog.GetLogger()
	l, ok := logger.(lager.Logger)
	if !ok {
		return
	}

	openlog.Info("Get lager e", openlog.WithTags(openlog.Tags{
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
	case seclog.DEBUG:
		lagerLogLevel = lager.DEBUG
	case seclog.INFO:
		lagerLogLevel = lager.INFO
	case seclog.WARN:
		lagerLogLevel = lager.WARN
	case seclog.ERROR:
		lagerLogLevel = lager.ERROR
	case seclog.FATAL:
		lagerLogLevel = lager.FATAL
	default:
		openlog.Info("ops..., got unknown logger level")
		return
	}

	switch e.EventType {
	case common.Update:
		l.SetLogLevel(lagerLogLevel)
	}
}
