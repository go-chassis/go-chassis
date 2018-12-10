package registry

import (
	"testing"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/stretchr/testify/assert"
)

var microServiceInstances = []*MicroServiceInstance{
	{InstanceID: "09", Metadata: map[string]string{"version": "0.0.1", "project": "dev"}},
	{InstanceID: "08", Metadata: map[string]string{"version": "0.0.1", "project": "test"}},
	{InstanceID: "07", Metadata: map[string]string{"version": "0.0.2", "project": "dev"}},
	{InstanceID: "06", Metadata: map[string]string{"version": "0.0.2", "project": "test"}},
	{InstanceID: "05", Metadata: map[string]string{"version": "0.0.3", "project": "dev"}},
	{InstanceID: "04", Metadata: map[string]string{"version": "0.0.3", "project": "dev"}},
	{InstanceID: "03", Metadata: map[string]string{"version": "0.0.3", "project": "test"}},
	{InstanceID: "02", Metadata: map[string]string{"version": "0.0.4", "project": "dev"}},
	{InstanceID: "01", Metadata: map[string]string{"version": "0.0.5", "project": "dev"}},
	{InstanceID: "10", Metadata: map[string]string{"version": "0.1", "project": "dev"}},
}

func TestNoIndexCache(t *testing.T) {
	cache := NewIndexCache()
	cache.Set("TestServer", microServiceInstances)
	instance, _ := cache.Get("TestServer", nil)
	assert.Equal(t, len(microServiceInstances), len(instance))
}

func TestIndexCache(t *testing.T) {
	cache := NewIndexCache()
	cache.Set("TestServer", microServiceInstances)
	//tag2 := map[string]string{"version": "latest", "project": "dev"}

	x, _ := cache.Get("TestServer", map[string]string{"version": "0.0.2", "project": "dev"})
	assert.Equal(t, 1, len(x))
	assert.Equal(t, "0.0.2", x[0].Metadata[common.BuildinTagVersion])
	assert.Equal(t, "dev", x[0].Metadata["project"])

	x, _ = cache.Get("TestServer", map[string]string{"version": "0.0.1"})
	assert.Equal(t, 2, len(x))
	assert.Equal(t, "0.0.1", x[0].Metadata[common.BuildinTagVersion])
	assert.Equal(t, "0.0.1", x[1].Metadata[common.BuildinTagVersion])

	microServiceInstances = append(microServiceInstances,
		&MicroServiceInstance{Metadata: map[string]string{"version": "0.0.1", "project": "dev"}})
	cache.Set("TestServer", microServiceInstances)
	x, _ = cache.Get("TestServer", map[string]string{"version": "0.0.1"})
	for _, i := range x {
		t.Log(i.InstanceID)
	}
	assert.Equal(t, 3, len(x))

	x, _ = cache.Get("TestServer", map[string]string{"version": "latest"})
	assert.Equal(t, 1, len(x))
	assert.Equal(t, "0.1", x[0].Metadata[common.BuildinTagVersion])

	cache.Delete("TestServer")
}
func TestIndexCache_Get(t *testing.T) {
	k1 := getIndexedCacheKey("service1", map[string]string{
		"a": "b",
		"c": "d",
	})
	k2 := getIndexedCacheKey("service1", map[string]string{
		"c": "d",
		"a": "b",
	})
	t.Log(k1)
	assert.Equal(t, k2, k1)
}
func BenchmarkNoIndexGet(b *testing.B) {
	cache := NewIndexCache()
	cache.Set("TestServer", microServiceInstances)
	tag := map[string]string{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("TestServer", tag)
	}
	b.ReportAllocs()
}

func BenchmarkIndexCacheGet(b *testing.B) {
	cache := NewIndexCache()
	cache.Set("TestServer", microServiceInstances)

	tag := map[string]string{"version": "0.0.3", "project": "dev"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("TestServer", tag)
	}
	b.ReportAllocs()
}
