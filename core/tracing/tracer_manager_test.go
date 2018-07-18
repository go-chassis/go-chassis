package tracing_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/tracing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
)

func TestTracerManager(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)
	// use nil option to init
	err := tracing.Init(nil)
	assert.Error(t, err)

	// when no config is provided should do nothing
	err = tracing.Init(&tracing.Option{})
	assert.NoError(t, err)
	_, ok := tracing.ConsumerTracer().(opentracing.NoopTracer)
	assert.True(t, ok)

	// when collector is invalid, should return err and do nothing
	err = tracing.Init(&tracing.Option{
		CollectorType:   "invalidType",
		CollectorTarget: "target",
	})
	assert.Error(t, err)
	_, ok = tracing.ConsumerTracer().(opentracing.NoopTracer)
	assert.True(t, ok)

	// when collector is valid and protocol endpoint is nil,
	// only consumer tracer should init
	err = tracing.Init(&tracing.Option{
		CollectorType:   "zipkin",
		CollectorTarget: "http://localhost:9411/api/v1/spans",
	})
	assert.NoError(t, err)
	_, ok = tracing.ConsumerTracer().(opentracing.NoopTracer)
	assert.False(t, ok)
	_, err = tracing.ProviderTracer(common.ProtocolRest)
	assert.Error(t, err)

	// when provide collector and protocol endpoint,
	// provider tracer should init
	err = tracing.Init(&tracing.Option{
		CollectorType:   "zipkin",
		CollectorTarget: "http://localhost:9411/api/v1/spans",
		ServiceName:     "test",
		ProtocolEndpointMap: map[string]string{
			common.ProtocolRest: "0.0.0.0:0",
		},
	})
	assert.NoError(t, err)
	_, err = tracing.ProviderTracer(common.ProtocolRest)
	assert.NoError(t, err)
}
