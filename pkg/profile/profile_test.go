package profile

import (
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/router"
	_ "github.com/go-chassis/go-chassis/v2/core/router/servicecomb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProfile(t *testing.T) {
	err := router.BuildRouter("cse")
	assert.NoError(t, err)
	rr := map[string][]*config.RouteRule{"test": {{Precedence: 10}}}
	router.DefaultRouter.SetRouteRule(rr)

	registry.MicroserviceInstanceIndex = registry.NewIndexCache()
	registry.MicroserviceInstanceIndex.Set("test", []*registry.MicroServiceInstance{{InstanceID: "id"}})

	p := newProfile()

	assert.Equal(t, 10, p.RouteRule["test"][0].Precedence)
	assert.Equal(t, "id", p.Discovery["test"][0].InstanceID)
}
