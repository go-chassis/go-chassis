// Package loadbalance is a way to load balance service nodes
package loadbalance

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
)

// constant strings for load balance variables
const (
	StrategyRoundRobin        = "RoundRobin"
	StrategyRandom            = "Random"
	StrategySessionStickiness = "SessionStickiness"
	StrategyLatency           = "WeightedResponse"
)

// Selector builds on the registry as a mechanism to pick nodes
// and mark their status. This allows host pools and other things
// to be built using various algorithms.
type Selector interface {
	Init(opts ...Option) error
	Options() Options
	// Select returns a function which should return the next node
	Select(microserviceName, version string, opts ...SelectOption) (Next, error)
	// Name of the selector
	String() string
}

// Next is a function that returns the next node
// based on the selector's strategy
type Next func() (*registry.MicroServiceInstance, error)

// Filter is used to filter a service during the selection process
type Filter func([]*registry.MicroServiceInstance) []*registry.MicroServiceInstance

// Strategy is a selection strategy e.g random, round robin
type Strategy func([]*registry.MicroServiceInstance, interface{}) Next

var (
	// DefaultSelector is the object of selector
	DefaultSelector Selector
	// ErrNoneAvailable is to represent load balance error
	ErrNoneAvailable = LBError{Message: "No available"}
)

// LBError load balance error
type LBError struct {
	Message string
}

// Error for to return load balance error message
func (e LBError) Error() string {
	return "lb: " + e.Message
}

// Enable function is for to enable load balance strategy
func Enable() error {
	lager.Logger.Info("Enable LoadBalancing")
	InstallStrategy(StrategyRandom, Random)
	InstallStrategy(StrategyRoundRobin, RoundRobin)
	InstallStrategy(StrategySessionStickiness, SessionStickiness)
	InstallStrategy(StrategyLatency, WeightedResponse)

	var strategyName string

	strategyName = config.GlobalDefinition.Cse.Loadbalance.Strategy["name"]

	if strategyName == "" && archaius.Get("cse.loadbalance.strategy.name") == "" {
		lager.Logger.Info("Empty strategy configuration, use RoundRobin as default")
		DefaultSelector = newDefaultSelector()
		return nil
	}

	if archaius.GetString("cse.loadbalance.strategy.name", "") != "" {
		strategyName = archaius.GetString("cse.loadbalance.strategy.name", "")
	} else {
		DefaultSelector = newDefaultSelector()
		return nil
	}
	strategy, err := GetStrategyPlugin(strategyName)
	if err != nil {
		lager.Logger.Errorf(err, "Get strategy ["+strategyName+"] failed")
		return err
	}
	lager.Logger.Info("Load balancing strategy is " + strategyName)
	strategyFunc := SetStrategy(strategy)
	DefaultSelector = newDefaultSelector(strategyFunc)
	return nil
}
