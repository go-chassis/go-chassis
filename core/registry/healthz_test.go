package registry

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestWrapInstance(t *testing.T) {
	wi := WrapInstance{
		AppID:       "1",
		ServiceName: "2",
		Version:     "3",
		Instance:    &MicroServiceInstance{InstanceID: "4"},
	}
	assert.Equal(t, "2:3:1:4", wi.String())
	assert.Equal(t, "2:3:1", wi.ServiceKey())
}

func TestRefreshCache(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	enableRegistryCache()

	// case: new nil cache
	RefreshCache("test", nil)
	// case: refresh nil cache
	RefreshCache("test", nil)

	// case: new instances
	RefreshCache("test", []*MicroServiceInstance{})
	RefreshCache("test", []*MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus},
		{InstanceID: "2", Status: common.DefaultStatus}}) // 2

	is, ok := MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(is.([]*MicroServiceInstance)))

	// case: down one
	RefreshCache("test", []*MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus},
		{InstanceID: "3", Status: common.DefaultStatus}})

	is, ok = MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 2, len(is.([]*MicroServiceInstance)))

	// case: down one with non-up status
	RefreshCache("test", []*MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus},
		{InstanceID: "3", Status: "xxx"}})

	is, ok = MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(is.([]*MicroServiceInstance)))

	// case: coming in with non-up status
	RefreshCache("test", []*MicroServiceInstance{
		{InstanceID: "1", Status: common.DefaultStatus},
		{InstanceID: "3", Status: "xxx"}})

	is, ok = MicroserviceInstanceIndex.Get("test", nil)
	assert.Equal(t, true, ok)
	assert.Equal(t, 1, len(is.([]*MicroServiceInstance)))
}
