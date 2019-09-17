package loadbalancer_test

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
func TestEnable(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "client"))
	t.Log(os.Getenv("CHASSIS_HOME"))
	config.Init()

	LBstr := make(map[string]string)

	LBstr["name"] = "RoundRobin"
	config.GetLoadBalancing().Strategy = LBstr
	loadbalancer.Enable(archaius.GetString("cse.loadbalance.strategy.name", ""))
	assert.Equal(t, "RoundRobin", config.GetLoadBalancing().Strategy["name"])

	LBstr["name"] = ""
	config.GetLoadBalancing().Strategy = LBstr
	loadbalancer.Enable(archaius.GetString("cse.loadbalance.strategy.name", ""))
	assert.Equal(t, "", config.GetLoadBalancing().Strategy["name"])

	LBstr["name"] = "ABC"
	config.GetLoadBalancing().Strategy = LBstr
	loadbalancer.Enable(archaius.GetString("cse.loadbalance.strategy.name", ""))
	assert.Equal(t, "ABC", config.GetLoadBalancing().Strategy["name"])

}

func TestBuildStrategy(t *testing.T) {
	t.Log("testing default selector")
	testData1 := []*registry.MicroService{
		{
			ServiceName: "test1",
			AppID:       "default",
			Level:       "FRONT",
			Version:     "1.0",
			Status:      "UP",
		},
	}
	testData2 := []*registry.MicroServiceInstance{
		{
			InstanceID:   "01",
			HostName:     "test1",
			Status:       "UP",
			EndpointsMap: map[string]string{"rest": "127.0.0.1", "highway": "10.0.0.3:8080"},
		},
		{
			InstanceID:   "02",
			HostName:     "test2",
			Status:       "UP",
			EndpointsMap: map[string]string{"highway": "10.0.0.3:8080"},
		},
	}

	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	config.Init()
	registry.Enable()
	registry.DoRegister()
	t.Log("System init finished")
	sid, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[0])
	assert.NoError(t, err)

	_, _, err = registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[1])
	assert.NoError(t, err)
	loadbalancer.Enable(archaius.GetString("cse.loadbalance.strategy.name", ""))
	registry.Enable()
	registry.DoRegister()
	runtime.ServiceID = sid
	t.Log(runtime.ServiceID)
	time.Sleep(1 * time.Second)

	inv := &invocation.Invocation{
		SourceServiceID:  sid,
		MicroServiceName: "test1",
		RouteTags:        utiltags.Tags{},
	}
	s, err := loadbalancer.BuildStrategy(inv, nil)
	assert.NoError(t, err)
	ins, err := s.Pick()
	t.Log(ins.EndpointsMap)
	assert.NoError(t, err)
	ins, err = s.Pick()
	assert.NoError(t, err)
	t.Log(ins.EndpointsMap)

	inv = &invocation.Invocation{
		SourceServiceID:  sid,
		MicroServiceName: "fake",
		RouteTags:        utiltags.Tags{},
	}
	s, err = loadbalancer.BuildStrategy(inv, nil)
	assert.Error(t, err)
	t.Log(err)
	switch err.(type) {
	case loadbalancer.LBError:
	default:
		t.Log("Should return lb err")
		t.Fail()
	}
	loadbalancer.SetLatency(1*time.Second,
		"127.0.0.1", "service",
		utiltags.NewDefaultTag("1.0", "app"), "rest")
}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
func BenchmarkDefaultSelector_Select(b *testing.B) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "client"))
	config.Init()
	registry.Enable()
	registry.DoRegister()
	loadbalancer.Enable(archaius.GetString("cse.loadbalance.strategy.name", ""))
	testData1 := []*registry.MicroService{
		{
			ServiceName: "test2",
			AppID:       "default",
			Level:       "FRONT",
			Version:     "1.0",
			Status:      "UP",
		},
	}
	testData2 := []*registry.MicroServiceInstance{
		{
			HostName:     "test1",
			Status:       "UP",
			EndpointsMap: map[string]string{"highway": "10.0.0.4:1234"},
		},
		{
			HostName:     "test2",
			Status:       "UP",
			EndpointsMap: map[string]string{"highway": "10.0.0.3:1234"},
		},
	}
	_, _, _ = registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[0])
	_, _, _ = registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[1])
	time.Sleep(1 * time.Second)
	b.ResetTimer()

	inv := &invocation.Invocation{
		SourceServiceID:  runtime.ServiceID,
		MicroServiceName: "test2",
		RouteTags:        utiltags.Tags{},
	}
	for i := 0; i < b.N; i++ {
		_, _ = loadbalancer.BuildStrategy(inv, nil)
	}

}
