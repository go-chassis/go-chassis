package servicecenter_test

import (
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCriteria(t *testing.T) {
	registry.ProvidersMicroServiceCache = cache.New(0, 0)
	registry.AddProviderToCache("service1", "1")
	registry.AddProviderToCache("service1", "2")
	registry.AddProviderToCache("service2", "2")
	c := servicecenter.GetCriteria()
	assert.Equal(t, 3, len(c))
	c = servicecenter.GetCriteriaByService("service1")
	assert.Equal(t, 2, len(c))
	c = servicecenter.GetCriteriaByService("service2")
	assert.Equal(t, 1, len(c))
}
