package loadbalancer_test

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"
	"github.com/ServiceComb/go-chassis/core/registry"
	_ "github.com/ServiceComb/go-chassis/core/registry/servicecenter"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEnable(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "client"))
	t.Log(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()

	LBstr := make(map[string]string)

	LBstr["name"] = "RoundRobin"
	config.GetLoadBalancing().Strategy = LBstr
	loadbalancer.Enable()
	assert.Equal(t, "RoundRobin", config.GetLoadBalancing().Strategy["name"])

	LBstr["name"] = ""
	config.GetLoadBalancing().Strategy = LBstr
	loadbalancer.Enable()
	assert.Equal(t, "", config.GetLoadBalancing().Strategy["name"])

	LBstr["name"] = "ABC"
	config.GetLoadBalancing().Strategy = LBstr
	loadbalancer.Enable()
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
			HostName:     "test1",
			Status:       "UP",
			EndpointsMap: map[string]string{"rest": "127.0.0.1", "highway": "10.0.0.3:8080"},
		},
		{
			HostName:     "test2",
			Status:       "UP",
			EndpointsMap: map[string]string{"highway": "10.0.0.3:8080"},
		},
	}

	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	registry.Enable()
	registry.DoRegister()
	t.Log("System init finished")
	sid, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[0])
	assert.NoError(t, err)

	_, _, err = registry.DefaultRegistrator.RegisterServiceAndInstance(testData1[0], testData2[1])
	assert.NoError(t, err)
	loadbalancer.Enable()
	registry.Enable()
	registry.DoRegister()
	config.SelfServiceID = sid
	t.Log(config.SelfServiceID)
	time.Sleep(1 * time.Second)
	s, err := loadbalancer.BuildStrategy(sid, "test1", "", common.LatestVersion, "", "", nil, nil, nil)
	assert.NoError(t, err)
	ins, err := s.Pick()
	t.Log(ins.EndpointsMap)
	assert.NoError(t, err)
	ins, err = s.Pick()
	assert.NoError(t, err)
	t.Log(ins.EndpointsMap)
	s, err = loadbalancer.BuildStrategy(sid, "fake", "", "0.1", "", "", nil, nil, nil)
	assert.Error(t, err)
	t.Log(err)
	switch err.(type) {
	case loadbalancer.LBError:
	default:
		t.Log("Should return lb err")
		t.Fail()
	}
}

func BenchmarkDefaultSelector_Select(b *testing.B) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "client"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	registry.Enable()
	registry.DoRegister()
	loadbalancer.Enable()
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
	for i := 0; i < b.N; i++ {
		_, _ = loadbalancer.BuildStrategy(config.SelfServiceID, "test2", "", "1.0", "", "", nil, nil, nil)
	}

}
