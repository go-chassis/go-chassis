package servicecenter

import (
	"errors"
	"fmt"

	scregistry "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/go-chassis/v2/health"

	"net/url"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/openlog"
	"github.com/go-chassis/sc-client"
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
	registryClient *sc.Client
}

// AutoSync automatically sync the running instances
func (c *CacheManager) AutoSync() {
	c.refreshCache()
	if config.GetServiceDiscoveryWatch() {
		err := c.registryClient.WatchMicroService(runtime.ServiceID, watch)
		if err != nil {
			openlog.Error(fmt.Sprintf("watch failed. Self Micro service Id:%s. %s", runtime.ServiceID, err))
		}
		openlog.Debug("Watching Instances change events.")
	}
	var ticker *time.Ticker
	refreshInterval := config.GetServiceDiscoveryRefreshInterval()
	if refreshInterval == "" {
		ticker = time.NewTicker(DefaultRefreshInterval)
	} else {
		timeValue, err := time.ParseDuration(refreshInterval)
		if err != nil {
			openlog.Error(fmt.Sprintf("refeshInterval is invalid. So use Default value, err %s", err))
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
			openlog.Error(fmt.Sprintf("get sc endpoints failed: %s", err))
		}
	}
	err := c.pullMicroServiceInstance()
	if err != nil {
		openlog.Error(fmt.Sprintf("AutoUpdateMicroserviceInstance failed: %s", err))
		//connection with sc may lost, reset the revision
		c.registryClient.ResetRevision()
	}

	if archaius.GetBool("servicecomb.registry.autoSchemaIndex", false) {
		err = c.MakeSchemaIndex()
		if err != nil {
			openlog.Error(fmt.Sprintf("MakeSchemaIndex failed: %s", err))
		}
	}

	if archaius.GetBool("servicecomb.registry.autoIPIndex", false) {
		err = c.MakeIPIndex()
		if err != nil {
			openlog.Error(fmt.Sprintf("Auto Update IP index failed: %s", err))
		}
	}

}

// MakeIPIndex make ip index
// if store instance metadata into tags
// it will be used in route management
func (c *CacheManager) MakeIPIndex() error {
	openlog.Debug("Make IP index")
	services, err := c.registryClient.GetAllResources("instances")
	if err != nil {
		openlog.Error(fmt.Sprintf("Get instances failed: %s", err))
		return err
	}
	for _, service := range services {
		for _, inst := range service.Instances {
			for _, uri := range inst.Endpoints {
				u, err := url.Parse(uri)
				if err != nil {
					openlog.Error(fmt.Sprintf("Wrong URI %s: %s", uri, err))
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

	openlog.Debug("Make Schema index")
	microServiceList, err := c.registryClient.GetAllMicroServices()
	if err != nil {
		openlog.Error(fmt.Sprintf("Get instances failed: %s", err))
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
					openlog.Debug(fmt.Sprintf("New Interface added in the Index Cache : %s", interfaceName))
				} else {
					val, _ := value.([]*scregistry.MicroService)
					if !checkIfMicroServiceExistInList(val, ms.ServiceId) {
						val = append(val, ms)
						registry.SchemaInterfaceIndexedCache.Set(interfaceName, val, 0)
						openlog.Debug(fmt.Sprintf("New Interface added in the Index Cache : %s", interfaceName))
					}
				}

				svcValue, ok := registry.SchemaServiceIndexedCache.Get(serviceID)
				if !ok {
					var allMicroServices []*scregistry.MicroService
					allMicroServices = append(allMicroServices, ms)
					registry.SchemaServiceIndexedCache.Set(serviceID, allMicroServices, 0)
					openlog.Debug(fmt.Sprintf("New Service added in the Index Cache : %s", serviceID))
				} else {
					val, _ := svcValue.([]*scregistry.MicroService)
					if !checkIfMicroServiceExistInList(val, ms.ServiceId) {
						val = append(val, ms)
						registry.SchemaServiceIndexedCache.Set(serviceID, val, 0)
						openlog.Debug(fmt.Sprintf("New Service added in the Index Cache : %s", serviceID))
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
		openlog.Info("no providers")
		return nil
	}
	//fetch remote based on app and service
	response, err := c.registryClient.BatchFindInstances(runtime.ServiceID, services)
	if err != nil {
		if errors.Is(err, sc.ErrNotModified) || errors.Is(err, sc.ErrEmptyCriteria) {
			openlog.Debug(err.Error())
		} else {
			openlog.Error("Refresh local instance cache failed: " + err.Error())
			return fmt.Errorf("refresh local instance cache failed: %w", err)
		}
	}
	instances := RegroupInstances(services, response)
	filterAndCache(serviceNameSet, instances)

	return nil
}

func (c *CacheManager) compareAndDeleteOutdatedProviders(newProviders sets.String) {
	oldProviders := registry.MicroserviceInstanceIndex.FullCache().Items()
	for old := range oldProviders {
		if !newProviders.Has(old) { //provider is outdated, delete it
			registry.MicroserviceInstanceIndex.Delete(old)
			openlog.Info(fmt.Sprintf("Delete the service [%s] in the cache", old))
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
			openlog.Warn("FindService info is empty")
			continue
		}
		if service.Service == nil {
			openlog.Warn("provider info is empty")
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

// set app into instance metadata, split instances into ups and downs
// set instance to cache by service name
func filterAndCache(services sets.String, providerInstances map[string][]*registry.MicroServiceInstance) {
	//append instances from different app and same service name into one unified slice
	downs := make(map[string]struct{})
	if len(providerInstances) == 0 {
		openlog.Warn("can not find instance in service center, " +
			"but must set empty cache to avoid frequent call to service center")
		for service, _ := range services {
			registry.MicroserviceInstanceIndex.Set(service, make([]*registry.MicroServiceInstance, 0))
		}
	}

	for serviceName, instances := range providerInstances {
		up := make([]*registry.MicroServiceInstance, 0)
		for _, ins := range instances {
			switch {
			case ins.Version == "":
				openlog.Warn("do not support old service center, plz upgrade")
				continue
			case ins.Status != common.DefaultStatus && ins.Status != common.TESTINGStatus:
				downs[ins.InstanceID] = struct{}{}
				openlog.Debug(fmt.Sprintf("do not cache the instance in '%s' status, instanceId = %s/%s",
					ins.Status, ins.ServiceID, ins.InstanceID))
				continue
			default:
				up = append(up, ins.WithAppID(ins.App))
			}
		}
		health.RefreshCache(serviceName, up, downs) //save cache after get all instances of a service name
	}

}

// watch watching micro-service instance status
func watch(response *sc.MicroServiceInstanceChangedEvent) {
	if response.Instance.Status != sc.MSInstanceUP {
		response.Action = common.Delete
	}
	switch response.Action {
	case sc.EventCreate:
		createAction(response)
	case sc.EventDelete:
		deleteAction(response)
	case sc.EventUpdate:
		updateAction(response)
	case sc.EventError:
		openlog.Warn(fmt.Sprintf("MicroServiceInstanceChangedEvent action is error, MicroServiceInstanceChangedEvent = %v", response))
	default:
		openlog.Warn(fmt.Sprintf("Do not support this Action = %s", response.Action))
		return
	}
}

// createAction added micro-service instance to the cache
func createAction(response *sc.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	microServiceInstances, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		openlog.Error(fmt.Sprintf("ServiceID does not exist in MicroServiceInstanceCache,action is EVT_CREATE.key = %s", key))
		return
	}
	if response.Instance.Status != sc.MSInstanceUP {
		openlog.Warn(fmt.Sprintf("createAction failed,MicroServiceInstance status is not MSI_UP,MicroServiceInstanceChangedEvent = %v", response))
		return
	}
	msi := ToMicroServiceInstance(response.Instance).WithAppID(response.Key.AppId)
	microServiceInstances = append(microServiceInstances, msi)
	registry.MicroserviceInstanceIndex.Set(key, microServiceInstances)
	openlog.Info(fmt.Sprintf("Cached Instances,action is EVT_CREATE, sid = %s, instances length = %d", response.Instance.ServiceId, len(microServiceInstances)))
}

// deleteAction delete micro-service instance
func deleteAction(response *sc.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	openlog.Debug(fmt.Sprintf("Received event EVT_DELETE, sid = %s, endpoints = %s", response.Instance.ServiceId, response.Instance.Endpoints))
	if err := health.HealthCheck(key, response.Key.Version, response.Key.AppId, ToMicroServiceInstance(response.Instance)); err == nil {
		return
	}
	microServiceInstances, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		openlog.Error(fmt.Sprintf("ServiceID does not exist in MicroserviceInstanceCache, action is EVT_DELETE, key = %s", key))
		return
	}
	var newInstances = make([]*registry.MicroServiceInstance, 0)
	for _, v := range microServiceInstances {
		if v.InstanceID != response.Instance.InstanceId {
			newInstances = append(newInstances, v)
		}
	}

	registry.MicroserviceInstanceIndex.Set(key, newInstances)
	openlog.Debug(fmt.Sprintf("Cached [%d] Instances of service [%s]", len(newInstances), key))
}

// updateAction update micro-service instance event
func updateAction(response *sc.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	microServiceInstances, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		openlog.Error(fmt.Sprintf("ServiceID does not exist in MicroserviceInstanceCache, action is EVT_UPDATE, sid = %s", key))
		return
	}
	if response.Instance.Status != sc.MSInstanceUP {
		openlog.Warn(fmt.Sprintf("updateAction failed, MicroServiceInstance status is not MSI_UP, MicroServiceInstanceChangedEvent = %v", response))
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
		openlog.Warn(fmt.Sprintf("updateAction error, iid:%s", response.Instance.InstanceId))
	}
	registry.MicroserviceInstanceIndex.Set(key, microServiceInstances)
	openlog.Info(fmt.Sprintf("Cached Instances,action is EVT_UPDATE, sid = %s, instances length = %d", response.Instance.ServiceId, len(microServiceInstances)))
}
