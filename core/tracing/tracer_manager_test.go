package tracing_test

import (
	"errors"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func test(o map[string]string) (opentracing.Tracer, error) {
	return nil, nil
}
func fake(o map[string]string) (opentracing.Tracer, error) {
	return nil, errors.New("123")
}
func TestTracerManager(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{}
	tracing.InstallTracer("test", test)
	tracing.InstallTracer("fake", fake)
	err := tracing.Init()
	assert.NoError(t, err)
	config.GlobalDefinition.Tracing = model.TracingStruct{
		Tracer: "test",
	}
	err = tracing.Init()
	assert.NoError(t, err)

	config.GlobalDefinition.Tracing = model.TracingStruct{
		Tracer: "fake",
	}
	err = tracing.Init()
	assert.Error(t, err)
}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
