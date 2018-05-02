package pilot

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
)

// PilotPlugin is the constant string of the plugin name
const PilotPlugin = "pilot"

// Registrator is the struct to do service discovery from istio pilot server
type Registrator struct {
	Name           string
	registryClient *EnvoyDSClient
}

// Close : Close all connection.
func (r *Registrator) Close() error {
	return close(r.registryClient)
}

// ServiceDiscovery is the struct to do service discovery from istio pilot server
type ServiceDiscovery struct {
	Name           string
	registryClient *EnvoyDSClient
}

// GetMicroServiceID : 获取指定微服务的MicroServiceID
func (r *ServiceDiscovery) GetMicroServiceID(appID, microServiceName, version, env string) (string, error) {
	_, err := r.registryClient.GetServiceHosts(microServiceName)
	if err != nil {
		lager.Logger.Errorf(err, "GetMicroServiceID failed")
		return "", err
	}
	lager.Logger.Debugf("GetMicroServiceID success")
	return microServiceName, nil
}

// GetAllMicroServices : Get all MicroService information.
func (r *ServiceDiscovery) GetAllMicroServices() ([]*registry.MicroService, error) {
	svcs, err := r.registryClient.GetAllServices()
	if err != nil {
		lager.Logger.Errorf(err, "GetAllMicroServices failed")
		return nil, err
	}

	var mss []*registry.MicroService
	for _, s := range svcs {
		mss = append(mss, ToMicroService(s))
	}
	return mss, nil
}

// GetMicroService : 根据microServiceID获取对应的微服务信息
func (r *ServiceDiscovery) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	hs, err := r.registryClient.GetServiceHosts(microServiceID)
	if err != nil {
		lager.Logger.Errorf(err, "GetMicroServiceID failed")
		return nil, err
	}
	lager.Logger.Debugf("GetMicroServices success, MicroService: %s", microServiceID)
	return ToMicroService(&Service{
		ServiceKey: microServiceID,
		Hosts:      hs.Hosts,
	}), nil
}

// GetMicroServiceInstances : 获取指定微服务的所有实例
func (r *ServiceDiscovery) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	hs, err := r.registryClient.GetServiceHosts(providerID)
	if err != nil {
		lager.Logger.Errorf(err, "GetMicroServiceInstances failed.")
		return nil, err
	}
	instances := filterInstances(hs.Hosts)
	lager.Logger.Debugf("GetMicroServiceInstances success, consumerID/providerID: %s/%s", consumerID, providerID)
	return instances, nil
}

// FindMicroServiceInstances find micro-service instances
func (r *ServiceDiscovery) FindMicroServiceInstances(consumerID, appID, microServiceName, version, env string) ([]*registry.MicroServiceInstance, error) {
	value, boo := registry.MicroserviceInstanceCache.Get(microServiceName)
	if !boo || value == nil {
		lager.Logger.Warnf("%s Get instances from remote, key: %s", consumerID, microServiceName)
		hs, err := r.registryClient.GetServiceHosts(microServiceName)
		if err != nil {
			return nil, fmt.Errorf("FindMicroServiceInstances failed, ProviderID: %s, err: %s",
				microServiceName, err)
		}

		filterRestore(hs.Hosts, microServiceName)
		value, boo = registry.MicroserviceInstanceCache.Get(microServiceName)
		if !boo || value == nil {
			lager.Logger.Debugf("Find no microservice instances for %s from cache", microServiceName)
			return nil, nil
		}
	}
	microServiceInstance, ok := value.([]*registry.MicroServiceInstance)
	if !ok {
		lager.Logger.Errorf(nil, "FindMicroServiceInstances failed, Type asserts failed. consumerIDL: %s",
			consumerID)
	}
	return microServiceInstance, nil
}

// AutoSync updating the cache manager
func (r *ServiceDiscovery) AutoSync() {
	c := &CacheManager{
		registryClient: r.registryClient,
	}
	c.AutoSync()
}

// Close : Close all connection.
func (r *ServiceDiscovery) Close() error {
	return close(r.registryClient)
}

func newDiscoveryService(options registry.Options) registry.ServiceDiscovery {
	c := &EnvoyDSClient{}
	c.Initialize(Options{
		Addrs:     options.Addrs,
		TLSConfig: options.TLSConfig,
	})
	return &ServiceDiscovery{
		Name:           PilotPlugin,
		registryClient: c,
	}
}

// register pilot registry plugin when import this package
func init() {
	registry.InstallServiceDiscovery(PilotPlugin, newDiscoveryService)
}
