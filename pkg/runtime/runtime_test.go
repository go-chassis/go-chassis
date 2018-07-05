package runtime_test

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.MicroserviceDefinition = &model.MicroserviceCfg{
		ServiceDescription: model.MicServiceStruct{
			Hostname: "test",
		},
	}
	err := runtime.Init()
	assert.NoError(t, err)
	assert.Equal(t, "test", runtime.HostName)
	assert.NotEmpty(t, runtime.HostName)
}
