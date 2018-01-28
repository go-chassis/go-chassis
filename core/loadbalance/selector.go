// Package loadbalance is a way to load balance service nodes
package loadbalance

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
)

// constant strings for load balance variables
const (
	StrategyRoundRobin        = "RoundRobin"
	StrategyRandom            = "Random"
	StrategySessionStickiness = "SessionStickiness"
	StrategyLatency           = "WeightedResponse"
)

var (
	// DefaultSelector is the object of selector
	DefaultSelector selector.Selector
)

// Enable function is for to enable load balance strategy
func Enable() error {
	lager.Logger.Info("Enable LoadBalancing")
	InstallStrategy(StrategyRandom, selector.Random)
	InstallStrategy(StrategyRoundRobin, selector.RoundRobin)
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
	strategyFunc := selector.SetStrategy(strategy)
	DefaultSelector = newDefaultSelector(strategyFunc)
	return nil
}
