package archaius

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
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
func (p *Panel) GetLoadBalancing(inv invocation.Invocation) model.LoadBalancingSpec {
	return model.LoadBalancingSpec{}

}

//GetRateLimiting get rate limiting config
func (p *Panel) GetRateLimiting(inv invocation.Invocation, serviceType string) model.FlowControl {
	return model.FlowControl{}
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
