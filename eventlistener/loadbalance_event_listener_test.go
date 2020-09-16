package eventlistener_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/loadbalancer"
	"github.com/go-chassis/go-chassis/v2/eventlistener"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLbEventError(t *testing.T) {

	eventlistener.Init()
	lbEventListener := &eventlistener.LoadbalancingEventListener{}
	e := &event.Event{EventType: "UPDATE", Key: "servicecomb.loadbalance.strategy.name", Value: "SessionStickiness"}
	lbEventListener.Event(e)
	assert.Equal(t, loadbalancer.StrategySessionStickiness, archaius.GetString("servicecomb.loadbalance.strategy.name", ""))
	assert.Equal(t, loadbalancer.StrategySessionStickiness, config.GetStrategyName("", ""))
	e2 := &event.Event{EventType: "DELETE", Key: "servicecomb.loadbalance.strategy.name", Value: "RoundRobin"}
	lbEventListener.Event(e2)
	archaius.Delete("servicecomb.loadbalance.strategy.name")
	assert.NotEqual(t, loadbalancer.StrategySessionStickiness, archaius.GetString("servicecomb.loadbalancer.strategy.name", ""))

}

func TestLbEvent(t *testing.T) {

	loadbalancer.Enable(archaius.GetString("servicecomb.loadbalance.strategy.name", ""))
	eventlistener.Init()
	archaius.Set("servicecomb.loadbalance.strategy.name", "SessionStickiness")
	lbEventListener := &eventlistener.LoadbalancingEventListener{}
	e := &event.Event{EventType: "UPDATE", Key: "servicecomb.loadbalance.strategy.name", Value: "SessionStickiness"}
	lbEventListener.Event(e)
	assert.Equal(t, loadbalancer.StrategySessionStickiness, archaius.GetString("servicecomb.loadbalance.strategy.name", ""))
	assert.Equal(t, loadbalancer.StrategySessionStickiness, config.GetStrategyName("", ""))
	e2 := &event.Event{EventType: "DELETE", Key: "servicecomb.loadbalance.strategy.name", Value: "RoundRobin"}
	lbEventListener.Event(e2)
	archaius.Delete("servicecomb.loadbalance.strategy.name")
	assert.NotEqual(t, loadbalancer.StrategySessionStickiness, archaius.GetString("servicecomb.loadbalance.strategy.name", ""))

}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.loadbalance.strategy.name", "SessionStickiness")
	config.ReadLBFromArchaius()
}
