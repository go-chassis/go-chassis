package tracing_test

import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/tracing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTracerManager(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.GlobalDefinition = &model.GlobalCfg{}

	err := tracing.Init()
	assert.NoError(t, err)
	err = tracing.Init()
	assert.NoError(t, err)
}
