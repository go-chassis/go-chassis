package servicecenter

import (
	"net/url"
	"strings"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-sc-client"

	"github.com/go-chassis/go-chassis/third_party/forked/k8s.io/apimachinery/pkg/util/sets"
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
			lager.Logger.Errorf("Watch failed. Self Micro service Id:%s. %s", runtime.ServiceID, err)
		}
		lager.Logger.Debugf("Watching Intances change events.")
	}
	var ticker *time.Ticker
	refreshInterval := config.GetServiceDiscoveryRefreshInterval()
	if refreshInterval == "" {
		ticker = time.NewTicker(DefaultRefreshInterval)
	} else {
		timeValue, err := time.ParseDuration(refreshInterval)
		if err != nil {
			lager.Logger.Errorf("refeshInterval is invalid. So use Default value, err %s", err)
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
	if archaius.GetBool("cse.service.registry.autodiscovery", false) {
		err := c.registryClient.SyncEndpoints()
		if err != nil {
			lager.Logger.Errorf("SyncSCEndpoints failed: %s", err)
		}
	}
	err := c.pullMicroserviceInstance()
	if err != nil {
		lager.Logger.Errorf("AutoUpdateMicroserviceInstance failed: %s", err)
		//connection with sc may lost, reset the revision
		c.registryClient.ResetRevision()
	}

	if archaius.GetBool("cse.service.registry.autoSchemaIndex", false) {
		err = c.MakeSchemaIndex()
		if err != nil {
			lager.Logger.Errorf("MakeSchemaIndex failed: %s", err)
		}
	}

	if archaius.GetBool("cse.service.registry.autoIPIndex", false) {
		err = c.MakeIPIndex()
		if err != nil {
			lager.Logger.Errorf("Auto Update IP index failed: %s", err)
		}
	}

}

// MakeIPIndex make ip index
// if store instance metadata into tags
// it will be used in route management
func (c *CacheManager) MakeIPIndex() error {
	lager.Logger.Debug("Make IP index")
	services, err := c.registryClient.GetAllResources("instances")
	if err != nil {
		lager.Logger.Errorf("Get instances failed: %s", err)
		return err
	}
	for _, service := range services {
		for _, inst := range service.Instances {
			for _, uri := range inst.Endpoints {
				u, err := url.Parse(uri)
				if err != nil {
					lager.Logger.Errorf("Wrong URI %s: %s", uri, err)
					continue
				}
				si := &registry.SourceInfo{}
				si.Tags = inst.Properties
				if si.Tags == nil {
					si.Tags = make(map[string]string)
				}
				si.Name = service.MicroService.ServiceName
				si.Tags[common.BuildinTagApp] = service.MicroService.AppID
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

	lager.Logger.Debug("Make Schema index")
	microServiceList, err := c.registryClient.GetAllMicroServices()
	if err != nil {
		lager.Logger.Errorf("Get instances failed: %s", err)
		return err
	}

	for _, ms := range microServiceList {
		serviceID, err := c.registryClient.GetMicroServiceID(ms.AppID, ms.ServiceName, ms.Version, ms.Environment)
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
					var allMicroServices []*client.MicroService
					allMicroServices = append(allMicroServices, ms)
					registry.SchemaInterfaceIndexedCache.Set(interfaceName, allMicroServices, 0)
					lager.Logger.Debugf("New Interface added in the Index Cache : %s", interfaceName)
				} else {
					val, _ := value.([]*client.MicroService)
					if !checkIfMicroServiceExistInList(val, ms.ServiceID) {
						val = append(val, ms)
						registry.SchemaInterfaceIndexedCache.Set(interfaceName, val, 0)
						lager.Logger.Debugf("New Interface added in the Index Cache : %s", interfaceName)
					}
				}

				svcValue, ok := registry.SchemaServiceIndexedCache.Get(serviceID)
				if !ok {
					var allMicroServices []*client.MicroService
					allMicroServices = append(allMicroServices, ms)
					registry.SchemaServiceIndexedCache.Set(serviceID, allMicroServices, 0)
					lager.Logger.Debugf("New Service added in the Index Cache : %s", serviceID)
				} else {
					val, _ := svcValue.([]*client.MicroService)
					if !checkIfMicroServiceExistInList(val, ms.ServiceID) {
						val = append(val, ms)
						registry.SchemaServiceIndexedCache.Set(serviceID, val, 0)
						lager.Logger.Debugf("New Service added in the Index Cache : %s", serviceID)
					}
				}
			}
		}
	}
	return nil
}

// This functions checks if the microservices exist in the list passed in argument
func checkIfMicroServiceExistInList(microserviceList []*client.MicroService, serviceID string) bool {
	msIsPresentInList := false
	for _, interfaceMicroserviceList := range microserviceList {
		if interfaceMicroserviceList.ServiceID == serviceID {
			msIsPresentInList = true
			break
		}
	}
	return msIsPresentInList
}

// pullMicroserviceInstance pull micro-service instance
func (c *CacheManager) pullMicroserviceInstance() error {
	//Get Providers
	microServicesCache := registry.GetProvidersFromCache()
	services := make([]*client.MicroService, 0)
	for _, service := range microServicesCache {
		services = append(services, &client.MicroService{
			ServiceName: service.ServiceName,
			Version:     service.Version,
			AppID:       service.AppID,
		})
	}
	serviceNameSet, serviceNameAppIDKeySet := c.getServiceSet(services)
	c.compareAndDeleteOutdatedProviders(serviceNameSet)

	for key := range serviceNameAppIDKeySet {
		service := strings.Split(key, ":")
		if len(service) != 2 {
			lager.Logger.Errorf("Invalid serviceStore %s for providers %s", key, runtime.ServiceID)
			continue
		}

		providerInstances, err := c.registryClient.FindMicroServiceInstances(runtime.ServiceID, service[1],
			service[0], findVersionRule(service[0]))
		if err != nil {
			if err == client.ErrNotModified {
				lager.Logger.Debug(err.Error())
				continue
			}
			if err == client.ErrMicroServiceNotExists {
				registry.ProvidersMicroServiceCache.Delete(strings.Join([]string{service[0], service[1]}, "|"))
			}
			lager.Logger.Error("Refresh local instance cache failed: " + err.Error())
			continue
		}

		filterReIndex(providerInstances, service[0], service[1])
	}
	return nil
}

func (c *CacheManager) compareAndDeleteOutdatedProviders(newProviders sets.String) {
	oldProviders := registry.MicroserviceInstanceIndex.Items()
	for old := range oldProviders {
		if !newProviders.Has(old) { //provider is outdated, delete it
			registry.MicroserviceInstanceIndex.Delete(old)
		}
	}
}

// getServiceSet returns service sets
func (c *CacheManager) getServiceSet(exist []*client.MicroService) (sets.String, sets.String) {
	//get Provider's instances
	serviceNameSet := sets.NewString()         // key is serviceName
	serviceNameAppIDKeySet := sets.NewString() // key is "serviceName:appId"
	if exist == nil || len(exist) == 0 {
		return serviceNameSet, serviceNameAppIDKeySet
	}

	for _, microservice := range exist {
		if microservice == nil {
			continue
		}
		serviceNameSet.Insert(microservice.ServiceName)
		key := strings.Join([]string{microservice.ServiceName, microservice.AppID}, ":")
		serviceNameAppIDKeySet.Insert(key)
	}
	return serviceNameSet, serviceNameAppIDKeySet
}

func filterReIndex(providerInstances []*client.MicroServiceInstance, serviceName string, appID string) {
	ups := make([]*registry.MicroServiceInstance, 0, len(providerInstances))
	downs := make(map[string]struct{})
	for _, ins := range providerInstances {
		switch {
		case ins.Version == "":
			lager.Logger.Warn("do not support old service center, plz upgrade")
			continue
		case ins.Status != common.DefaultStatus:
			downs[ins.InstanceID] = struct{}{}
			lager.Logger.Debugf("do not cache the instance in '%s' status, instanceId = %s/%s",
				ins.Status, ins.ServiceID, ins.InstanceID)
			continue
		default:
			ups = append(ups, ToMicroServiceInstance(ins).WithAppID(appID))
		}
	}
	registry.RefreshCache(serviceName, ups, downs)
}

// findVersionRule returns version rules for microservice
func findVersionRule(microservice string) string {
	if ref, ok := config.GlobalDefinition.Cse.References[microservice]; ok {
		return ref.Version
	}
	return common.AllVersion
}

// watch watching micro-service instance status
func watch(response *client.MicroServiceInstanceChangedEvent) {
	if response.Instance.Status != client.MSInstanceUP {
		response.Action = common.Delete
	}
	switch response.Action {
	case client.EventCreate:
		createAction(response)
		break
	case client.EventDelete:
		deleteAction(response)
		break
	case client.EventUpdate:
		updateAction(response)
		break
	case client.EventError:
		lager.Logger.Warnf("MicroServiceInstanceChangedEvent action is error, MicroServiceInstanceChangedEvent = %s", response)
		break
	default:
		lager.Logger.Warnf("Do not support this Action = %s", response.Action)
		return
	}
}

// createAction added micro-service instance to the cache
func createAction(response *client.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	value, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		lager.Logger.Errorf("ServiceID does not exist in MicroserviceInstanceCache,action is EVT_CREATE.key = %s", key)
		return
	}
	microServiceInstances, ok := value.([]*registry.MicroServiceInstance)
	if !ok {
		lager.Logger.Errorf("Type asserts failed.action is EVT_CREATE,sid = %s", response.Instance.ServiceID)
		return
	}
	if response.Instance.Status != client.MSInstanceUP {
		lager.Logger.Warnf("createAction failed,MicroServiceInstance status is not MSI_UP,MicroServiceInstanceChangedEvent = %s", response)
		return
	}
	msi := ToMicroServiceInstance(response.Instance).WithAppID(response.Key.AppID)
	microServiceInstances = append(microServiceInstances, msi)
	registry.MicroserviceInstanceIndex.Set(key, microServiceInstances)
	lager.Logger.Debugf("Cached Instances,action is EVT_CREATE, sid = %s, instances length = %d", response.Instance.ServiceID, len(microServiceInstances))
}

// deleteAction delete micro-service instance
func deleteAction(response *client.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	lager.Logger.Debugf("Received event EVT_DELETE, sid = %s, endpoints = %s", response.Instance.ServiceID, response.Instance.Endpoints)
	if err := registry.HealthCheck(key, response.Key.Version, response.Key.AppID, ToMicroServiceInstance(response.Instance)); err == nil {
		return
	}
	value, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		lager.Logger.Errorf("ServiceID does not exist in MicroserviceInstanceCache, action is EVT_DELETE, key = %s", key)
		return
	}
	microServiceInstances, ok := value.([]*registry.MicroServiceInstance)
	if !ok {
		lager.Logger.Errorf("Type asserts failed.action is EVT_DELETE, sid = %s", response.Instance.ServiceID)
		return
	}
	var newInstances = make([]*registry.MicroServiceInstance, 0)
	for _, v := range microServiceInstances {
		if v.InstanceID != response.Instance.InstanceID {
			newInstances = append(newInstances, v)
		}
	}

	registry.MicroserviceInstanceIndex.Set(key, newInstances)
	lager.Logger.Debugf("Cached [%d] Instances of service [%s]", len(newInstances), key)
}

// updateAction update micro-service instance event
func updateAction(response *client.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName
	value, ok := registry.MicroserviceInstanceIndex.Get(key, nil)
	if !ok {
		lager.Logger.Errorf("ServiceID does not exist in MicroserviceInstanceCache, action is EVT_UPDATE, sid = %s", key)
		return
	}
	microServiceInstances, ok := value.([]*registry.MicroServiceInstance)
	if !ok {
		lager.Logger.Errorf("Type asserts failed.action is EVT_UPDATE, sid = %s", response.Instance.ServiceID)
		return
	}
	if response.Instance.Status != client.MSInstanceUP {
		lager.Logger.Warnf("updateAction failed, MicroServiceInstance status is not MSI_UP, MicroServiceInstanceChangedEvent = %s", response)
		return
	}
	msi := ToMicroServiceInstance(response.Instance).WithAppID(response.Key.AppID)
	var iidExist = InstanceIDIsNotExist
	var arrayNum int
	for k, v := range microServiceInstances {
		if v.InstanceID == response.Instance.InstanceID {
			iidExist = InstanceIDIsExist
			arrayNum = k
		}
	}
	switch iidExist {
	case InstanceIDIsExist:
		microServiceInstances[arrayNum] = msi
		break
	case InstanceIDIsNotExist:
		microServiceInstances = append(microServiceInstances, msi)
		break
	default:
		lager.Logger.Warnf("updateAction error, iid:%s", response.Instance.InstanceID)
	}
	registry.MicroserviceInstanceIndex.Set(key, microServiceInstances)
	lager.Logger.Debugf("Cached Instances,action is EVT_UPDATE, sid = %s, instances length = %d", response.Instance.ServiceID, len(microServiceInstances))
}
