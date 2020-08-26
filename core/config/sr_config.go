package config

import (
	"github.com/go-chassis/go-archaius"
)

// GetRegistratorType returns the Type of service registry
func GetRegistratorType() string {
	return GlobalDefinition.ServiceComb.Registry.Type
}

// GetRegistratorAddress returns the Address of service registry
func GetRegistratorAddress() string {
	return GlobalDefinition.ServiceComb.Registry.Address
}

// GetRegistratorScope returns the Scope of service registry
func GetRegistratorScope() string {
	return GlobalDefinition.ServiceComb.Registry.Scope
}

// GetRegistratorAutoRegister returns the AutoRegister of service registry
func GetRegistratorAutoRegister() string {
	return GlobalDefinition.ServiceComb.Registry.AutoRegister
}

// GetRegistratorAPIVersion returns the APIVersion of service registry
func GetRegistratorAPIVersion() string {
	return GlobalDefinition.ServiceComb.Registry.APIVersion.Version
}

// GetRegistratorDisable returns the Disable of service registry
func GetRegistratorDisable() bool {
	return archaius.GetBool("servicecomb.registry.disabled", false)
}
