package archaius_test

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/control/archaius"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/stretchr/testify/assert"
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
