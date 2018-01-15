package registry

import (
	client "github.com/ServiceComb/go-sc-client"
	cache "github.com/patrickmn/go-cache"
)

const (
	//DefaultExpireTime default expiry time is kept as 0
	DefaultExpireTime = 0
)

//MicroserviceInstanceCache key: ServiceName:Version:AppID, value: []instance
var MicroserviceInstanceCache *cache.Cache

//SelfInstancesCache key: serviceID, value: []instanceID
var SelfInstancesCache *cache.Cache

//IPIndexedCache key: instance ip, value: microservice
var IPIndexedCache *cache.Cache

//SchemaInterfaceIndexedCache key: schema interface name value: []*microservice
var SchemaInterfaceIndexedCache *cache.Cache

//SchemaServiceIndexedCache key: schema service name value: []*microservice
var SchemaServiceIndexedCache *cache.Cache

//CacheManager cache manager struct
type CacheManager struct {
	registryClient *client.RegistryClient
}

func initCache() *cache.Cache {
	var value *cache.Cache
	value = cache.New(DefaultExpireTime, 0)
	return value
}

func init() {
	MicroserviceInstanceCache = initCache()
	SelfInstancesCache = initCache()
	IPIndexedCache = initCache()
	SchemaServiceIndexedCache = initCache()
	SchemaInterfaceIndexedCache = initCache()
}
