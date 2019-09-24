package config_test

import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRouterType(t *testing.T) {
	config.OldRouterDefinition = &config.RouterConfig{Router: config.Router{
		Infra: "",
	}}
	assert.Equal(t, "cse", config.GetRouterType())

	config.OldRouterDefinition = &config.RouterConfig{Router: config.Router{
		Infra: "test",
	}}
	assert.Equal(t, "test", config.GetRouterType())

	config.OldRouterDefinition = &config.RouterConfig{Router: config.Router{
		Infra:   "test",
		Address: "123",
	}}
	assert.Equal(t, "123", config.GetRouterEndpoints())
}
