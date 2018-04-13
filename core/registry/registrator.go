package registry

import (
	"fmt"
	"strings"
	"sync"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"log"
)

// constant values for registry parameters
const (
	DefaultRegistratorPlugin       = "servicecenter"
	DefaultServiceDiscoveryPlugin  = "servicecenter"
	DefaultContractDiscoveryPlugin = "servicecenter"
	Name                           = "registry"
	SDTag                          = "serviceDiscovery"
	CDTag                          = "contractDiscovery"
	RTag                           = "registrator"
	Auto                           = "auto"
	Manual                         = "manual"
)

// IsEnabled check enable
var IsEnabled bool
var mu sync.Mutex

// DefaultRegistrator is the client of registry, you can call the method of it to interact with microservice registry
var DefaultRegistrator Registrator

// DefaultAddr default address of service center
var DefaultAddr = "http://127.0.0.1:30100"

// registryFunc registry function
var registryFunc = make(map[string]func(opts Options) Registrator)

// HBService variable of heartbeat service
var HBService = &HeartbeatService{
	instances: make(map[string]*HeartbeatTask),
}

// Registrator is the interface for developer to update information in service registry
type Registrator interface {
	//Close destroy connection between the registry client and server
	Close() error
	//RegisterService register a microservice to registry, if it is duplicated in registry, it returns error
	RegisterService(microService *MicroService) (string, error)
	//RegisterServiceInstance register a microservice instance to registry
	RegisterServiceInstance(sid string, instance *MicroServiceInstance) (string, error)
	RegisterServiceAndInstance(microService *MicroService, instance *MicroServiceInstance) (string, string, error)
	Heartbeat(microServiceID, microServiceInstanceID string) (bool, error)
	AddDependencies(dep *MicroServiceDependency) error
	UnRegisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error
	UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error
	UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error
	UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error
	AddSchemas(microServiceID, schemaName, schemaInfo string) error
}

func enableRegistrator(opts Options) {
	rt := config.GlobalDefinition.Cse.Service.Registry.Type
	if rt == "" {
		rt = DefaultRegistratorPlugin
	}
	f := registryFunc[rt]
	if f == nil {
		panic("No registry plugin")
	}
	DefaultRegistrator = f(opts)
	lager.Logger.Warnf(nil, "Enable %s registry.", DefaultRegistrator)
}

// InstallRegistrator install registrator plugin
func InstallRegistrator(name string, f func(opts Options) Registrator) {
	registryFunc[name] = f
	log.Printf("Installed registry plugin: %s.\n", name)
}
func setSpecifiedOptions(oR, oSD, oCD Options) error {
	lager.Logger.Info("Doesn't set address for registry, so use Registrator, Discovery separated configs")
	hostsR, schemeR, err := URIs2Hosts(strings.Split(config.GlobalDefinition.Cse.Service.Registrator.Address, ","))
	if err != nil {
		return err
	}
	oR.Addrs = hostsR
	oR.Tenant = config.GlobalDefinition.Cse.Service.Registrator.Tenant
	oR.Version = config.GlobalDefinition.Cse.Service.Registrator.APIVersion.Version
	oR.TLSConfig, err = getTLSConfig(schemeR, RTag)
	if err != nil {
		return err
	}
	if oR.TLSConfig != nil {
		oR.EnableSSL = true
	}
	hostsSD, schemeSD, err := URIs2Hosts(strings.Split(config.GlobalDefinition.Cse.Service.ServiceDiscovery.Address, ","))
	if err != nil {
		return err
	}
	oSD.Addrs = hostsSD
	oSD.Tenant = config.GlobalDefinition.Cse.Service.ServiceDiscovery.Tenant
	oSD.Version = config.GlobalDefinition.Cse.Service.ServiceDiscovery.APIVersion.Version
	oSD.TLSConfig, err = getTLSConfig(schemeSD, SDTag)
	if err != nil {
		return err
	}
	if oSD.TLSConfig != nil {
		oSD.EnableSSL = true
	}
	hostsCD, schemeCD, err := URIs2Hosts(strings.Split(config.GlobalDefinition.Cse.Service.ContractDiscovery.Address, ","))
	if err != nil {
		return err
	}
	oCD.Addrs = hostsCD
	oCD.Tenant = config.GlobalDefinition.Cse.Service.ContractDiscovery.Tenant
	oCD.Version = config.GlobalDefinition.Cse.Service.ContractDiscovery.APIVersion.Version
	oCD.TLSConfig, err = getTLSConfig(schemeCD, CDTag)
	if err != nil {
		return err
	}
	if oCD.TLSConfig != nil {
		oCD.EnableSSL = true
	}
	return nil
}

// Enable create DefaultRegistrator
func Enable() error {
	mu.Lock()
	defer mu.Unlock()
	if IsEnabled {
		return nil
	}

	if config.GlobalDefinition.Cse.Service.Registry.Tenant == "" {
		config.GlobalDefinition.Cse.Service.Registry.Tenant = common.DefaultApp
	}

	var scheme string
	hosts, scheme, err := URIs2Hosts(strings.Split(config.GlobalDefinition.Cse.Service.Registry.Address, ","))
	if err != nil {
		return err
	}
	tlsConfig, err := getTLSConfig(scheme, Name)
	if err != nil {
		return err
	}
	var secure bool
	if tlsConfig != nil {
		secure = true
	}
	oR := Options{
		Addrs:     hosts,
		Tenant:    config.GlobalDefinition.Cse.Service.Registry.Tenant,
		EnableSSL: secure,
		TLSConfig: tlsConfig,
		Version:   config.GlobalDefinition.Cse.Service.Registry.APIVersion.Version,
	}
	oSD := Options{
		Addrs:     hosts,
		Tenant:    config.GlobalDefinition.Cse.Service.Registry.Tenant,
		EnableSSL: secure,
		TLSConfig: tlsConfig,
		Version:   config.GlobalDefinition.Cse.Service.Registry.APIVersion.Version,
	}
	oCD := Options{
		Addrs:     hosts,
		Tenant:    config.GlobalDefinition.Cse.Service.Registry.Tenant,
		EnableSSL: secure,
		TLSConfig: tlsConfig,
		Version:   config.GlobalDefinition.Cse.Service.Registry.APIVersion.Version,
	}
	if len(hosts) == 0 {
		if err := setSpecifiedOptions(oR, oSD, oCD); err != nil {
			return err
		}
	}

	enableRegistrator(oR)
	enableContractDiscovery(oSD)
	enableServiceDiscovery(oCD)
	if err := RegisterMicroservice(); err != nil {
		lager.Logger.Errorf(err, "start backoff for register microservice")
		startBackOff(RegisterMicroservice)
	}
	go HBService.Start()
	DefaultServiceDiscoveryService.AutoSync()
	lager.Logger.Info("Enabled Registry")
	IsEnabled = true
	return nil
}

// DoRegister for registering micro-service instances
func DoRegister() error {
	var isAutoRegister bool
	switch config.GlobalDefinition.Cse.Service.Registry.AutoRegister {
	case "":
		isAutoRegister = true
	case Auto:
		isAutoRegister = true
	case Manual:
		isAutoRegister = false
	default:
		{
			tmpErr := fmt.Errorf("parameter incorrect, autoregister: %s", config.GlobalDefinition.Cse.Service.Registry.AutoRegister)
			lager.Logger.Error(tmpErr.Error(), nil)
			return tmpErr
		}
	}
	if isAutoRegister {
		if err := RegisterMicroserviceInstances(); err != nil {
			lager.Logger.Errorf(err, "start back off for register microservice instances background")
			go startBackOff(RegisterMicroserviceInstances)
		}
	}
	return nil
}
