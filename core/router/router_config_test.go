package router_test

import (
	"errors"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/router"
	_ "github.com/go-chassis/go-chassis/v2/initiator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	config.Init()
	t.Run("install and build a wrong router,return err", func(t *testing.T) {
		router.InstallRouterPlugin("wrong", func() (router.Router, error) {
			return nil, errors.New("1")
		})
		err := router.BuildRouter("wrong")
		assert.Error(t, err)
	})
	t.Run("validate rule, exact 100, should success", func(t *testing.T) {
		ok := router.ValidateRule(map[string][]*config.RouteRule{
			"service1": {
				{Routes: []*config.RouteTag{{Weight: 50}, {Weight: 50}}},
			},
		})
		assert.True(t, ok)
	})
	t.Run("validate rule, greater than 100, should fail", func(t *testing.T) {
		ok := router.ValidateRule(map[string][]*config.RouteRule{
			"service2": {
				{Routes: []*config.RouteTag{{Weight: 900}, {Weight: 50}}},
			},
		})
		assert.False(t, ok)
	})

}

func BenchmarkGenWeightPoolKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		router.GenWeightPoolKey("test", 1)
	}
}
