package registry

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	cache "github.com/patrickmn/go-cache"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	//DefaultExpireTime default expiry time is kept as 0
	DefaultExpireTime = 0
	//Version defines tag name of version
	Version = "version"
	//App defines tag name of app
	App = "app"
)

//MicroserviceInstanceIndex key: ServiceName, value: []instance
var MicroserviceInstanceIndex CacheIndex

//SelfInstancesCache key: serviceID, value: []instanceID
var SelfInstancesCache *cache.Cache

//IPIndexedCache key: instance ip, value: microservice
var IPIndexedCache *cache.Cache

//SchemaInterfaceIndexedCache key: schema interface name value: []*microservice
var SchemaInterfaceIndexedCache *cache.Cache

//SchemaServiceIndexedCache key: schema service name value: []*microservice
var SchemaServiceIndexedCache *cache.Cache

func initCache() *cache.Cache { return cache.New(DefaultExpireTime, 0) }

func enableRegistryCache() {
	MicroserviceInstanceIndex = newCacheIndex()
	SelfInstancesCache = initCache()
	IPIndexedCache = initCache()
	SchemaServiceIndexedCache = initCache()
	SchemaInterfaceIndexedCache = initCache()
}

// CacheIndex defines interface for cache and index used by registry
type CacheIndex interface {
	SetIndexTags(tags sets.String)
	Get(k string, tags map[string]string) (interface{}, bool)
	Set(k string, x interface{})

	Items() map[string]cache.Item
	Delete(k string)
}

// newCacheIndex returns default index implemention
func newCacheIndex() CacheIndex {
	if archaius.GetBool("cse.service.registry.cacheIndex", false) {
		return newIndexCache()
	}
	return newNoIndexCache()
}

// Tags defines query conditions
type Tags map[string]string

// NewDefaultTag returns Tags with version and appID
func NewDefaultTag(version, appID string) Tags {
	return Tags{Version: version, App: appID}
}
