package eventlistener_test

import (
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"
	"github.com/ServiceComb/go-chassis/eventlistener"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLbEventError(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")

	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	eventlistener.Init()
	archaius.AddKeyValue("cse.loadbalance.strategy.name", "SessionStickiness")
	lbEventListener := &eventlistener.LoadbalanceEventListener{}
	e := &core.Event{EventType: "UPDATE", Key: "cse.loadbalance.strategy.name", Value: "SessionStickiness"}
	lbEventListener.Event(e)
	assert.Equal(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalance.strategy.name", ""))
	assert.Equal(t, loadbalancer.StrategySessionStickiness, config.GetStrategyName("", ""))
	e2 := &core.Event{EventType: "DELETE", Key: "cse.loadbalance.strategy.name", Value: "RoundRobin"}
	lbEventListener.Event(e2)
	archaius.DeleteKeyValue("cse.loadbalance.strategy.name", "SessionStickiness")
	assert.NotEqual(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalancer.strategy.name", ""))

}

func TestLbEvent(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")

	config.Init()
	loadbalancer.Enable()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	eventlistener.Init()
	archaius.AddKeyValue("cse.loadbalance.strategy.name", "SessionStickiness")
	lbEventListener := &eventlistener.LoadbalanceEventListener{}
	e := &core.Event{EventType: "UPDATE", Key: "cse.loadbalance.strategy.name", Value: "SessionStickiness"}
	lbEventListener.Event(e)
	assert.Equal(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalance.strategy.name", ""))
	assert.Equal(t, loadbalancer.StrategySessionStickiness, config.GetStrategyName("", ""))
	e2 := &core.Event{EventType: "DELETE", Key: "cse.loadbalance.strategy.name", Value: "RoundRobin"}
	lbEventListener.Event(e2)
	archaius.DeleteKeyValue("cse.loadbalance.strategy.name", "SessionStickiness")
	assert.NotEqual(t, loadbalancer.StrategySessionStickiness, archaius.GetString("cse.loadbalance.strategy.name", ""))

}
