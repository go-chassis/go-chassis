package servicecomb

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/control"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/resilience/rate"
	"github.com/go-chassis/go-chassis/v2/third_party/forked/afex/hystrix-go/hystrix"
)

//Panel pull configs from archaius
type Panel struct {
}

func newPanel(options control.Options) control.Panel {
	SaveToLBCache(config.GetLoadBalancing())
	SaveToCBCache(config.GetHystrixConfig())
	return &Panel{}
}

//GetCircuitBreaker return command , and circuit breaker settings
func (p *Panel) GetCircuitBreaker(inv invocation.Invocation, serviceType string) (string, hystrix.CommandConfig) {
	key := GetCBCacheKey(inv.MicroServiceName, serviceType)
	command := control.NewCircuitName(serviceType, config.GetHystrixConfig().CircuitBreakerProperties.Scope, inv)
	c, ok := CBConfigCache.Get(key)
	if !ok {
		c, _ := CBConfigCache.Get(serviceType)
		return command, c.(hystrix.CommandConfig)

	}
	return command, c.(hystrix.CommandConfig)
}

//GetLoadBalancing get load balancing config
func (p *Panel) GetLoadBalancing(inv invocation.Invocation) control.LoadBalancingConfig {
	c, ok := LBConfigCache.Get(inv.MicroServiceName)
	if !ok {
		c, ok := LBConfigCache.Get("")
		if !ok {
			return DefaultLB

		}
		return c.(control.LoadBalancingConfig)

	}
	return c.(control.LoadBalancingConfig)

}

//GetRateLimiting get rate limiting config
func (p *Panel) GetRateLimiting(inv invocation.Invocation, serviceType string) control.RateLimitingConfig {
	rl := control.RateLimitingConfig{}
	rl.Enabled = archaius.GetBool("servicecomb.flowcontrol."+serviceType+".qps.enabled", true)
	if serviceType == common.Consumer {
		keys := GetConsumerKey(inv.SourceMicroService, inv.MicroServiceName, inv.SchemaID, inv.OperationID)
		rl.Rate, rl.Key = GetQPSRateWithPriority(
			keys.OperationQualifiedName, keys.SchemaQualifiedName, keys.MicroServiceName)
	} else {
		keys := GetProviderKey(inv.SourceMicroService)
		rl.Rate, rl.Key = GetQPSRateWithPriority(
			keys.ServiceOriented, keys.Global)
	}

	return rl
}

//GetFaultInjection get Fault injection config
func (p *Panel) GetFaultInjection(inv invocation.Invocation) model.Fault {
	return model.Fault{}

}

//GetEgressRule get egress config
func (p *Panel) GetEgressRule() []control.EgressConfig {
	return []control.EgressConfig{}
}

func init() {
	control.InstallPlugin("archaius", newPanel)
}

// GetQPSRateWithPriority get qps rate with priority
func GetQPSRateWithPriority(cmd ...string) (int, string) {
	var (
		qpsVal      int
		configExist bool
	)
	for _, c := range cmd {
		qpsVal, configExist = GetQPSRate(c)
		if configExist {
			return qpsVal, c
		}
	}

	return rate.DefaultRate, cmd[len(cmd)-1]

}

// GetQPSRate get qps rate
func GetQPSRate(rateConfig string) (int, bool) {
	qpsRate := archaius.GetInt(rateConfig, rate.DefaultRate)
	if qpsRate == rate.DefaultRate {
		return qpsRate, false
	}

	return qpsRate, true
}
