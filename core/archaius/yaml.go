package archaius

import (
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/lager"
	"os"
)

// PathExist to check the existence of the file path
func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// ConfigListener is provides the object to listen the events
type ConfigListener struct{}

// Event is for to receive the events based on registered key and object pairs
func (cl *ConfigListener) Event(e *core.Event) {
	lager.Logger.Infof("the value of %s is change to %v, get the values: %v", e.Key, e.Value, Get(e.Key))
}
