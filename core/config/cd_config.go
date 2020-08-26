package config

import "github.com/go-chassis/go-archaius"

// GetContractDiscoveryType returns the Type of contract discovery registry
func GetContractDiscoveryType() string {
	return GlobalDefinition.ServiceComb.Registry.Type
}

// GetContractDiscoveryAddress returns the Address of contract discovery registry
func GetContractDiscoveryAddress() string {
	return GlobalDefinition.ServiceComb.Registry.Address
}

// GetContractDiscoveryAPIVersion returns the APIVersion of contract discovery registry
func GetContractDiscoveryAPIVersion() string {
	return GlobalDefinition.ServiceComb.Registry.APIVersion.Version
}

// GetContractDiscoveryDisable returns the Disable of contract discovery registry
func GetContractDiscoveryDisable() bool {
	return archaius.GetBool("servicecomb.registry.disabled", false)
}
