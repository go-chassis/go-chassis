package archaius_test

import (
	"os"
	"testing"

	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/control/archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
)

func TestSaveToLBCache(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	archaius.SaveToLBCache(&model.LoadBalancing{
		Strategy: map[string]string{
			"name": loadbalancer.StrategyRoundRobin,
		},
		AnyService: map[string]model.LoadBalancingSpec{
			"test": {
				Strategy: map[string]string{
					"name": loadbalancer.StrategyRoundRobin,
				},
			},
		},
	})
	c, _ := archaius.LBConfigCache.Get("test")
	assert.Equal(t, loadbalancer.StrategyRoundRobin, c.(control.LoadBalancingConfig).Strategy)
}
func TestSaveDefaultToLBCache(t *testing.T) {
	t.Log("==delete outdated key")
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	archaius.SaveToLBCache(&model.LoadBalancing{
		Strategy: map[string]string{
			"name": loadbalancer.StrategyRoundRobin,
		},
		AnyService: map[string]model.LoadBalancingSpec{
			"test": {
				Strategy: map[string]string{
					"name": loadbalancer.StrategyRoundRobin,
				},
			},
		},
	})
	_, ok := archaius.LBConfigCache.Get("test")
	assert.True(t, ok)
	archaius.SaveToLBCache(&model.LoadBalancing{})
	_, ok = archaius.LBConfigCache.Get("test")
	assert.False(t, ok)
}

func TestSaveToCBCache(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.GlobalDefinition = &model.GlobalCfg{
		Panel: model.ControlPanel{
			Infra: "",
		},
	}
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	err := config.Init()
	assert.NoError(t, err)
	err = control.Init()
	archaius.SaveToCBCache(config.GetHystrixConfig())
	c, _ := archaius.CBConfigCache.Get("Consumer")
	assert.Equal(t, 1000, c.(hystrix.CommandConfig).Timeout)
}
