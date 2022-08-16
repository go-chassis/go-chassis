package servicecomb

import (
	"github.com/go-chassis/go-chassis/v2/control"
	"github.com/go-chassis/go-chassis/v2/core/loadbalancer"
	"github.com/go-chassis/go-chassis/v2/resilience/retry"
	"github.com/patrickmn/go-cache"
)

// save configs
var (
	//key is service name
	LBConfigCache = cache.New(0, 0)
	//key is [Provider|Consumer]:service_name or [Consumer|Provider]
	CBConfigCache     = cache.New(0, 0)
	RLConfigCache     = cache.New(0, 0)
	EgressConfigCache = cache.New(0, 0)
	FIConfigCache     = cache.New(0, 0)
)

// Default values
var (
	DefaultLB = control.LoadBalancingConfig{
		Strategy:    loadbalancer.StrategyRoundRobin,
		BackOffKind: retry.DefaultBackOffKind,
	}
)
