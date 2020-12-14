package health_test

import (
	"testing"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/health"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.registry.address", "http://127.0.0.1:30100")
	archaius.Set("servicecomb.service.name", "Client")
	runtime.HostName = "localhost"
	config.MicroserviceDefinition = &model.ServiceSpec{}
	archaius.UnmarshalConfig(config.MicroserviceDefinition)
	config.ReadGlobalConfigFromArchaius()
}
func TestWrapInstance(t *testing.T) {
	wi := health.WrapInstance{
		AppID:       "1",
		ServiceName: "2",
		Version:     "3",
		Instance:    &registry.MicroServiceInstance{InstanceID: "4"},
	}
	assert.Equal(t, "2:3:1:4", wi.String())
	assert.Equal(t, "2:3:1", wi.ServiceKey())
}

func TestRefreshCache(t *testing.T) {
	registry.EnableRegistryCache()

	// case: new nil simpleCache
	health.RefreshCache("test", nil, nil)
	// case: refresh nil simpleCache
	health.RefreshCache("test", nil, nil)

	// case: new instances
	health.RefreshCache("test", []*registry.MicroServiceInstance{}, nil)
	health.RefreshCache("test", []*registry.MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus},
		{InstanceID: "2", Status: common.DefaultStatus}}, nil) // 2

	is, ok := registry.MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(is))

	// case: unregister one
	health.RefreshCache("test", []*registry.MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus},
		{InstanceID: "3", Status: common.DefaultStatus}}, nil)

	is, ok = registry.MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(is))

	// case: down one with non-up status
	health.RefreshCache("test", []*registry.MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus}},
		map[string]struct{}{"3": {}})

	is, ok = registry.MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(is))

	// case: coming in with non-up status
	health.RefreshCache("test", []*registry.MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus}},
		map[string]struct{}{"3": {}})

	is, ok = registry.MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(is))
}
