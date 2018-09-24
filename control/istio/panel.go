package istio

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
)

func init() {
	control.InstallPlugin("pilot", newPilotPanel)
}

//PilotPanel pull configs from istio pilot
type PilotPanel struct {
}

func newPilotPanel(options control.Options) control.Panel {
	return &PilotPanel{}
}

//GetEgressRule get egress config
func (p *PilotPanel) GetEgressRule() []control.EgressConfig {
	return []control.EgressConfig{}
}

//GetCircuitBreaker return command , and circuit breaker settings
func (p *PilotPanel) GetCircuitBreaker(inv invocation.Invocation, serviceType string) (string, hystrix.CommandConfig) {
	return "", hystrix.CommandConfig{}

}

//GetLoadBalancing get load balancing config
func (p *PilotPanel) GetLoadBalancing(inv invocation.Invocation) control.LoadBalancingConfig {
	return control.LoadBalancingConfig{}

}

//GetRateLimiting get rate limiting config
func (p *PilotPanel) GetRateLimiting(inv invocation.Invocation, serviceType string) control.RateLimitingConfig {
	return control.RateLimitingConfig{}
}

//GetFaultInjection get Fault injection config
func (p *PilotPanel) GetFaultInjection(inv invocation.Invocation) model.Fault {
	return model.Fault{}
}
