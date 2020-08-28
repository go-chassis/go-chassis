package servicecenter

import (
	"fmt"
	scregistry "github.com/apache/servicecomb-service-center/pkg/registry"

	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/scclient"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-chassis/openlog"
)

const (
	// ServiceCenter constant string
	ServiceCenter = "servicecenter"
)

// Registrator to represent the object of service center to call the APIs of service center
type Registrator struct {
	Name           string
	registryClient *client.RegistryClient
	opts           client.Options
}

// RegisterService : 注册微服务
func (r *Registrator) RegisterService(ms *registry.MicroService) (string, error) {
	serviceKey := registry.Microservice2ServiceKeyStr(ms)
	microservice := ToSCService(ms)
	sid, err := r.registryClient.GetMicroServiceID(microservice.AppId, microservice.ServiceName, microservice.Version, microservice.Environment)
	if err != nil {
		openlog.Error(fmt.Sprintf("Get service [%s] failed, err %s", serviceKey, err))
		return "", err
	}
	if sid == "" {
		openlog.Warn(fmt.Sprintf("service [%s] not exists in registry, register it", serviceKey))
		sid, err = r.registryClient.RegisterService(microservice)
		if err != nil {
			openlog.Error(fmt.Sprintf("register service [%s] failed, err: %s", serviceKey, err))
			return "", err
		}
	} else {
		openlog.Info(fmt.Sprintf("[%s] exists in registry", serviceKey))
	}

	return sid, nil
}

// RegisterServiceInstance : 注册微服务
func (r *Registrator) RegisterServiceInstance(sid string, cIns *registry.MicroServiceInstance) (string, error) {
	instance := ToSCInstance(cIns)
	instance.ServiceId = sid
	instanceID, err := r.registryClient.RegisterMicroServiceInstance(instance)
	if err != nil {
		openlog.Error("register instance failed.")
		return "", err
	}
	return instanceID, nil
}

// RegisterServiceAndInstance : 注册微服务
func (r *Registrator) RegisterServiceAndInstance(cMicroService *registry.MicroService, cInstance *registry.MicroServiceInstance) (string, string, error) {
	microService := ToSCService(cMicroService)
	instance := ToSCInstance(cInstance)
	microServiceID, err := r.registryClient.GetMicroServiceID(microService.AppId,
		microService.ServiceName, microService.Version, microService.Environment)
	if err != nil {
		return "", "", err
	}
	if microServiceID == "" {
		microServiceID, err = r.registryClient.RegisterService(microService)
		if err != nil {
			openlog.Error(fmt.Sprintf("RegisterMicroService failed %s", err))
			return "", "", err
		}
		openlog.Debug(fmt.Sprintf("RegisterMicroService success, microServiceID: %s", microServiceID))
	}
	instance.ServiceId = microServiceID
	instanceID, err := r.registryClient.RegisterMicroServiceInstance(instance)
	if err != nil {
		openlog.Error(fmt.Sprintf("register instance failed. %s", err))
		return microServiceID, "", err
	}
	openlog.Info(fmt.Sprintf("register instance success, microServiceID/instanceID: %s/%s.", microServiceID, instanceID))
	return microServiceID, instanceID, nil
}

// UnRegisterMicroServiceInstance : 去注册微服务实例
func (r *Registrator) UnRegisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error {
	isSuccess, err := r.registryClient.UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID)
	if !isSuccess || err != nil {
		openlog.Error(fmt.Sprintf("unregister instance failed, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID))
		return err
	}
	openlog.Debug(fmt.Sprintf("unregister instance success, microServiceID/instanceID = %s/%s.", microServiceID, microServiceInstanceID))
	return nil
}

// Heartbeat : Keep instance heartbeats.
func (r *Registrator) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	bo, err := r.registryClient.Heartbeat(microServiceID, microServiceInstanceID)
	if err != nil {
		openlog.Error(fmt.Sprintf("Heartbeat failed, microServiceID/instanceID: %s/%s. %s",
			microServiceID, microServiceInstanceID, err))
		return false, err
	}
	if !bo {
		openlog.Error(fmt.Sprintf("Heartbeat failed, microServiceID/instanceID: %s/%s. %s",
			microServiceID, microServiceInstanceID, err))
		return bo, err
	}
	openlog.Debug(fmt.Sprintf("heartbeat success, microServiceID/instanceID: %s/%s.",
		microServiceID, microServiceInstanceID))
	return bo, nil
}

// AddSchemas to service center
func (r *Registrator) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	if err := r.registryClient.AddSchemas(microServiceID, schemaName, schemaInfo); err != nil {
		openlog.Error(fmt.Sprintf("add schemas failed: %s", err))
		return err
	}
	openlog.Info("add all schemas success.")
	return nil
}

// UpdateMicroServiceInstanceStatus : 更新微服务实例状态信息
func (r *Registrator) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error {
	isSuccess, err := r.registryClient.UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status)
	if !isSuccess {
		openlog.Error(fmt.Sprintf("UpdateMicroServiceInstanceStatus failed, microServiceID/instanceID = %s/%s.",
			microServiceID, microServiceInstanceID))
		return err
	}
	openlog.Debug(fmt.Sprintf("UpdateMicroServiceInstanceStatus success, microServiceID/instanceID = %s/%s.",
		microServiceID, microServiceInstanceID))
	return nil
}

// UpdateMicroServiceProperties 更新微服务properties信息
func (r *Registrator) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	microService := &scregistry.MicroService{
		Properties: properties,
	}
	isSuccess, err := r.registryClient.UpdateMicroServiceProperties(microServiceID, microService)
	if !isSuccess {
		openlog.Error(
			fmt.Sprintf("UpdateMicroService Properties failed, microServiceID/instanceID = %s.", microServiceID))
		return err
	}
	openlog.Debug(fmt.Sprintf("UpdateMicroService Properties success, microServiceID/instanceID = %s.",
		microServiceID))
	return nil
}

// UpdateMicroServiceInstanceProperties : 更新微服务实例properties信息
func (r *Registrator) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error {
	microServiceInstance := &scregistry.MicroServiceInstance{
		Properties: properties,
	}
	isSuccess, err := r.registryClient.UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID, microServiceInstance)
	if !isSuccess {
		openlog.Error(fmt.Sprintf("UpdateMicroServiceInstanceProperties failed, microServiceID/instanceID = %s/%s.",
			microServiceID, microServiceInstanceID))
		return err
	}
	openlog.Debug(fmt.Sprintf("UpdateMicroServiceInstanceProperties success, microServiceID/instanceID = %s/%s.",
		microServiceID, microServiceInstanceID))
	return nil
}

// Close : Close all client connection.
func (r *Registrator) Close() error {
	return closeClient(r.registryClient)
}

// ServiceDiscovery to represent the object of service center to call the APIs of service center
type ServiceDiscovery struct {
	Name           string
	registryClient *client.RegistryClient
	opts           client.Options
}

// GetAllMicroServices : Get all MicroService information.
func (r *ServiceDiscovery) GetAllMicroServices() ([]*registry.MicroService, error) {
	microServices, err := r.registryClient.GetAllMicroServices()
	if err != nil {
		openlog.Error(fmt.Sprintf("GetAllMicroServices failed: %s", err))
		return nil, err
	}
	mss := make([]*registry.MicroService, 0)
	for _, s := range microServices {
		mss = append(mss, ToMicroService(s))
	}
	openlog.Debug(fmt.Sprintf("GetAllMicroServices success, MicroService: %v", microServices))
	return mss, nil
}

// GetAllApplications : Get all Applications information.
func (r *ServiceDiscovery) GetAllApplications() ([]string, error) {
	apps, err := r.registryClient.GetAllApplications()
	if err != nil {
		openlog.Error(fmt.Sprintf("GetAllApplications failed: %s", err))
		return nil, err
	}
	openlog.Debug(fmt.Sprintf("GetAllApplications success, Applications: %s", apps))
	return apps, nil
}

// GetMicroService : 根据microServiceID获取对应的微服务信息
func (r *ServiceDiscovery) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	microService, err := r.registryClient.GetMicroService(microServiceID)
	if err != nil {
		openlog.Error(fmt.Sprintf("GetMicroService failed: %s", err))
		return nil, err
	}
	openlog.Debug(fmt.Sprintf("GetMicroServices success, MicroService: %v", microService))
	return ToMicroService(microService), nil
}

// GetMicroServiceInstances : 获取指定微服务的所有实例
func (r *ServiceDiscovery) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	providerInstances, err := r.registryClient.GetMicroServiceInstances(consumerID, providerID)
	if err != nil {
		openlog.Error(fmt.Sprintf("GetMicroServiceInstances failed: %s", err))
		return nil, err
	}
	instances := filterInstances(providerInstances)
	openlog.Debug(fmt.Sprintf("GetMicroServiceInstances success, consumerID/providerID: %s/%s", consumerID, providerID))
	return instances, nil
}

// FindMicroServiceInstances find micro-service instances
func (r *ServiceDiscovery) FindMicroServiceInstances(consumerID, microServiceName string, tags utiltags.Tags) ([]*registry.MicroServiceInstance, error) {
	// TODO: wrap default tags for service center
	// because sc need version and appID to generate tags
	tags = wrapTagsForServiceCenter(tags)
	microServiceInstance, boo := registry.MicroserviceInstanceIndex.Get(microServiceName, tags.KV)
	appID := tags.AppID()
	if appID == "" {
		appID = runtime.App
	}
	registry.AddProviderToCache(microServiceName, appID)
	if !boo || microServiceInstance == nil {
		criteria := GetCriteriaByService(microServiceName)
		openlog.Warn(fmt.Sprintf("%s Get instances from remote, key: %s:%s:%s", consumerID, appID, microServiceName, tags.Version()))
		providerInstancesResponse, err := r.registryClient.BatchFindInstances(consumerID, criteria)
		if err != nil {
			return nil, fmt.Errorf("FindMicroServiceInstances failed, ProviderID: %s, err: %s", microServiceName, err)
		}
		providerInstances := RegroupInstances(criteria, providerInstancesResponse)
		filter(providerInstances)
		microServiceInstance, boo = registry.MicroserviceInstanceIndex.Get(microServiceName, tags.KV)
		if !boo || microServiceInstance == nil {
			openlog.Debug(fmt.Sprintf("Find no micro service instances for %s from cache", microServiceName))
			return nil, nil
		}
	}

	return microServiceInstance, nil
}

//RegroupInstances organize raw data to better format
func RegroupInstances(keys []*scregistry.FindService, response *scregistry.BatchFindInstancesResponse) map[string][]*registry.MicroServiceInstance {
	instanceMap := make(map[string][]*registry.MicroServiceInstance)
	if response.Services != nil {
		for _, result := range response.Services.Updated {
			if len(result.Instances) == 0 {
				continue
			}
			for _, instance := range result.Instances {
				ni := ToMicroServiceInstance(instance)
				ni.App = keys[result.Index].Service.AppId
				ni.Version = keys[result.Index].Service.Version
				instances, ok := instanceMap[keys[result.Index].Service.ServiceName]
				if !ok {
					instances = []*registry.MicroServiceInstance{ni}
					instanceMap[keys[result.Index].Service.ServiceName] = instances
					continue
				}
				instanceMap[keys[result.Index].Service.ServiceName] = append(instances, ni)
			}

		}
	}
	return instanceMap
}

// WatchMicroService : 支持用户自调用主动监听实例变化功能
func (r *ServiceDiscovery) WatchMicroService(selfMicroServiceID string, callback func(*client.MicroServiceInstanceChangedEvent)) {
	err := r.registryClient.WatchMicroService(selfMicroServiceID, callback)
	if err != nil {
		openlog.Error(err.Error())
	}

}

// AutoSync updating the cache manager
func (r *ServiceDiscovery) AutoSync() {
	c := &CacheManager{
		registryClient: r.registryClient,
	}
	c.AutoSync()
}

// Close : Close all websocket connection.
func (r *ServiceDiscovery) Close() error {
	return closeClient(r.registryClient)
}

// ContractDiscovery to represent the object of service center to call the APIs of service center
type ContractDiscovery struct {
	Name           string
	registryClient *client.RegistryClient
	opts           client.Options
}

// GetMicroServicesByInterface get micro-services by interface
func (r *ContractDiscovery) GetMicroServicesByInterface(interfaceName string) (microService []*registry.MicroService) {
	value, ok := registry.SchemaInterfaceIndexedCache.Get(interfaceName)
	if !ok || value == nil {
		r.fillSchemaInterfaceIndexCache(nil, interfaceName)
		value, _ = registry.SchemaInterfaceIndexedCache.Get(interfaceName)
	}

	microServiceModel, ok := value.([]*scregistry.MicroService)

	if !ok {
		openlog.Error("GetMicroServicesByInterface failed, Type asserts failed")
	}

	for _, v := range microServiceModel {
		ms := ToMicroService(v)
		microService = append(microService, ms)
	}
	return microService
}

// GetSchemaContentByInterface get schema content by interface
func (r *ContractDiscovery) GetSchemaContentByInterface(interfaceName string) (schemas registry.SchemaContent) {
	value, ok := registry.SchemaInterfaceIndexedCache.Get(interfaceName)
	if !ok || value == nil {
		return r.fillSchemaInterfaceIndexCache(nil, interfaceName)
	}

	val, ok := value.([]*scregistry.MicroService)

	if !ok {
		return schemas
	}

	return r.fillSchemaInterfaceIndexCache(val, interfaceName)

}

// GetSchemaContentByServiceName get schema content by service name
func (r *ContractDiscovery) GetSchemaContentByServiceName(svcName, version, appID, env string) (schemas []*registry.SchemaContent) {
	serviceID, err := r.registryClient.GetMicroServiceID(appID, svcName, version, env)
	if err != nil {
		return schemas
	}
	value, ok := registry.SchemaServiceIndexedCache.Get(serviceID)
	if !ok || value == nil {
		return r.fillSchemaServiceIndexCache(nil, serviceID)
	}

	val, ok := value.([]*scregistry.MicroService)

	if !ok {
		return schemas
	}

	return r.fillSchemaServiceIndexCache(val, serviceID)

}

// fillSchemaServiceIndexCache fill schema service index cache
func (r *ContractDiscovery) fillSchemaServiceIndexCache(ms []*scregistry.MicroService, serviceID string) (content []*registry.SchemaContent) {
	if ms == nil {
		microServiceList, err := r.registryClient.GetAllMicroServices()
		if err != nil {
			openlog.Error(fmt.Sprintf("Get instances failed: %s", err))
			return content
		}

		return r.fillCacheAndGetServiceSchemaContent(microServiceList, serviceID)
	}

	return r.fillCacheAndGetServiceSchemaContent(ms, serviceID)
}

// fillCacheAndGetServiceSchemaContent fill cache and get services schema content
func (r *ContractDiscovery) fillCacheAndGetServiceSchemaContent(microServiceList []*scregistry.MicroService, serviceID string) (schemaContent []*registry.SchemaContent) {

	for _, ms := range microServiceList {

		if ms.ServiceId == serviceID {

			for _, schemaName := range ms.Schemas {
				content, err := r.registryClient.GetSchema(ms.ServiceId, schemaName)
				if err != nil {
					continue
				}

				schema, err := unmarshalSchemaContent(content)
				if err != nil {
					continue
				}
				_, ok := registry.SchemaServiceIndexedCache.Get(serviceID)
				if !ok {
					var allServices []*scregistry.MicroService
					allServices = append(allServices, ms)
					registry.SchemaServiceIndexedCache.Set(serviceID, allServices, 0)
				}
				schemaContent = append(schemaContent, schema)
			}
		}
	}

	return
}

// fillSchemaInterfaceIndexCache fill schema interface index cache
func (r *ContractDiscovery) fillSchemaInterfaceIndexCache(ms []*scregistry.MicroService, interfaceName string) (content registry.SchemaContent) {
	if ms == nil {
		microServiceList, err := r.registryClient.GetAllMicroServices()
		if err != nil {
			openlog.Error(fmt.Sprintf("Get instances failed: %s", err))
			return content
		}

		return r.fillCacheAndGetInterfaceSchemaContent(microServiceList, interfaceName)
	}

	return r.fillCacheAndGetInterfaceSchemaContent(ms, interfaceName)
}

// fillCacheAndGetInterfaceSchemaContent fill cache and get interface schema content
func (r *ContractDiscovery) fillCacheAndGetInterfaceSchemaContent(microServiceList []*scregistry.MicroService, interfaceName string) (schemaContent registry.SchemaContent) {

	for _, ms := range microServiceList {
		serviceID, err := r.registryClient.GetMicroServiceID(ms.AppId, ms.ServiceName, ms.Version, ms.Environment)
		if err != nil {
			continue
		}

		for _, schemaName := range ms.Schemas {
			content, err := r.registryClient.GetSchema(serviceID, schemaName)
			if err != nil {
				continue
			}

			schemaContent, err = parseSchemaContent(content)
			if err != nil {
				continue
			}

			interfaceValue := schemaContent.Info["x-java-interface"]
			if interfaceValue == "" {
				continue
			}

			value, ok := registry.SchemaInterfaceIndexedCache.Get(interfaceName)
			if !ok {
				var allServices []*scregistry.MicroService
				allServices = append(allServices, ms)
				registry.SchemaInterfaceIndexedCache.Set(interfaceValue, allServices, 0)
			} else {
				val, _ := value.([]*scregistry.MicroService)
				val = append(val, ms)
				registry.SchemaInterfaceIndexedCache.Set(interfaceValue, val, 0)

			}

			if interfaceName == interfaceValue {
				return schemaContent
			}
		}
	}

	return
}

// GetSchema from service center
func (r *ContractDiscovery) GetSchema(microServiceID, schemaName string) ([]byte, error) {
	var schemaContent []byte
	var err error
	if schemaContent, err = r.registryClient.GetSchema(microServiceID, schemaName); err != nil {
		openlog.Error(fmt.Sprintf("GetSchema failed: %s", err))
		return []byte(""), err
	}
	openlog.Debug("GetSchema success.")
	return schemaContent, nil

}

//Close close client connection
func (r *ContractDiscovery) Close() error {
	return closeClient(r.registryClient)
}

//NewRegistrator new Service center registrator
func NewRegistrator(options registry.Options) registry.Registrator {
	sco := ToSCOptions(options)
	r := &client.RegistryClient{}
	if err := r.Initialize(sco); err != nil {
		openlog.Error(fmt.Sprintf("RegistryClient initialization failed, err %s", err))
	}

	return &Registrator{
		Name:           ServiceCenter,
		registryClient: r,
		opts:           sco,
	}
}

//NewServiceDiscovery new service center discovery
func NewServiceDiscovery(options registry.Options) registry.ServiceDiscovery {
	sco := ToSCOptions(options)
	r := &client.RegistryClient{}
	if err := r.Initialize(sco); err != nil {
		openlog.Error(fmt.Sprintf("RegistryClient initialization failed. %s", err))
	}

	return &ServiceDiscovery{
		Name:           ServiceCenter,
		registryClient: r,
		opts:           sco,
	}
}
func newContractDiscovery(options registry.Options) registry.ContractDiscovery {
	sco := ToSCOptions(options)
	r := &client.RegistryClient{}
	if err := r.Initialize(sco); err != nil {
		openlog.Error(fmt.Sprintf("RegistryClient initialization failed: %s", err))
	}

	return &ContractDiscovery{
		Name:           ServiceCenter,
		registryClient: r,
		opts:           sco,
	}
}

// init initialize the plugin of service center registry
func init() {
	registry.InstallRegistrator(ServiceCenter, NewRegistrator)
	registry.InstallServiceDiscovery(ServiceCenter, NewServiceDiscovery)
	registry.InstallContractDiscovery(ServiceCenter, newContractDiscovery)

}
