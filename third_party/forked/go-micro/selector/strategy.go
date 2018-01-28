package selector

import (
	"math/rand"
	"sync"
	"time"

	"fmt"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
)

var strategies map[string]func([]*registry.MicroServiceInstance, interface{}) Next = make(map[string]func([]*registry.MicroServiceInstance, interface{}) Next)
var i int

func init() {
	rand.Seed(time.Now().UnixNano())
	i = rand.Int()
}

// InstallStrategy install strategy
func InstallStrategy(name string, strategy func([]*registry.MicroServiceInstance, interface{}) Next) {
	strategies[name] = strategy
	lager.Logger.Debugf("Installed strategy plugin: %s.", name)
}

// GetStrategyPlugin get strategy plugin
func GetStrategyPlugin(name string) (func([]*registry.MicroServiceInstance, interface{}) Next, error) {
	s, ok := strategies[name]
	if !ok {
		return nil, fmt.Errorf("Don't support strategyName [%s]", name)
	}

	return s, nil
}

// Random is a random strategy algorithm for node selection
func Random(instances []*registry.MicroServiceInstance, metadata interface{}) Next {
	return func() (*registry.MicroServiceInstance, error) {
		if len(instances) == 0 {
			return nil, ErrNoneAvailable
		}

		k := rand.Int() % len(instances)
		return instances[k], nil
	}
}

// RoundRobin is a roundrobin strategy algorithm for node selection
func RoundRobin(instances []*registry.MicroServiceInstance, metadata interface{}) Next {
	var mtx sync.Mutex
	return func() (*registry.MicroServiceInstance, error) {
		if len(instances) == 0 {
			return nil, ErrNoneAvailable
		}

		mtx.Lock()
		node := instances[i%len(instances)]
		i++
		mtx.Unlock()

		return node, nil
	}
}
