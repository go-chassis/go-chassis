package tracing_test

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/tracing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTracerManager(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.GlobalDefinition = &model.GlobalCfg{}

	config.GlobalDefinition.Tracing.CollectorType = tracing.TracingZipkinCollector
	config.GlobalDefinition.Tracing.CollectorTarget = "localhost:9441/v1/spans"
	err := tracing.Init()
	assert.NoError(t, err)
	err = tracing.Init()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tracing.TracerMap))
}

func TestTracerManagerError(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Tracing.CollectorType = "errortype"
	config.GlobalDefinition.Tracing.CollectorTarget = "localhost:9441/v1/spans"
	tracing.GetTracer("calltracer")
	err := tracing.Init()
	assert.Error(t, err)
}
