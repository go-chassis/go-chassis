package fault

import (
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/invocation"
)

// InjectFault inject fault
type InjectFault func(model.Fault, *invocation.Invocation) error

// FaultInjectors fault injectors
var FaultInjectors = make(map[string]InjectFault)

//FaultError fault injection error
type FaultError struct {
	Message string
}

func (e FaultError) Error() string {
	return e.Message
}

// InstallFaultInjectionPlugin install fault injection plugin
func InstallFaultInjectionPlugin(name string, f InjectFault) {
	FaultInjectors[name] = f
}

func init() {
	InstallFaultInjectionPlugin("rest", faultInject)
	InstallFaultInjectionPlugin("highway", faultInject)
	InstallFaultInjectionPlugin("dubbo", faultInject)
}

func faultInject(rule model.Fault, inv *invocation.Invocation) error {
	return ValidateAndApplyFault(&rule, inv)
}
