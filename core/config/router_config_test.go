package config_test

import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRouterType(t *testing.T) {
	config.RouterDefinition = &model.RouterConfig{Router: model.Router{
		Infra: "",
	}}
	assert.Equal(t, "cse", config.GetRouterType())

	config.RouterDefinition = &model.RouterConfig{Router: model.Router{
		Infra: "test",
	}}
	assert.Equal(t, "test", config.GetRouterType())

	config.RouterDefinition = &model.RouterConfig{Router: model.Router{
		Infra:   "test",
		Address: "123",
	}}
	assert.Equal(t, "123", config.GetRouterEndpoints())
}
