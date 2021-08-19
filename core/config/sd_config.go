package config

import "github.com/go-chassis/go-archaius"

// GetServiceDiscoveryType returns the Type of SD registry
func GetServiceDiscoveryType() string {
	return GlobalDefinition.ServiceComb.Registry.Type
}

// GetServiceDiscoveryAddress returns the Address of SD registry
func GetServiceDiscoveryAddress() string {
	return GlobalDefinition.ServiceComb.Registry.Address
}

// GetServiceDiscoveryRefreshInterval returns the RefreshInterval of SD registry
func GetServiceDiscoveryRefreshInterval() string {
	return GlobalDefinition.ServiceComb.Registry.RefreshInterval
}

// GetServiceDiscoveryWatch returns the Watch of SD registry
func GetServiceDiscoveryWatch() bool {
	return GlobalDefinition.ServiceComb.Registry.Watch
}

// GetServiceDiscoveryAPIVersion returns the APIVersion of SD registry
func GetServiceDiscoveryAPIVersion() string {
	return GlobalDefinition.ServiceComb.Registry.APIVersion.Version
}

// GetServiceDiscoveryDisable returns the Disable of SD registry
func GetServiceDiscoveryDisable() bool {
	return archaius.GetBool("servicecomb.registry.discovery.disabled", false)
}

// GetServiceDiscoveryHealthCheck returns the HealthCheck of SD registry
func GetServiceDiscoveryHealthCheck() bool {
	return archaius.GetBool("servicecomb.registry.healthCheck", false)
}

// GetServiceDiscoveryUploadSchema returns if should register schema of SD registry
func GetServiceDiscoveryUploadSchema() bool {
	return archaius.GetBool("servicecomb.registry.uploadSchema", false)
}

// DefaultConfigPath set the default config path
const DefaultConfigPath = "/etc/.kube/config"

// GetServiceDiscoveryConfigPath returns the configpath of SD registry
func GetServiceDiscoveryConfigPath() string {
	if GlobalDefinition.ServiceComb.Registry.ConfigPath != "" {
		return GlobalDefinition.ServiceComb.Registry.ConfigPath
	}
	return DefaultConfigPath
}
