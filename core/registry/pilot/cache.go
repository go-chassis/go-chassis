package pilot

import (
	"errors"
	"fmt"
	"time"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
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
	old := registry.MicroserviceInstanceIndex.Items()
	labels := registry.MicroserviceInstanceIndex.GetIndexTags()

	for serviceKey, store := range old {
		for key := range store.Items() {
			tags := pilotTags(labels, key)
			hs, err := c.registryClient.GetHostsByKey(serviceKey, tags)
			if err != nil {
				continue
			}
			filterRestore(hs.Hosts, serviceKey, tags)
		}
	}
	return nil
}

// filterRestore filter and restore instances to cache
func filterRestore(hs []*Host, serviceKey string, tags registry.Tags) {
	if len(hs) == 0 {
		registry.MicroserviceInstanceIndex.Delete(serviceKey)
		return
	}

	store := make([]*registry.MicroServiceInstance, 0, len(hs))
	for _, host := range hs {
		msi := ToMicroServiceInstance(host, tags)
		store = append(store, msi)
	}
	registry.MicroserviceInstanceIndex.Set(serviceKey, store)
}
