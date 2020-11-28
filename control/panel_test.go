package control_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/control"
	_ "github.com/go-chassis/go-chassis/v2/control/servicecomb"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstallPlugin(t *testing.T) {
	control.InstallPlugin("test", nil)

}
func TestInit(t *testing.T) {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	config.GlobalDefinition = &model.GlobalCfg{
		Panel: model.ControlPanel{
			Infra: "",
		},
	}
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	archaius.Init(archaius.WithMemorySource())
	err := control.Init(opts)
	assert.NoError(t, err)
	opts.Infra = "xxx"
	err = control.Init(opts)
	t.Log(err)
	assert.Error(t, err)
}

func TestNewCircuitCmd(t *testing.T) {
	config.HystrixConfig = &model.HystrixConfigWrapper{
		HystrixConfig: model.HystrixConfig{
			CircuitBreakerProperties: model.CircuitWrapper{
				Scope: "",
			},
		},
	}
	i := invocation.Invocation{
		MicroServiceName: "mall",
		SchemaID:         "rest",
		OperationID:      "/test",
		Endpoint:         "127.0.0.1:8081",
	}
	cmd := control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall.rest./test", cmd)

	config.GetHystrixConfig().CircuitBreakerProperties.Scope = "instance"
	cmd = control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall.127.0.0.1:8081", cmd)

	config.GetHystrixConfig().CircuitBreakerProperties.Scope = "instance-api"
	cmd = control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall.127.0.0.1:8081.rest./test", cmd)

	config.GetHystrixConfig().CircuitBreakerProperties.Scope = "api"
	cmd = control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall.rest./test", cmd)

	config.GetHystrixConfig().CircuitBreakerProperties.Scope = "service"
	cmd = control.NewCircuitName("Consumer", config.GetHystrixConfig().CircuitBreakerProperties.Scope, i)
	assert.Equal(t, "Consumer.mall", cmd)
}
