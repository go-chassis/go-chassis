package archaius

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/qpslimiter"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"strings"
)

//Panel pull configs from archaius
type Panel struct {
}

func newPanel(options control.Options) control.Panel {
	return &Panel{}
}

//GetCircuitBreaker return command , and circuit breaker settings
func (p *Panel) GetCircuitBreaker(inv invocation.Invocation, serviceType string) (string, hystrix.CommandConfig) {
	command := serviceType
	if inv.MicroServiceName != "" {
		command = strings.Join([]string{serviceType, inv.MicroServiceName}, ".")
	}
	return command, hystrix.CommandConfig{
		ForceFallback:          config.GetForceFallback(inv.MicroServiceName, serviceType),
		TimeoutEnabled:         config.GetTimeoutEnabled(inv.MicroServiceName, serviceType),
		Timeout:                config.GetTimeout(command, serviceType),
		MaxConcurrentRequests:  config.GetMaxConcurrentRequests(command, serviceType),
		ErrorPercentThreshold:  config.GetErrorPercentThreshold(command, serviceType),
		RequestVolumeThreshold: config.GetRequestVolumeThreshold(command, serviceType),
		SleepWindow:            config.GetSleepWindow(command, serviceType),
		ForceClose:             config.GetForceClose(inv.MicroServiceName, serviceType),
		ForceOpen:              config.GetForceOpen(inv.MicroServiceName, serviceType),
		CircuitBreakerEnabled:  config.GetCircuitBreakerEnabled(command, serviceType),
	}
}

//GetLoadBalancing get load balancing config
func (p *Panel) GetLoadBalancing(inv invocation.Invocation) control.LoadBalancingConfig {
	c := control.LoadBalancingConfig{}
	c.Strategy = config.GetStrategyName(inv.SourceMicroService, inv.MicroServiceName)
	c.Filters = config.GetServerListFilters()
	c.RetryEnabled = config.RetryEnabled(inv.SourceMicroService, inv.MicroServiceName)
	c.RetryOnSame = config.GetRetryOnSame(inv.SourceMicroService, inv.MicroServiceName)
	c.RetryOnNext = config.GetRetryOnNext(inv.SourceMicroService, inv.MicroServiceName)
	c.BackOffKind = config.BackOffKind(inv.SourceMicroService, inv.MicroServiceName)
	c.BackOffMax = config.BackOffMaxMs(inv.SourceMicroService, inv.MicroServiceName)
	c.BackOffMin = config.BackOffMinMs(inv.SourceMicroService, inv.MicroServiceName)
	return c

}

//GetRateLimiting get rate limiting config
func (p *Panel) GetRateLimiting(inv invocation.Invocation, serviceType string) control.RateLimitingConfig {
	rl := control.RateLimitingConfig{}
	rl.Enabled = archaius.GetBool("cse.flowcontrol."+serviceType+".qps.enabled", true)
	operationMeta := qpslimiter.InitSchemaOperations(&inv)
	rl.Rate, rl.Key = qpslimiter.GetQPSTrafficLimiter().GetQPSRateWithPriority(operationMeta)
	return rl
}

//GetFaultInjection get Fault injection config
func (p *Panel) GetFaultInjection(inv invocation.Invocation) model.Fault {
	return model.Fault{}

}

//GetEgressRule get egress config
func (p *Panel) GetEgressRule(inv invocation.Invocation) {

}

func init() {
	control.InstallPlugin("archaius", newPanel)
}
