package servicecenter

import (
	scregistry "github.com/apache/servicecomb-service-center/pkg/registry"
	"net/url"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/scclient"

	"github.com/go-mesh/openlogging"
	"k8s.io/apimachinery/pkg/util/sets"
)

// constant values for default expiration time, and refresh interval
const (
	DefaultExpireTime      = 0
	DefaultRefreshInterval = time.Second * 30
)

// constant values for checking instance ID status
const (
	InstanceIDIsExist    = "instanceIdIsExist"
	InstanceIDIsNotExist = "instanceIdIsNotExist"
)

// CacheManager cache manager
type CacheManager struct {
	registryClient *client.RegistryClient
}

// AutoSync automatically sync the running instances
func (c *CacheManager) AutoSync() {
	c.refreshCache()
	if config.GetServiceDiscoveryWatch() {
		err := c.registryClient.WatchMicroService(runtime.ServiceID, watch)
		if err != nil {
			openlogging.GetLogger().Errorf("Watch failed. Self Micro service Id:%s. %s", runtime.ServiceID, err)
		}
		openlogging.GetLogger().Debugf("Watching Instances change events.")
	}
	var ticker *time.Ticker
	refreshInterval := config.GetServiceDiscoveryRefreshInterval()
	if refreshInterval == "" {
		ticker = time.NewTicker(DefaultRefreshInterval)
	} else {
		timeValue, err := time.ParseDuration(refreshInterval)
		if err != nil {
			openlogging.GetLogger().Errorf("refeshInterval is invalid. So use Default value, err %s", err)
			timeValue = DefaultRefreshInterval
		}
		ticker = time.NewTicker(timeValue)
	}
	go func() {
		for range ticker.C {
			c.refreshCache()
		}
	}()
}

// refreshCache refresh cache
func (c *CacheManager) refreshCache() {
	if archaius.GetBool("servicecomb.registry.autodiscovery", false) {
		err := c.registryClient.SyncEndpoints()
		if err != nil {
			openlogging.GetLogger().Errorf("get sc endpoints failed: %s", err)
		}
	}
	err := c.pullMicroServiceInstance()
	if err != nil {
		openlogging.GetLogger().Errorf("AutoUpdateMicroserviceInstance failed: %s", err)
		//connection with sc may lost, reset the revision
		c.registryClient.ResetRevision()
	}

	if archaius.GetBool("servicecomb.registry.autoSchemaIndex", false) {
		err = c.MakeSchemaIndex()
		if err != nil {
			openlogging.GetLogger().Errorf("MakeSchemaIndex failed: %s", err)
		}
	}

	if archaius.GetBool("servicecomb.registry.autoIPIndex", false) {
		err = c.MakeIPIndex()
		if err != nil {
			openlogging.GetLogger().Errorf("Auto Update IP index failed: %s", err)
		}
	}

}

// MakeIPIndex make ip index
// if store instance metadata into tags
// it will be used in route management
func (c *CacheManager) MakeIPIndex() error {
	openlogging.GetLogger().Debug("Make IP index")
	services, err := c.registryClient.GetAllResources("instances")
	if err != nil {
		openlogging.GetLogger().Errorf("Get instances failed: %s", err)
		return err
	}
	for _, service := range services {
		for _, inst := range service.Instances {
			for _, uri := range inst.Endpoints {
				u, err := url.Parse(uri)
				if err != nil {
					openlogging.GetLogger().Errorf("Wrong URI %s: %s", uri, err)
					continue
				}
				si := &registry.SourceInfo{}
				si.Tags = inst.Properties
				if si.Tags == nil {
					si.Tags = make(map[string]string)
				}
				si.Name = service.MicroService.ServiceName
				si.Tags[common.BuildinTagApp] = service.MicroService.AppId
				si.Tags[common.BuildinTagVersion] = service.MicroService.Version
				registry.SetIPIndex(u.Hostname(), si)
				//no need to analyze each endpoint, so break
				break
			}
		}
	}
	return nil
}

// MakeSchemaIndex make schema index
func (c *CacheManager) MakeSchemaIndex() error {

	openlogging.GetLogger().Debug("Make Schema index")
	microServiceList, err := c.registryClient.GetAllMicroServices()
	if err != nil {
		openlogging.GetLogger().Errorf("Get instances failed: %s", err)
		return err
	}

	for _, ms := range microServiceList {
		serviceID, err := c.registryClient.GetMicroServiceID(ms.AppId, ms.ServiceName, ms.Version, ms.Environment)
		if err != nil {
			continue
		}

		for _, schemaName := range ms.Schemas {

			content, err := c.registryClient.GetSchema(serviceID, schemaName)
			if err != nil {
				continue
			}

			schemaContent, err := parseSchemaContent(content)
			if err != nil {
				continue
			}

			interfaceName := schemaContent.Info["x-java-interface"]
			if interfaceName != "" {
				value, ok := registry.SchemaInterfaceIndexedCache.Get(interfaceName)
				if !ok {
					var allMicroServices []*scregistry.MicroService
					allMicroServices = append(allMicroServices, ms)
					registry.SchemaInterfaceIndexedCache.Set(interfaceName, allMicroServices, 0)
					openlogging.GetLogger().Debugf("New Interface added in the Index Cache : %s", interfaceName)
				} else {
					val, _ := value.([]*scregistry.MicroService)
					if !checkIfMicroServiceExistInList(val, ms.ServiceId) {
						val = append(val, ms)
						registry.SchemaInterfaceIndexedCache.Set(interfaceName, val, 0)
						openlogging.GetLogger().Debugf("New Interface added in the Index Cache : %s", interfaceName)
					}
				}

				svcValue, ok := registry.SchemaServiceIndexedCache.Get(serviceID)
				if !ok {
					var allMicroServices []*scregistry.MicroService
					allMicroServices = append(allMicroServices, ms)
					registry.SchemaServiceIndexedCache.Set(serviceID, allMicroServices, 0)
					openlogging.GetLogger().Debugf("New Service added in the Index Cache : %s", serviceID)
				} else {
					val, _ := svcValue.([]*scregistry.MicroService)
					if !checkIfMicroServiceExistInList(val, ms.ServiceId) {
						val = append(val, ms)
						registry.SchemaServiceIndexedCache.Set(serviceID, val, 0)
						openlogging.GetLogger().Debugf("New Service added in the Index Cache : %s", serviceID)
					}
				}
			}
		}
	}
	return nil
}

// This functions checks if the microservices exist in the list passed in argument
func checkIfMicroServiceExistInList(microserviceList []*scregistry.MicroService, serviceID string) bool {
	msIsPresentInList := false
	for _, interfaceMicroserviceList := range microserviceList {
		if interfaceMicroserviceList.ServiceId == serviceID {
			msIsPresentInList = true
			break
		}
	}
	return msIsPresentInList
}

// pullMicroServiceInstance pull micro-service instance
func (c *CacheManager) pullMicroServiceInstance() error {
	//Get Providers
	services := GetCriteria()
	serviceNameSet, _ := getServiceSet(services)
	c.compareAndDeleteOutdatedProviders(serviceNameSet)
	if len(services) == 0 {
		openlogging.Info("no providers")
		return nil
	}
	//fetch remote based on app and service
	response, err := c.registryClient.BatchFindInstances(runtime.ServiceID, services)
	if err != nil {
		if err == client.ErrNotModified || err == client.ErrEmptyCriteria {
			openlogging.Debug(err.Error())
		} else {
			openlogging.Error("Refresh local instance cache failed: " + err.Error())
		}
	}
	instances := RegroupInstances(services, response)
	filter(instances)

	return nil
}

func (c *CacheManager) compareAndDeleteOutdatedProviders(newProviders sets.String) {
	oldProviders := registry.MicroserviceInstanceIndex.FullCache().Items()
	for old := range oldProviders {
		if !newProviders.Has(old) { //provider is outdated, delete it
			registry.MicroserviceInstanceIndex.Delete(old)
			openlogging.GetLogger().Infof("Delete the service [%s] in the cache", old)
		}
	}
}

// getServiceSet regroup the providers by service name
func getServiceSet(exist []*scregistry.FindService) (sets.String, map[string]sets.String) {
	//get Provider's instances
	serviceNameSet := sets.NewString()                     // key is serviceName
	serviceNameAppIDKeySet := make(map[string]sets.String) // key is "serviceName" value is app sets
	if len(exist) == 0 {
		return serviceNameSet, serviceNameAppIDKeySet
	}
	for _, service := range exist {
		if service == nil {
			openlogging.Warn("FindService info is empty")
			continue
		}
		if service.Service == nil {
			openlogging.Warn("provider info is empty")
			continue
		}
		serviceNameSet.Insert(service.Service.ServiceName)
		m, ok := serviceNameAppIDKeySet[service.Service.ServiceName]
		if ok {
			m.Insert(service.Service.AppId)
		} else {
			serviceNameAppIDKeySet[service.Service.ServiceName] = sets.NewString()
			serviceNameAppIDKeySet[service.Service.ServiceName].Insert(service.Service.AppId)
		}
	}
	return serviceNameSet, serviceNameAppIDKeySet
}

//set app into instance metadata, split instances into ups and downs
//set instance to cache by service name
func filter(providerInstances map[string][]*registry.MicroServiceInstance) {
	//append instances from different app and same service name into one unified slice
	downs := make(map[string]struct{})
	for serviceName, instances := range providerInstances {
		up := make([]*registry.MicroServiceInstance, 0)
		for _, ins := range instances {
			switch {
			case ins.Version == "":
				openlogging.Warn("do not support old service center, plz upgrade")
				continue
			case ins.Status != common.DefaultStatus && ins.Status != common.TESTINGStatus:
				downs[ins.InstanceID] = struct{}{}
				openlogging.GetLogger().Debugf("do not cache the instance in '%s' status, instanceId = %s/%s",
					ins.Status, ins.ServiceID, ins.InstanceID)
				continue
			default:
				up = append(up, ins.WithAppID(ins.App))
			}
		}
		registry.RefreshCache(serviceName, up, downs) //save cache after get all instances of a service name
	}

}

// watch watching micro-service instance status
func watch(response *client.MicroServiceInstanceChangedEvent) {
	if response.Instance.Status != client.MSInstanceUP {
		response.Action = common.Delete
	}
	switch response.Action {
	case client.EventCreate:
		createAction(response)
	case client.EventDelete:
		deleteAction(response)
	case client.EventUpdate:
		updateAction(response)
	case client.EventError:
		openlogging.GetLogger().Warnf("MicroServiceInstanceChangedEvent action is error, MicroServiceInstanceChangedEvent = %s", response)
	default:
		openlogging.GetLogger().Warnf("Do not support this Action = %s", response.Action)
		return
	}
}

// createAction added micro-service instance to the cache
func createAction(response *client.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	microServiceInstances, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		openlogging.GetLogger().Errorf("ServiceID does not exist in MicroServiceInstanceCache,action is EVT_CREATE.key = %s", key)
		return
	}
	if response.Instance.Status != client.MSInstanceUP {
		openlogging.GetLogger().Warnf("createAction failed,MicroServiceInstance status is not MSI_UP,MicroServiceInstanceChangedEvent = %s", response)
		return
	}
	msi := ToMicroServiceInstance(response.Instance).WithAppID(response.Key.AppId)
	microServiceInstances = append(microServiceInstances, msi)
	registry.MicroserviceInstanceIndex.Set(key, microServiceInstances)
	openlogging.GetLogger().Infof("Cached Instances,action is EVT_CREATE, sid = %s, instances length = %d", response.Instance.ServiceId, len(microServiceInstances))
}

// deleteAction delete micro-service instance
func deleteAction(response *client.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	openlogging.GetLogger().Debugf("Received event EVT_DELETE, sid = %s, endpoints = %s", response.Instance.ServiceId, response.Instance.Endpoints)
	if err := registry.HealthCheck(key, response.Key.Version, response.Key.AppId, ToMicroServiceInstance(response.Instance)); err == nil {
		return
	}
	microServiceInstances, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		openlogging.GetLogger().Errorf("ServiceID does not exist in MicroserviceInstanceCache, action is EVT_DELETE, key = %s", key)
		return
	}
	var newInstances = make([]*registry.MicroServiceInstance, 0)
	for _, v := range microServiceInstances {
		if v.InstanceID != response.Instance.InstanceId {
			newInstances = append(newInstances, v)
		}
	}

	registry.MicroserviceInstanceIndex.Set(key, newInstances)
	openlogging.GetLogger().Debugf("Cached [%d] Instances of service [%s]", len(newInstances), key)
}

// updateAction update micro-service instance event
func updateAction(response *client.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	microServiceInstances, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		openlogging.GetLogger().Errorf("ServiceID does not exist in MicroserviceInstanceCache, action is EVT_UPDATE, sid = %s", key)
		return
	}
	if response.Instance.Status != client.MSInstanceUP {
		openlogging.GetLogger().Warnf("updateAction failed, MicroServiceInstance status is not MSI_UP, MicroServiceInstanceChangedEvent = %s", response)
		return
	}
	msi := ToMicroServiceInstance(response.Instance).WithAppID(response.Key.AppId)
	var iidExist = InstanceIDIsNotExist
	var arrayNum int
	for k, v := range microServiceInstances {
		if v.InstanceID == response.Instance.InstanceId {
			iidExist = InstanceIDIsExist
			arrayNum = k
		}
	}
	switch iidExist {
	case InstanceIDIsExist:
		microServiceInstances[arrayNum] = msi
	case InstanceIDIsNotExist:
		microServiceInstances = append(microServiceInstances, msi)
	default:
		openlogging.GetLogger().Warnf("updateAction error, iid:%s", response.Instance.InstanceId)
	}
	registry.MicroserviceInstanceIndex.Set(key, microServiceInstances)
	openlogging.GetLogger().Infof("Cached Instances,action is EVT_UPDATE, sid = %s, instances length = %d", response.Instance.ServiceId, len(microServiceInstances))
}
