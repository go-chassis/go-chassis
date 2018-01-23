package registry

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
)

// constant values for registry parameters
const (
	DefaultRegistryPlugin = "servicecenter"
	Name                  = "registry"
	Auto                  = "auto"
	Manual                = "manual"
)

// IsEnabled check enable
var IsEnabled bool
var mu sync.Mutex

// RegistryService is the client of registry, you can call the method of it to interact with microservice registry
var RegistryService Registry

// DefaultAddr default address of service center
var DefaultAddr = "http://127.0.0.1:30100"

// registryFunc registry function
var registryFunc = make(map[string]func(opts ...Option) Registry)

// HBService variable of heartbeat service
var HBService = &HeartbeatService{
	instances: make(map[string]*HeartbeatTask),
}

// Registry is the interface for developer to interact with microservice registry
type Registry interface {
	//Close destroy connection between the registry client and server
	Close() error
	//RegisterService register a microservice to registry, if it is duplicated in registry, it returns error
	RegisterService(microService *MicroService) (string, error)
	//RegisterServiceInstance register a microservice instance to registry
	RegisterServiceInstance(sid string, instance *MicroServiceInstance) (string, error)
	RegisterServiceAndInstance(microService *MicroService, instance *MicroServiceInstance) (string, string, error)
	Heartbeat(microServiceID, microServiceInstanceID string) (bool, error)
	AddDependencies(dep *MicroServiceDependency) error
	GetMicroServiceID(appID, microServiceName, version string) (string, error)
	GetAllMicroServices() ([]*MicroService, error)
	GetMicroService(microServiceID string) (*MicroService, error)
	GetMicroServiceInstances(consumerID, providerID string) ([]*MicroServiceInstance, error)
	FindMicroServiceInstances(consumerID, appID, microServiceName, version string) ([]*MicroServiceInstance, error)
	UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error
	UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error
	UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error
	UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error
	String() string
	AutoSync()
	AddSchemas(microServiceID, schemaName, schemaInfo string) error
	GetMicroServicesByInterface(interfaceName string) (microservices []*MicroService)
	GetSchemaContentByInterface(interfaceName string) SchemaContent
	GetSchemaContentByServiceName(svcName, version, appID, env string) []*SchemaContent
}

// enableRegistry enable registry
func enableRegistry(opts ...Option) {
	rt := config.GlobalDefinition.Cse.Service.Registry.Type
	if rt == "" {
		rt = DefaultRegistryPlugin
	}
	f := registryFunc[rt]
	// TODO check whether the registry exists.
	if f == nil {
		panic("No registry plugin")
	}
	RegistryService = f(opts...)
	lager.Logger.Warnf(nil, "Enable %s registry.", RegistryService)
}

// InstallPlugin install plugin
func InstallPlugin(name string, f func(opts ...Option) Registry) {
	registryFunc[name] = f
}

// Enable create RegistryService
func Enable() error {
	mu.Lock()
	defer mu.Unlock()
	if IsEnabled {
		return nil
	}

	if config.GlobalDefinition.Cse.Service.Registry.Type != common.FileRegistry {
		if config.GlobalDefinition.Cse.Service.Registry.Address == "" {
			config.GlobalDefinition.Cse.Service.Registry.Address = DefaultAddr
		}
	}

	if config.GlobalDefinition.Cse.Service.Registry.Tenant == "" {
		config.GlobalDefinition.Cse.Service.Registry.Tenant = common.DefaultApp
	}

	var tlsConfig *tls.Config
	var scheme string
	addrs := strings.Split(config.GlobalDefinition.Cse.Service.Registry.Address, ",")
	hosts := make([]string, len(addrs))

	for index, addr := range addrs {
		if strings.Contains(addr, "://") {
			u, e := url.Parse(addr)
			if e != nil {
				return e
			}
			if len(scheme) != 0 && u.Scheme != scheme {
				return fmt.Errorf("inconsistent scheme found in registry address")
			}
			scheme = u.Scheme
			hosts[index] = u.Host
		}
	}

	secure := (scheme == common.HTTPS)
	if secure {
		sslTag := Name + "." + common.Consumer
		tmpTLSConfig, sslConfig, err := chassisTLS.GetTLSConfigByService(Name, "", common.Consumer)
		if err != nil {
			if chassisTLS.IsSSLConfigNotExist(err) {
				tmpErr := fmt.Errorf("%s tls mode, but no ssl config", sslTag)
				lager.Logger.Error(tmpErr.Error(), err)
				return tmpErr
			}
			lager.Logger.Errorf(err, "Load %s TLS config failed.", sslTag)
			return err
		}
		lager.Logger.Warnf(nil, "%s TLS mode, verify peer: %t, cipher plugin: %s.",
			sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin)
		tlsConfig = tmpTLSConfig
	}

	if config.GlobalDefinition.Cse.Service.Registry.Type == common.File {
		hosts = append(hosts, config.GlobalDefinition.Cse.Service.Registry.Address)
	}

	enableRegistry(
		Addrs(hosts...),
		Tenant(config.GlobalDefinition.Cse.Service.Registry.Tenant),
		EnableSSL(secure),
		TLSConfig(tlsConfig),
		Version(config.GlobalDefinition.Cse.Service.Registry.APIVersion.Version))
	if err := RegisterMicroservice(); err != nil {
		return err
	}
	go HBService.Start()
	RegistryService.AutoSync()
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
			return err
		}
	}
	return nil
}
