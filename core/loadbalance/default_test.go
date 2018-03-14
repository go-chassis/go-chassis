package loadbalance_test

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/registry"
	_ "github.com/ServiceComb/go-chassis/core/registry/servicecenter"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultSelector_Init(t *testing.T) {
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
	sid, _, err := registry.RegistryService.RegisterServiceAndInstance(testData1[0], testData2[0])
	assert.NoError(t, err)

	_, _, err = registry.RegistryService.RegisterServiceAndInstance(testData1[0], testData2[1])
	assert.NoError(t, err)
	loadbalance.Enable()
	registry.Enable()
	registry.DoRegister()
	lb := loadbalance.DefaultSelector
	config.SelfServiceID = sid
	t.Log(config.SelfServiceID)
	time.Sleep(1 * time.Second)
	next, err := lb.Select("test1", common.LatestVersion, selector.WithConsumerID(sid))
	assert.NoError(t, err)
	ins, err := next()
	t.Log(ins.EndpointsMap)
	assert.NoError(t, err)
	ins, err = next()
	assert.NoError(t, err)
	t.Log(ins.EndpointsMap)
	next, err = lb.Select("fakeServer", "0.1", selector.WithAppID("fake"), selector.WithConsumerID(sid))
	assert.Error(t, err)
	t.Log(err)
	switch err.(type) {
	case selector.LBError:
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
	loadbalance.Enable()
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
	_, _, _ = registry.RegistryService.RegisterServiceAndInstance(testData1[0], testData2[0])
	_, _, _ = registry.RegistryService.RegisterServiceAndInstance(testData1[0], testData2[1])
	lb := loadbalance.DefaultSelector
	time.Sleep(1 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lb.Select("test2", "1.0", selector.WithConsumerID(config.SelfServiceID))
	}

}
