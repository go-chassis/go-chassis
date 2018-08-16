package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
