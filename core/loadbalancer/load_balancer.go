// Package loadbalancer is client side load balancer
package loadbalancer

import (
	"fmt"
	"strings"

	"github.com/go-chassis/go-chassis/core/archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
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
func BuildStrategy(consumerID, serviceName, protocol, sessionID string, fs []string,
	s Strategy, tags utiltags.Tags) (Strategy, error) {
	if s == nil {
		s = &RoundRobinStrategy{}
	}

	var isFilterExist = true
	for _, filter := range fs {
		if filter == "" {
			isFilterExist = false
		}

	}

	instances, err := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(consumerID, serviceName, tags)
	if err != nil {
		lbErr := LBError{err.Error()}
		lager.Logger.Errorf("Lb err: %s", err)
		return nil, lbErr
	}

	if isFilterExist {
		filterFuncs := make([]Filter, 0)
		//append filters in config
		for _, fName := range fs {
			f := Filters[fName]
			if f != nil {
				filterFuncs = append(filterFuncs, f)
				continue
			}
		}
		for _, filter := range filterFuncs {
			instances = filter(instances, nil)
		}
	}

	if len(instances) == 0 {
		lbErr := LBError{fmt.Sprintf("No available instance, key: %s(%v)", serviceName, tags)}
		lager.Logger.Error(lbErr.Error())
		return nil, lbErr
	}

	serviceKey := strings.Join([]string{serviceName, tags.String()}, "|")
	s.ReceiveData(instances, serviceKey, protocol, sessionID)
	return s, nil
}

// Strategy is load balancer algorithm , call Pick to return one instance
type Strategy interface {
	ReceiveData(instances []*registry.MicroServiceInstance, serviceKey, protocol, sessionID string)
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
