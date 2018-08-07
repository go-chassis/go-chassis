package archaius_test

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/control/archaius"
	archaius2 "github.com/go-chassis/go-chassis/core/archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
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
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	archaius.SaveToLBCache(&model.LoadBalancing{})
	c, _ := archaius.LBConfigCache.Get("test")
	assert.Equal(t, loadbalancer.StrategyRoundRobin, c.(control.LoadBalancingConfig).Strategy)
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
	err = archaius2.Init()
	assert.NoError(t, err)
	err = control.Init()
	archaius.SaveToCBCache(config.GetHystrixConfig())
	c, _ := archaius.CBConfigCache.Get("Consumer")
	assert.Equal(t, 1000, c.(hystrix.CommandConfig).Timeout)
}
