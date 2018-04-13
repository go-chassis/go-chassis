// Package loadbalancer is client side load balancer
package loadbalancer

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"

	"fmt"
	"github.com/ServiceComb/go-chassis/core/registry"
)

// constant strings for load balance variables
const (
	StrategyRoundRobin        = "RoundRobin"
	StrategyRandom            = "Random"
	StrategySessionStickiness = "SessionStickiness"
	StrategyLatency           = "WeightedResponse"
	OperatorEqual             = "="
	OperatorGreater           = ">"
	OperatorSmaller           = "<"
	OperatorPattern           = "Pattern"
)

var (
	// ErrNoneAvailableInstance is to represent load balance error
	ErrNoneAvailableInstance = LBError{Message: "None available instance"}
)

// LBError load balance error
type LBError struct {
	Message string
}

// Error for to return load balance error message
func (e LBError) Error() string {
	return "lb: " + e.Message
}

// BuildStrategy query instance list and give it to Strategy then return Strategy
func BuildStrategy(consumerID, serviceName, app, version, protocol, sessionID string, fs []Filter, s Strategy, metadata interface{}) (Strategy, error) {
	if s == nil {
		s = &RoundRobinStrategy{}
	}

	var isFilterExist = true
	for _, filter := range fs {
		if filter == nil {
			isFilterExist = false
		}

	}

	if app == "" {
		app = config.GlobalDefinition.AppID
	}

	instances, err := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(consumerID, app, serviceName, version, "")
	if err != nil {
		lbErr := LBError{err.Error()}
		lager.Logger.Errorf(lbErr, "Lb err")
		return nil, lbErr
	}

	if isFilterExist {
		for _, filter := range fs {
			instances = filter(instances, nil)
		}

	}

	if len(instances) == 0 {
		lbErr := LBError{fmt.Sprintf("No available instance, key: %s:%s:%s", app, serviceName, version)}
		lager.Logger.Error(lbErr.Error(), nil)
		return nil, lbErr
	}
	s.ReceiveData(instances, serviceName, protocol, sessionID)
	return s, nil
}

// Strategy is load balancer algorithm , call Pick to return one instance
type Strategy interface {
	ReceiveData(instances []*registry.MicroServiceInstance, serviceName, protocol, sessionID string)
	Pick() (*registry.MicroServiceInstance, error)
}

//Criteria is rule for filter
type Criteria struct {
	Key      string
	Operator string
	Value    string
}

// Filter receive instances and criteria, it will filter instances based on criteria you defined,criteria is optional, you can give nil for it
type Filter func(instances []*registry.MicroServiceInstance, criteria []*Criteria) []*registry.MicroServiceInstance

// Enable function is for to enable load balance strategy
func Enable() error {
	lager.Logger.Info("Enable LoadBalancing")
	InstallStrategy(StrategyRandom, newRandomStrategy)
	InstallStrategy(StrategyRoundRobin, newRoundRobinStrategy)
	InstallStrategy(StrategySessionStickiness, newSessionStickinessStrategy)
	InstallStrategy(StrategyLatency, newWeightedResponseStrategy)

	var strategyName string

	strategyName = config.GetLoadBalancing().Strategy["name"]
	strategyNameFromArchaius := archaius.GetString("cse.loadbalance.strategy.name", "")
	if strategyName == "" && archaius.Get("cse.loadbalance.strategy.name") == "" {
		lager.Logger.Info("Empty strategy configuration, use RoundRobin as default")
		return nil
	}
	lager.Logger.Info("Strategy is " + strategyNameFromArchaius)

	return nil
}
