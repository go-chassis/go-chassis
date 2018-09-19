package registry

import (
	"testing"

	"strings"

	"fmt"
	"sort"

	"github.com/stretchr/testify/assert"
)

func TestSetIPIndex(t *testing.T) {
	enableRegistryCache()
	SetIPIndex("10.1.0.1", &SourceInfo{
		Name: "ServerA",
	})
	si := GetIPIndex("10.1.0.1")
	assert.Equal(t, "ServerA", si.Name)

	si = GetIPIndex("10.1.1.1")
	assert.Nil(t, si)
}
func TestGetProvidersFromCache(t *testing.T) {
	enableRegistryCache()

	AddProviderToCache("SERVER0", "default")
	AddProviderToCache("SERVER1", "default")
	AddProviderToCache("SERVER2", "default")

	services := GetProvidersFromCache()
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
	enableRegistryCache()
	testMap := map[string]string{"SERVER1": "default", "SERVER2": "default", "SERVER3": "default"}
	for key, value := range testMap {
		AddProviderToCache(key, value)
	}

	for key, value := range testMap {
		v, ok := ProvidersMicroServiceCache.Get(strings.Join([]string{key, value}, "|"))
		assert.Equal(t, ok, true)

		microService := v.(MicroService)
		assert.Equal(t, key, microService.ServiceName)
		assert.Equal(t, value, microService.AppID)
	}

}
