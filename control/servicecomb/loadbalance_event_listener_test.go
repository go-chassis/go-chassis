package servicecomb_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/control/servicecomb"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/loadbalancer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLbEventError(t *testing.T) {
	servicecomb.Init()
	archaius.Set("cse.loadbalance.strategy.name", loadbalancer.StrategySessionStickiness)
	assert.Equal(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalance.strategy.name", ""))

	archaius.Set("cse.loadbalance.strategy.Server.name", loadbalancer.StrategySessionStickiness)
	config.ReadLBFromArchaius()
	assert.Equal(t, loadbalancer.StrategySessionStickiness, config.GetStrategyName(""))

	archaius.Delete("cse.loadbalance.strategy.name")
	assert.NotEqual(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalance.strategy.name", ""))
}
