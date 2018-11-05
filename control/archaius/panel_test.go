package archaius_test

import (
	_ "github.com/go-chassis/go-chassis/initiator"

	"github.com/go-chassis/go-chassis/control"
	_ "github.com/go-chassis/go-chassis/control/archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPanel_GetLoadBalancing(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	err := config.Init()
	assert.NoError(t, err)
	config.GlobalDefinition.Panel.Infra = "archaius"
	err = control.Init()
	assert.NoError(t, err)

	t.Log("lb")
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "Server",
	}
	c := control.DefaultPanel.GetLoadBalancing(inv)
	assert.Equal(t, loadbalancer.StrategyRandom, c.Strategy)

	inv = invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "",
	}
	c = control.DefaultPanel.GetLoadBalancing(inv)
	assert.Equal(t, loadbalancer.StrategyLatency, c.Strategy)

	inv = invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "fake",
	}
	c = control.DefaultPanel.GetLoadBalancing(inv)
	assert.Equal(t, loadbalancer.StrategyLatency, c.Strategy)

	t.Log("cb ")
	command, cb := control.DefaultPanel.GetCircuitBreaker(inv, common.Consumer)
	assert.Equal(t, 1000, cb.Timeout)
	assert.Equal(t, "Consumer.fake", command)

	t.Log("rl ")
	inv.MicroServiceName = "Server"
	rl := control.DefaultPanel.GetRateLimiting(inv, common.Consumer)
	assert.Equal(t, 100, rl.Rate)
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.Server", rl.Key)
	assert.Equal(t, true, rl.Enabled)
}

func BenchmarkPanel_GetLoadBalancing(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	config.Init()
	config.GlobalDefinition.Panel.Infra = "archaius"
	control.Init()
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "Server",
	}
	for i := 0; i < b.N; i++ {

		control.DefaultPanel.GetLoadBalancing(inv)

	}
}
func BenchmarkPanel_GetLoadBalancing2(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	config.Init()
	config.GlobalDefinition.Panel.Infra = "archaius"
	control.Init()
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "",
	}
	for i := 0; i < b.N; i++ {

		control.DefaultPanel.GetLoadBalancing(inv)

	}
}
func BenchmarkPanel_GetCircuitBreaker(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	config.Init()
	config.GlobalDefinition.Panel.Infra = "archaius"
	control.Init()
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "",
	}
	for i := 0; i < b.N; i++ {

		control.DefaultPanel.GetCircuitBreaker(inv, common.Consumer)

	}
}
func BenchmarkPanel_GetRateLimiting(b *testing.B) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/client/")
	config.Init()
	config.GlobalDefinition.Panel.Infra = "archaius"
	control.Init()
	inv := invocation.Invocation{
		SourceMicroService: "",
		MicroServiceName:   "",
	}
	for i := 0; i < b.N; i++ {

		control.DefaultPanel.GetRateLimiting(inv, common.Consumer)

	}
}
