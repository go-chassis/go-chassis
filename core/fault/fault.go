package fault

import (
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/invocation"
)

// InjectFault inject fault
type InjectFault func(model.Fault, *invocation.Invocation) error

// FaultInjectors fault injectors
var FaultInjectors = make(map[string]InjectFault)

// InstallFaultInjectionPlugin install fault injection plugin
func InstallFaultInjectionPlugin(name string, f InjectFault) {
	FaultInjectors[name] = f
}

func init() {
	InstallFaultInjectionPlugin("rest", faultInject)
	InstallFaultInjectionPlugin("highway", faultInject)
}

func faultInject(rule model.Fault, inv *invocation.Invocation) error {
	return ValidateAndApplyFault(&rule, inv)
}
