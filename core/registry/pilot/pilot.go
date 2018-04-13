package pilot

import (
	"errors"
	"fmt"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"strings"
)

// PilotPlugin is the constant string of the plugin name
const PilotPlugin = "pilot"

// Registrator is the struct to do service discovery from istio pilot server
type Registrator struct {
	Name           string
	registryClient *EnvoyDSClient
}

// RegisterService : 注册微服务
func (r *Registrator) RegisterService(ms *registry.MicroService) (string, error) {
	lager.Logger.Warnf(errors.New("Not support operation"),
		"RegisterService [%s] failed.", ms.ServiceName)
	return ms.ServiceName, nil
}

// RegisterServiceInstance : 注册微服务
func (r *Registrator) RegisterServiceInstance(sid string, cIns *registry.MicroServiceInstance) (string, error) {
	if len(cIns.EndpointsMap) == 0 {
		err := errors.New("Required EndpointsMap")
		lager.Logger.Errorf(err, "RegisterMicroServiceInstance failed, Mid: %s", sid)
		return "", err
	}

	ep := cIns.EndpointsMap[common.ProtocolRest]
	if len(ep) == 0 {
		err := errors.New("Only support protocol '" + common.ProtocolRest + "'")
		lager.Logger.Errorf(err, "RegisterMicroServiceInstance failed, Mid: %s", sid)
		return "", err
	}

	instanceID := strings.Replace(ep, ":", "_", 1)
	value, ok := registry.SelfInstancesCache.Get(sid)
	if !ok {
		lager.Logger.Warnf(nil, "RegisterMicroServiceInstance get SelfInstancesCache failed, Mid/Sid: %s/%s",
			sid, instanceID)
	}
	instanceIDs, ok := value.([]string)
	if !ok {
		lager.Logger.Warnf(nil, "RegisterMicroServiceInstance type asserts failed,  Mid/Sid: %s/%s",
			sid, instanceID)
	}
	var isRepeat bool
	for _, va := range instanceIDs {
		if va == instanceID {
			isRepeat = true
		}
	}
	if !isRepeat {
		instanceIDs = append(instanceIDs, instanceID)
	}
	registry.SelfInstancesCache.Set(sid, instanceIDs, 0)
	lager.Logger.Infof("RegisterMicroServiceInstance success, microServiceID/instanceID: %s/%s.",
		sid, instanceID)
	return instanceID, nil
}

// RegisterServiceAndInstance : 注册微服务
func (r *Registrator) RegisterServiceAndInstance(cMicroService *registry.MicroService, cInstance *registry.MicroServiceInstance) (string, string, error) {
	microServiceID, err := r.RegisterService(cMicroService)
	if err != nil {
		return "", "", err
	}
	instanceID, err := r.RegisterServiceInstance(microServiceID, cInstance)
	if err != nil {
		return "", "", err
	}
	return microServiceID, instanceID, nil
}

// Heartbeat : Keep instance heartbeats.
func (r *Registrator) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	lager.Logger.Debugf("Heartbeat success, microServiceID/instanceID: %s/%s.", microServiceID, microServiceInstanceID)
	return true, nil
}

// AddDependencies ： 注册微服务的依赖关系
func (r *Registrator) AddDependencies(cDep *registry.MicroServiceDependency) error {
	lager.Logger.Debugf("AddDependencies success.")
	return nil
}

// AddSchemas to service center
func (r *Registrator) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	lager.Logger.Debugf("AddSchemas success.")
	return nil
}

// UnRegisterMicroServiceInstance : 去注册微服务实例
func (r *Registrator) UnRegisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error {
	lager.Logger.Errorf(errors.New("Not support operation"),
		"unregisterMicroServiceInstance failed, microServiceID/instanceID = %s/%s.",
		microServiceID, microServiceInstanceID)
	return nil
}

// UpdateMicroServiceInstanceStatus : 更新微服务实例状态信息
func (r *Registrator) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error {
	lager.Logger.Debugf(
		"UpdateMicroServiceInstanceStatus failed, microServiceID/instanceID = %s/%s. error: Not support operation",
		microServiceID, microServiceInstanceID)
	return nil
}

// UpdateMicroServiceProperties 更新微服务properties信息
func (r *Registrator) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	lager.Logger.Debugf(
		"UpdateMicroService Properties failed, microServiceID/instanceID = %s. error: Not support operation",
		microServiceID)
	return nil
}

// UpdateMicroServiceInstanceProperties : 更新微服务实例properties信息
func (r *Registrator) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error {
	lager.Logger.Debugf(
		"UpdateMicroServiceInstanceProperties failed, microServiceID/instanceID = %s/%s. error: Not support operation",
		microServiceID, microServiceInstanceID)
	return nil
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
		lager.Logger.Warnf(nil, "%s Get instances from remote, key: %s", consumerID, microServiceName)
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
func newRegistrator(options registry.Options) registry.Registrator {
	c := &EnvoyDSClient{}
	c.Initialize(Options{
		Addrs:     options.Addrs,
		TLSConfig: options.TLSConfig,
	})
	return &Registrator{
		Name:           PilotPlugin,
		registryClient: c,
	}
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
	registry.InstallRegistrator(PilotPlugin, newRegistrator)
	registry.InstallServiceDiscovery(PilotPlugin, newDiscoveryService)
}
