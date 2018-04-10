package pilot

import (
	"time"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"

	"errors"
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
)

// constant values for default expiration time, and refresh interval
const (
	DefaultExpireTime      = 0
	DefaultRefreshInterval = time.Second * 30
)

// CacheManager cache manager
type CacheManager struct {
	registryClient *EnvoyDSClient
}

// AutoSync automatically syncing with the running instances
func (c *CacheManager) AutoSync() {
	c.refreshCache()
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
		// TODO CDS
		lager.Logger.Errorf(errors.New("not supported"), "SyncPilotEndpoints failed.")
	}
	err := c.pullMicroserviceInstance()
	if err != nil {
		lager.Logger.Errorf(err, "AutoUpdateMicroserviceInstance failed.")
	}

	if archaius.GetBool("cse.service.registry.autoSchemaIndex", false) {
		lager.Logger.Errorf(errors.New("Not support operation"), "MakeSchemaIndex failed.")
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
	services, err := c.registryClient.GetAllServices()
	if err != nil {
		lager.Logger.Errorf(err, "Get instances failed")
		return err
	}
	for _, service := range services {
		for _, h := range service.Hosts {
			registry.IPIndexedCache.Set(fmt.Sprintf("%s:%d", h.Address, h.Port), service.ServiceKey, 0)
			//no need to analyze each endpoint
			break
		}
	}
	return nil
}

// pullMicroserviceInstance pull micro-service instance
func (c *CacheManager) pullMicroserviceInstance() error {
	//Get Providers
	services, err := c.registryClient.GetAllServices()
	if err != nil {
		lager.Logger.Errorf(err, "get Providers failed, sid = %s", config.SelfServiceID)
		return err
	}

	// just auto clean cache
	c.getServiceStore(services)

	for _, service := range services {
		filterRestore(service.Hosts, service.ServiceKey)
	}
	return nil
}

// getServiceStore returns service sets
func (c *CacheManager) getServiceStore(exist []*Service) sets.String {
	//get Provider's instances
	serviceStore := sets.NewString()
	for _, microservice := range exist {
		if !serviceStore.Has(microservice.ServiceKey) {
			serviceStore.Insert(microservice.ServiceKey)
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
		if !exist.Has(key) {
			delsets.Insert(key)
		}
	}

	for insKey := range delsets {
		registry.MicroserviceInstanceCache.Delete(insKey)
	}
}

// filterRestore filter and restore instances to cache
func filterRestore(hs []*Host, serviceName string) {
	var store []*registry.MicroServiceInstance
	for _, ins := range hs {
		msi := ToMicroServiceInstance(ins)
		store = append(store, msi)
	}
	registry.MicroserviceInstanceCache.Set(serviceName, store, 0)
}
