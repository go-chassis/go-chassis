package registry

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"log"
)

var sdFunc = make(map[string]func(opts Options) ServiceDiscovery)

var cdFunc = make(map[string]func(opts Options) ContractDiscovery)

//InstallServiceDiscovery install service discovery client
func InstallServiceDiscovery(name string, f func(opts Options) ServiceDiscovery) {
	sdFunc[name] = f
	log.Printf("Installed service discovery plugin: %s.\n", name)
}

//InstallContractDiscovery install contract service client
func InstallContractDiscovery(name string, f func(opts Options) ContractDiscovery) {
	cdFunc[name] = f
	log.Printf("Installed contract discovery plugin: %s.\n", name)
}

//ServiceDiscovery fetch service and instances from remote or local
type ServiceDiscovery interface {
	GetMicroServiceID(appID, microServiceName, version, env string) (string, error)
	GetAllMicroServices() ([]*MicroService, error)
	GetMicroService(microServiceID string) (*MicroService, error)
	GetMicroServiceInstances(consumerID, providerID string) ([]*MicroServiceInstance, error)
	// FindMicroServiceInstances find instances of a service specified by appID, microServiceName, version and env
	FindMicroServiceInstances(consumerID, appID, microServiceName, version, env string) ([]*MicroServiceInstance, error)
	AutoSync()
	Close() error
}

//DefaultServiceDiscoveryService supplies service discovery
var DefaultServiceDiscoveryService ServiceDiscovery

// DefaultContractDiscoveryService supplies contract discovery
var DefaultContractDiscoveryService ContractDiscovery

//ContractDiscovery fetch schema content from remote or local
type ContractDiscovery interface {
	GetMicroServicesByInterface(interfaceName string) (microservices []*MicroService)
	GetSchemaContentByInterface(interfaceName string) SchemaContent
	GetSchemaContentByServiceName(svcName, version, appID, env string) []*SchemaContent
	Close() error
}

func enableServiceDiscovery(opts Options) {
	if config.GetServiceDiscoveryDisable() {
		return
	}

	t := config.GetServiceDiscoveryType()
	if t == "" {
		t = DefaultServiceDiscoveryPlugin
	}
	f := sdFunc[t]
	if f == nil {
		panic("No service discovery plugin")
	}
	DefaultServiceDiscoveryService = f(opts)

	DefaultServiceDiscoveryService.AutoSync()

	lager.Logger.Infof("Enable %s service discovery.", t)
}

func enableContractDiscovery(opts Options) {
	if config.GetContractDiscoveryDisable() {
		return
	}

	t := config.GetContractDiscoveryType()
	if t == "" {
		t = DefaultContractDiscoveryPlugin
	}
	f := cdFunc[t]
	if f == nil {
		lager.Logger.Warn("No contract discovery plugin", nil)
		return
	}
	DefaultContractDiscoveryService = f(opts)
	lager.Logger.Infof("Enable %s contract discovery.", t)
}
