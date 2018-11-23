package control

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"strings"
)

var panelPlugin = make(map[string]func(options Options) Panel)

//DefaultPanel get fetch config
var DefaultPanel Panel

const (
	//ScopeAPI is config const
	ScopeAPI = "api"
)

//Panel is a abstraction of pulling configurations from various of systems, and transfer different configuration into standardized model
//you can use different panel implementation to pull different of configs from Istio or Archaius
//TODO able to set configs
type Panel interface {
	GetCircuitBreaker(inv invocation.Invocation, serviceType string) (string, hystrix.CommandConfig)
	GetLoadBalancing(inv invocation.Invocation) LoadBalancingConfig
	GetRateLimiting(inv invocation.Invocation, serviceType string) RateLimitingConfig
	GetFaultInjection(inv invocation.Invocation) model.Fault
	GetEgressRule() []EgressConfig
}

//Options is options
type Options struct {
	Address string
}

//InstallPlugin install implementation
func InstallPlugin(name string, f func(options Options) Panel) {
	panelPlugin[name] = f
}

//Init initialize DefaultPanel
func Init() error {
	infra := config.GlobalDefinition.Panel.Infra
	if infra == "" {
		infra = "archaius"
	}
	f, ok := panelPlugin[infra]
	if !ok {
		return fmt.Errorf("do not support [%s] panel", infra)
	}

	DefaultPanel = f(Options{
		Address: config.GlobalDefinition.Panel.Settings["address"],
	})
	return nil
}

// NewCircuitName create circuit command
func NewCircuitName(serviceType string, inv invocation.Invocation) string {
	var cmd = serviceType
	if inv.MicroServiceName != "" {
		cmd = strings.Join([]string{cmd, inv.MicroServiceName}, ".")
	}
	if config.GetHystrixConfig().CircuitBreakerProperties.Scope == "" {
		config.GetHystrixConfig().CircuitBreakerProperties.Scope = ScopeAPI
	}
	if config.GetHystrixConfig().CircuitBreakerProperties.Scope == ScopeAPI {
		if inv.SchemaID != "" {
			cmd = strings.Join([]string{cmd, inv.SchemaID}, ".")
		}
		if inv.OperationID != "" {
			cmd = strings.Join([]string{cmd, inv.OperationID}, ".")
		}
	}
	return cmd

}
