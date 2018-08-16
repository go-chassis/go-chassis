package registry

import (
	"github.com/patrickmn/go-cache"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	//DefaultExpireTime default expiry time is kept as 0
	DefaultExpireTime = 0
)

//MicroserviceInstanceIndex key: ServiceName, value: []instance
var MicroserviceInstanceIndex CacheIndex

//SelfInstancesCache key: serviceID, value: []instanceID
var SelfInstancesCache *cache.Cache

//ipIndexedCache is for caching map of instance IP and service information
//key: instance ip, value: SourceInfo
var ipIndexedCache *cache.Cache

//SchemaInterfaceIndexedCache key: schema interface name value: []*microservice
var SchemaInterfaceIndexedCache *cache.Cache

//SchemaServiceIndexedCache key: schema service name value: []*microservice
var SchemaServiceIndexedCache *cache.Cache

func initCache() *cache.Cache { return cache.New(DefaultExpireTime, 0) }

func enableRegistryCache() {
	MicroserviceInstanceIndex = newCacheIndex()
	SelfInstancesCache = initCache()
	ipIndexedCache = initCache()
	SchemaServiceIndexedCache = initCache()
	SchemaInterfaceIndexedCache = initCache()
}

// CacheIndex defines interface for cache and index used by registry
type CacheIndex interface {
	GetIndexTags() []string
	SetIndexTags(tags sets.String)
	Get(k string, tags map[string]string) (interface{}, bool)
	Set(k string, x interface{})
	Items() map[string]*cache.Cache
	Delete(k string)
}

//SetIPIndex save ip index
func SetIPIndex(ip string, si *SourceInfo) {
	ipIndexedCache.Set(ip, si, 0)
}

//GetIPIndex get ip corresponding source info
func GetIPIndex(ip string) *SourceInfo {
	cacheDatum, ok := ipIndexedCache.Get(ip)
	if !ok {
		return nil
	}
	si, ok := cacheDatum.(*SourceInfo)
	if !ok {
		return nil
	}
	return si
}

// SetNoIndexCache reset microservie instance index to no index cache
func SetNoIndexCache() { MicroserviceInstanceIndex = newNoIndexCache() }

// newCacheIndex returns index implemention according to config
func newCacheIndex() CacheIndex { return newIndexCache() }
