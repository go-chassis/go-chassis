package registry

import (
	"fmt"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-chassis/openlog"
)

var sdFunc = make(map[string]func(opts Options) ServiceDiscovery)

var cdFunc = make(map[string]func(opts Options) ContractDiscovery)

//InstallServiceDiscovery install service discovery client
func InstallServiceDiscovery(name string, f func(opts Options) ServiceDiscovery) {
	sdFunc[name] = f
	openlog.Info("Installed service discovery plugin: " + name)
}

//NewDiscovery create discovery service
func NewDiscovery(name string, opts Options) (ServiceDiscovery, error) {
	f := sdFunc[name]
	if f == nil {
		return nil, fmt.Errorf("no service discovery plugin: %s", name)
	}
	return f(opts), nil
}

//InstallContractDiscovery install contract service client
func InstallContractDiscovery(name string, f func(opts Options) ContractDiscovery) {
	cdFunc[name] = f
	openlog.Info("Installed contract discovery plugin: " + name)
}

//ServiceDiscovery fetch service and instances from remote or local
type ServiceDiscovery interface {
	GetMicroService(microServiceID string) (*MicroService, error)
	FindMicroServiceInstances(consumerID, microServiceName string, tags utiltags.Tags) ([]*MicroServiceInstance, error)
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

func enableServiceDiscovery(opts Options) error {
	if config.GetServiceDiscoveryDisable() {
		openlog.Warn("discovery is disabled")
		return nil
	}

	t := config.GetServiceDiscoveryType()
	if t == "" {
		t = DefaultServiceDiscoveryPlugin
	}
	f := sdFunc[t]
	if f == nil {
		panic("No service discovery plugin")
	}
	var err error
	DefaultServiceDiscoveryService, err = NewDiscovery(t, opts)
	if err != nil {
		return err
	}

	DefaultServiceDiscoveryService.AutoSync()

	openlog.Info(fmt.Sprintf("enable %s service discovery.", t))
	return nil
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
		openlog.Warn("No contract discovery plugin")
		return
	}
	DefaultContractDiscoveryService = f(opts)
	openlog.Info(fmt.Sprintf("Enable %s contract discovery.", t))
}
