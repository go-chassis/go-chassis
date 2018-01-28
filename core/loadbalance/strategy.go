package loadbalance

import (
	"math/rand"
	"time"

	"fmt"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
)

var strategies map[string]func([]*registry.MicroServiceInstance, interface{}) selector.Next = make(map[string]func([]*registry.MicroServiceInstance, interface{}) selector.Next)
var i int

func init() {
	rand.Seed(time.Now().UnixNano())
	i = rand.Int()
}

// InstallStrategy install strategy
func InstallStrategy(name string, strategy func([]*registry.MicroServiceInstance, interface{}) selector.Next) {
	strategies[name] = strategy
	lager.Logger.Debugf("Installed strategy plugin: %s.", name)
}

// GetStrategyPlugin get strategy plugin
func GetStrategyPlugin(name string) (func([]*registry.MicroServiceInstance, interface{}) selector.Next, error) {
	s, ok := strategies[name]
	if !ok {
		return nil, fmt.Errorf("Don't support strategyName [%s]", name)
	}

	return s, nil
}
