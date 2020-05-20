package registry_test

import (
	"github.com/go-chassis/go-chassis/core/registry"
	"testing"

	"strings"

	"fmt"
	"sort"

	"github.com/stretchr/testify/assert"
)

func TestSetIPIndex(t *testing.T) {
	registry.EnableRegistryCache()
	registry.SetIPIndex("10.1.0.1", &registry.SourceInfo{
		Name: "ServerA",
	})
	si := registry.GetIPIndex("10.1.0.1")
	assert.Equal(t, "ServerA", si.Name)

	si = registry.GetIPIndex("10.1.1.1")
	assert.Nil(t, si)
}
func TestGetProvidersFromCache(t *testing.T) {
	registry.EnableRegistryCache()

	registry.AddProviderToCache("SERVER0", "default")
	registry.AddProviderToCache("SERVER1", "default")
	registry.AddProviderToCache("SERVER2", "default")

	services := registry.GetProvidersFromCache()
	assert.Equal(t, len(services), 3)

	serverNames := []string{}
	for _, v := range services {
		serverNames = append(serverNames, v.ServiceName+"|"+v.AppID)
	}
	sort.Strings(serverNames)
	for i, v := range serverNames {
		serverName := fmt.Sprint("SERVER", i, "|", "default")
		assert.Equal(t, serverName, v)
	}
}
func TestAddProviderToCache(t *testing.T) {
	registry.EnableRegistryCache()
	testMap := map[string]string{"SERVER1": "default", "SERVER2": "default", "SERVER3": "default"}
	for key, value := range testMap {
		registry.AddProviderToCache(key, value)
	}

	for key, value := range testMap {
		v, ok := registry.ProvidersMicroServiceCache.Get(strings.Join([]string{key, value}, "|"))
		assert.Equal(t, ok, true)

		microService := v.(registry.MicroService)
		assert.Equal(t, key, microService.ServiceName)
		assert.Equal(t, value, microService.AppID)
	}

}
