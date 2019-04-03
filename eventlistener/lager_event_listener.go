package eventlistener

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/paas-lager"
	"strings"
	"github.com/go-chassis/go-archaius/core"
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
func (e *LagerEventListener) Event(event *core.Event) {
	logger := openlogging.GetLogger()
	logger.Debugf("Get lager event, key: %s, type: %s", event.Key, event.EventType)
	l, ok := logger.(lager.Logger)
	if !ok {
		return
	}

	var lagerLogLevel lager.LogLevel
	switch strings.ToUpper(event.Value.(string)) {
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
