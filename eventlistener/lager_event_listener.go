package eventlistener

import (
	"fmt"
	"github.com/go-chassis/go-archaius/core"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/paas-lager/third_party/forked/cloudfoundry/lager"
	"github.com/go-mesh/openlogging"
	"github.com/go-chassis/paas-lager"
	"strings"
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
func (e *LagerEventListener) Event(event *core.Event) {
	defer func() {
		if err := recover(); err != nil {
			openlogging.GetLogger().Errorf("%s", err)
		}
	}()
	logger := openlogging.GetLogger()
	logger.Debugf("Get lager event, key: %s, type: %s", event.Key, event.EventType)
	l, ok := logger.(lager.Logger)
	if !ok {
		return
	}
	var lagerLogLevel lager.LogLevel

	Value, ok := event.Value.(string)
	if !ok {
		openlogging.GetLogger().Errorf("event.Value Assertion err:%s", event.
			Value)
		return
	}

	switch strings.ToUpper(Value) {
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
		fmt.Printf("unknown logger level: %s", event.Value)
		return
	}

	switch event.EventType {
	case common.Update:
		l.SetLogLevel(lagerLogLevel)
	}
}
