package registry_test

import (
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrapInstance(t *testing.T) {
	wi := registry.WrapInstance{
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
	registry.RefreshCache("test", nil, nil)
	// case: refresh nil simpleCache
	registry.RefreshCache("test", nil, nil)

	// case: new instances
	registry.RefreshCache("test", []*registry.MicroServiceInstance{}, nil)
	registry.RefreshCache("test", []*registry.MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus},
		{InstanceID: "2", Status: common.DefaultStatus}}, nil) // 2

	is, ok := registry.MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(is))

	// case: unregister one
	registry.RefreshCache("test", []*registry.MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus},
		{InstanceID: "3", Status: common.DefaultStatus}}, nil)

	is, ok = registry.MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(is))

	// case: down one with non-up status
	registry.RefreshCache("test", []*registry.MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus}},
		map[string]struct{}{"3": {}})

	is, ok = registry.MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(is))

	// case: coming in with non-up status
	registry.RefreshCache("test", []*registry.MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus}},
		map[string]struct{}{"3": {}})

	is, ok = registry.MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(is))
}
