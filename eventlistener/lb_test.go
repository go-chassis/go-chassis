package eventlistener_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/eventlistener"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLbEventError(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")

	config.Init()
	eventlistener.Init()
	archaius.Set("cse.loadbalance.strategy.name", "SessionStickiness")
	lbEventListener := &eventlistener.LoadbalancingEventListener{}
	e := &event.Event{EventType: "UPDATE", Key: "cse.loadbalance.strategy.name", Value: "SessionStickiness"}
	lbEventListener.Event(e)
	assert.Equal(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalance.strategy.name", ""))
	assert.Equal(t, loadbalancer.StrategySessionStickiness, config.GetStrategyName("", ""))
	e2 := &event.Event{EventType: "DELETE", Key: "cse.loadbalance.strategy.name", Value: "RoundRobin"}
	lbEventListener.Event(e2)
	archaius.Delete("cse.loadbalance.strategy.name")
	assert.NotEqual(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalancer.strategy.name", ""))

}

func TestLbEvent(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")

	config.Init()
	loadbalancer.Enable(archaius.GetString("cse.loadbalance.strategy.name", ""))
	eventlistener.Init()
	archaius.Set("cse.loadbalance.strategy.name", "SessionStickiness")
	lbEventListener := &eventlistener.LoadbalancingEventListener{}
	e := &event.Event{EventType: "UPDATE", Key: "cse.loadbalance.strategy.name", Value: "SessionStickiness"}
	lbEventListener.Event(e)
	assert.Equal(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalance.strategy.name", ""))
	assert.Equal(t, loadbalancer.StrategySessionStickiness, config.GetStrategyName("", ""))
	e2 := &event.Event{EventType: "DELETE", Key: "cse.loadbalance.strategy.name", Value: "RoundRobin"}
	lbEventListener.Event(e2)
	archaius.Delete("cse.loadbalance.strategy.name")
	assert.NotEqual(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalance.strategy.name", ""))

}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
