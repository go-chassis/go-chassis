package servicecenter

import (
	"net/url"
	"strings"
	"time"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"

	"github.com/ServiceComb/go-sc-client"
	"github.com/ServiceComb/go-sc-client/model"
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
		err := c.registryClient.WatchMicroService(config.SelfServiceID, watch)
		if err != nil {
			lager.Logger.Errorf(err, "Watch failed. Self Micro service Id:%s.", config.SelfServiceID)
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
			lager.Logger.Errorf(err, "refeshInterval is invalid. So use Default value")
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
			lager.Logger.Errorf(err, "SyncSCEndpoints failed.")
		}
	}
	err := c.pullMicroserviceInstance()
	if err != nil {
		lager.Logger.Errorf(err, "AutoUpdateMicroserviceInstance failed.")
	}

	if archaius.GetBool("cse.service.registry.autoSchemaIndex", false) {
		err = c.MakeSchemaIndex()
		if err != nil {
			lager.Logger.Errorf(err, "MakeSchemaIndex failed.")
		}
	}

	if archaius.GetBool("cse.service.registry.autoIPIndex", false) {
		err = c.MakeIPIndex()
		if err != nil {
			lager.Logger.Errorf(err, "Auto Update IP index failed.")
		}
	}

}

// MakeIPIndex make ip index
func (c *CacheManager) MakeIPIndex() error {
	lager.Logger.Debug("Make IP index")
	services, err := c.registryClient.GetAllResources("instances")
	if err != nil {
		lager.Logger.Errorf(err, "Get instances failed")
		return err
	}
	for _, service := range services {
		for _, inst := range service.Instances {
			for _, uri := range inst.Endpoints {
				u, err := url.Parse(uri)
				if err != nil {
					lager.Logger.Error("Wrong URI", err)
					continue
				}
				u.Host = strings.Split(u.Host, ":")[0]
				registry.IPIndexedCache.Set(u.Host, service.MicroService, 0)
				//no need to analyze each endpoint
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
		lager.Logger.Errorf(err, "Get instances failed")
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

			interfaceName := schemaContent.Info["x-java-interfcae"]
			value, ok := registry.SchemaInterfaceIndexedCache.Get(interfaceName)
			if !ok {
				var allMicroServices []*model.MicroService
				allMicroServices = append(allMicroServices, ms)
				registry.SchemaInterfaceIndexedCache.Set(interfaceName, allMicroServices, 0)
			} else {
				val, _ := value.([]*model.MicroService)
				val = append(val, ms)
				registry.SchemaInterfaceIndexedCache.Set(interfaceName, val, 0)
			}
			svcValue, ok := registry.SchemaServiceIndexedCache.Get(serviceID)
			if !ok {
				var allMicroServices []*model.MicroService
				allMicroServices = append(allMicroServices, ms)
				registry.SchemaServiceIndexedCache.Set(serviceID, allMicroServices, 0)
			} else {
				val, _ := svcValue.([]*model.MicroService)
				val = append(val, ms)
				registry.SchemaServiceIndexedCache.Set(serviceID, val, 0)
			}

		}
	}

	return nil
}

// pullMicroserviceInstance pull micro-service instance
func (c *CacheManager) pullMicroserviceInstance() error {
	//Get Providers
	rsp, err := c.registryClient.GetProviders(config.SelfServiceID)
	if err != nil {
		lager.Logger.Errorf(err, "get Providers failed, sid = %s", config.SelfServiceID)
		return err
	}

	serviceStore := c.getServiceStore(rsp.Services)
	for key := range serviceStore {
		service := strings.Split(key, ":")
		if len(service) != 2 {
			lager.Logger.Errorf(err, "Invalid servcieStore %s for providers %s", key, config.SelfServiceID)
			continue
		}

		providerInstances, err := c.registryClient.FindMicroServiceInstances(config.SelfServiceID, service[1],
			service[0], findVersionRule(service[0]))
		if err != nil {
			lager.Logger.Errorf(err, "GetMicroServiceInstances failed")
			continue
		}

		filterRestore(providerInstances, service[0], service[1])
	}
	return nil
}

// getServiceStore returns service sets
func (c *CacheManager) getServiceStore(exist []*model.MicroService) sets.String {
	//get Provider's instances
	serviceStore := sets.NewString()
	for _, microservice := range exist {
		key := strings.Join([]string{microservice.ServiceName, microservice.AppID}, ":")
		if !serviceStore.Has(key) {
			serviceStore.Insert(key)
		}
	}

	if archaius.GetBool("cse.service.registry.autoClearCache", false) {
		c.autoClearCache(serviceStore)
	}
	return serviceStore
}

// autoClearCache clear cache for non exist service
func (c *CacheManager) autoClearCache(exist sets.String) {
	old := registry.MicroserviceInstanceCache.Items()
	delsets := sets.NewString()

	for key := range old {
		ins := strings.Split(key, ":")
		if len(ins) != 3 {
			continue
		}

		skey := strings.Join([]string{ins[0], ins[2]}, ":")
		if !exist.Has(skey) {
			delsets.Insert(key)
		}
	}

	for insKey := range delsets {
		registry.MicroserviceInstanceCache.Delete(insKey)
	}
}

// filterRestore filter and restore instances to cache
func filterRestore(providerInstances []*model.MicroServiceInstance, serviceName, appID string) {
	var (
		latestKey = strings.Join([]string{serviceName, common.LatestVersion, appID}, ":")
		store     = make(map[string][]*registry.MicroServiceInstance, 0)
		latest    string
		index     string
	)

	for _, ins := range providerInstances {
		if ins.Status != model.MSInstanceUP {
			continue
		}

		// keep this for compatibility with old service-center
		var key string
		if ins.Version != "" {
			key = strings.Join([]string{serviceName, ins.Version, appID}, ":")
			if ins.Version > latest {
				latest, index = ins.Version, key
			}
		} else {
			key = latestKey
			index = common.LatestVersion
		}

		msi := ToMicroServiceInstance(ins)
		if _, ok := store[key]; !ok {
			store[key] = make([]*registry.MicroServiceInstance, 0)
		}
		store[key] = append(store[key], msi)
	}

	for key, ins := range store {
		registry.MicroserviceInstanceCache.Set(key, ins, 0)
		lager.Logger.Debugf("Cached [%d] Instances of service [%s]", len(ins), key)
	}

	if index != common.LatestVersion {
		registry.MicroserviceInstanceCache.Set(latestKey, store[index], 0)
		lager.Logger.Debugf("Cached [%d] Instances of service [%s]", len(store[index]), latestKey)
	}
}

// findVersionRule returns version rules for microservice
func findVersionRule(microservice string) string {
	if ref, ok := config.GlobalDefinition.Cse.References[microservice]; ok {
		return ref.Version
	}
	return common.AllVersion
}

// watch watching micro-service instance status
func watch(response *model.MicroServiceInstanceChangedEvent) {
	if response.Instance.Status != model.MSInstanceUP {
		response.Action = common.Delete
	}
	switch response.Action {
	case model.EventCreate:
		createAction(response)
		break
	case model.EventDelete:
		deleteAction(response)
		break
	case model.EventUpdate:
		updateAction(response)
		break
	case model.EventError:
		lager.Logger.Warnf("MicroServiceInstanceChangedEvent action is error, MicroServiceInstanceChangedEvent = %s", response)
		break
	default:
		lager.Logger.Warnf("Do not support this Action = %s", response.Action)
		return
	}
}

// createAction added micro-service instance to the cache
func createAction(response *model.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName + ":" + response.Key.Version + ":" + response.Key.AppID
	value, ok := registry.MicroserviceInstanceCache.Get(key)
	if !ok {
		lager.Logger.Errorf(nil, "ServiceID does not exist in MicroserviceInstanceCache,action is EVT_CREATE.key = %s", key)
		return
	}
	microServiceInstances, ok := value.([]*registry.MicroServiceInstance)
	if !ok {
		lager.Logger.Errorf(nil, "Type asserts failed.action is EVT_CREATE,sid = %s", response.Instance.ServiceID)
		return
	}
	if response.Instance.Status != model.MSInstanceUP {
		lager.Logger.Warnf("createAction failed,MicroServiceInstance status is not MSI_UP,MicroServiceInstanceChangedEvent = %s", response)
		return
	}
	msi := ToMicroServiceInstance(response.Instance)
	microServiceInstances = append(microServiceInstances, msi)
	registry.MicroserviceInstanceCache.Set(key, microServiceInstances, 0)
	lager.Logger.Debugf("Cached Instances,action is EVT_CREATE, sid = %s, instances length = %d", response.Instance.ServiceID, len(microServiceInstances))
}

// deleteAction delete micro-service instance
func deleteAction(response *model.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName + ":" + response.Key.Version + ":" + response.Key.AppID
	value, ok := registry.MicroserviceInstanceCache.Get(key)
	if !ok {
		lager.Logger.Errorf(nil, "ServiceID does not exist in MicroserviceInstanceCache,action is EVT_DELETE, key = %s", key)
		return
	}
	microServiceInstances, ok := value.([]*registry.MicroServiceInstance)
	if !ok {
		lager.Logger.Errorf(nil, "Type asserts failed.action is EVT_DELETE,sid = %s", response.Instance.ServiceID)
		return
	}
	var newInstances []*registry.MicroServiceInstance = make([]*registry.MicroServiceInstance, 0)
	for _, v := range microServiceInstances {
		if v.InstanceID != response.Instance.InstanceID {
			newInstances = append(newInstances, v)
		}
	}
	registry.MicroserviceInstanceCache.Set(key, newInstances, 0)
	lager.Logger.Debugf("Cached Instances,action is EVT_DELETE,sid = %s, instances length = %d", response.Instance.ServiceID, len(newInstances))
}

// updateAction update micro-service instance event
func updateAction(response *model.MicroServiceInstanceChangedEvent) {
	key := response.Key.ServiceName + ":" + response.Key.Version + ":" + response.Key.AppID
	value, ok := registry.MicroserviceInstanceCache.Get(key)
	if !ok {
		lager.Logger.Errorf(nil, "ServiceID does not exist in MicroserviceInstanceCache,action is EVT_UPDATE,sid = %s", key)
		return
	}
	microServiceInstances, ok := value.([]*registry.MicroServiceInstance)
	if !ok {
		lager.Logger.Errorf(nil, "Type asserts failed.action is EVT_UPDATE,sid = %s", response.Instance.ServiceID)
		return
	}
	if response.Instance.Status != model.MSInstanceUP {
		lager.Logger.Warnf("updateAction failed,MicroServiceInstance status is not MSI_UP,MicroServiceInstanceChangedEvent = %s", response)
		return
	}
	msi := ToMicroServiceInstance(response.Instance)
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
	registry.MicroserviceInstanceCache.Set(key, microServiceInstances, 0)
	lager.Logger.Debugf("Cached Instances,action is EVT_UPDATE,sid = %s,instances length = %d", response.Instance.ServiceID, len(microServiceInstances))
}
