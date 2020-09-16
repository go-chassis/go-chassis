package servicecomb_test

import (
	"testing"

	"github.com/go-chassis/go-chassis/v2/control"
	"github.com/go-chassis/go-chassis/v2/control/servicecomb"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/loadbalancer"
	"github.com/go-chassis/go-chassis/v2/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
)

func TestSaveToLBCache(t *testing.T) {
	servicecomb.SaveToLBCache(&model.LoadBalancing{
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
	c, _ := servicecomb.LBConfigCache.Get("test")
	assert.Equal(t, loadbalancer.StrategyRoundRobin, c.(control.LoadBalancingConfig).Strategy)
}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
func TestSaveDefaultToLBCache(t *testing.T) {
	t.Log("==delete outdated key")
	servicecomb.SaveToLBCache(&model.LoadBalancing{
		Strategy: map[string]string{
			"name": loadbalancer.StrategyLatency,
		},
		AnyService: map[string]model.LoadBalancingSpec{
			"test": {
				Strategy: map[string]string{
					"name": loadbalancer.StrategyRoundRobin,
				},
			},
		},
	})
	_, ok := servicecomb.LBConfigCache.Get("test")
	assert.True(t, ok)
	servicecomb.SaveToLBCache(&model.LoadBalancing{})
	_, ok = servicecomb.LBConfigCache.Get("test")
	assert.False(t, ok)
}

func TestSaveToCBCache(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{
		Panel: model.ControlPanel{
			Infra: "",
		},
	}
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	control.Init(opts)
	servicecomb.SaveToCBCache(config.GetHystrixConfig())
	c, _ := servicecomb.CBConfigCache.Get("Consumer")
	assert.Equal(t, 100, c.(hystrix.CommandConfig).MaxConcurrentRequests)
}
