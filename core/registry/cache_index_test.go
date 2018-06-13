package registry

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/sets"
)

var microServiceInstances = []*MicroServiceInstance{
	{InstanceID: "10", Metadata: map[string]string{"version": "0.0.1", "project": "dev"}},
	{InstanceID: "09", Metadata: map[string]string{"version": "0.0.1", "project": "test"}},
	{InstanceID: "08", Metadata: map[string]string{"version": "0.0.1", "project": "test"}},
	{InstanceID: "07", Metadata: map[string]string{"version": "0.0.2", "project": "dev"}},
	{InstanceID: "06", Metadata: map[string]string{"version": "0.0.2", "project": "test"}},
	{InstanceID: "05", Metadata: map[string]string{"version": "0.0.3", "project": "dev"}},
	{InstanceID: "04", Metadata: map[string]string{"version": "0.0.3", "project": "dev"}},
	{InstanceID: "03", Metadata: map[string]string{"version": "0.0.3", "project": "test"}},
	{InstanceID: "02", Metadata: map[string]string{"version": "0.0.4", "project": "dev"}},
	{InstanceID: "01", Metadata: map[string]string{"version": "0.0.5", "project": "dev"}},
	{InstanceID: "11", Metadata: map[string]string{"version": "0.0.5", "project": "test"}},
}

func TestNoIndexCache(t *testing.T) {
	cache := newNoIndexCache()
	cache.Set("TestServer", microServiceInstances)
	tag1 := map[string]string{"version": "0.0.2", "project": "dev"}
	tag2 := map[string]string{"version": "latest", "project": "dev"}

	x, ok1 := cache.Get("TestServer", tag1)
	m, ok2 := x.([]*MicroServiceInstance)
	assert.Equal(t, ok1, true)
	assert.Equal(t, ok2, true)
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m[0].Metadata["version"], "0.0.2")
	assert.Equal(t, m[0].Metadata["project"], "dev")

	x, ok1 = cache.Get("TestServer", tag2)
	m, ok2 = x.([]*MicroServiceInstance)
	assert.Equal(t, ok1, true)
	assert.Equal(t, ok2, true)
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m[0].Metadata["version"], "0.0.5")
	assert.Equal(t, m[0].Metadata["project"], "dev")

	items := cache.Items()
	assert.Equal(t, len(items), 1)

	cache.Delete("TestServer")
	items = cache.Items()
	assert.Equal(t, len(items), 0)
}

func TestIndexCache(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	cache := newIndexCache()
	cache.SetIndexTags(sets.NewString("version", "project"))
	cache.Set("TestServer", microServiceInstances)
	tag1 := map[string]string{"version": "0.0.2", "project": "dev"}
	tag2 := map[string]string{"version": "latest", "project": "dev"}

	x, ok1 := cache.Get("TestServer", tag1)
	m, ok2 := x.([]*MicroServiceInstance)
	assert.Equal(t, ok1, true)
	assert.Equal(t, ok2, true)
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m[0].Metadata["version"], "0.0.2")
	assert.Equal(t, m[0].Metadata["project"], "dev")

	x, ok1 = cache.Get("TestServer", tag2)
	m, ok2 = x.([]*MicroServiceInstance)
	assert.Equal(t, ok1, true)
	assert.Equal(t, ok2, true)
	assert.Equal(t, len(m), 1)
	assert.Equal(t, m[0].Metadata["version"], "0.0.5")
	assert.Equal(t, m[0].Metadata["project"], "dev")

	items := cache.Items()
	assert.Equal(t, len(items), 1)

	cache.Delete("TestServer")
	items = cache.Items()
	assert.Equal(t, len(items), 0)
}
