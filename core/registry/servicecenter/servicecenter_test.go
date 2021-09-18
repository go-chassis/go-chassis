package servicecenter_test

import (
	scregistry "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	_ "github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	_ "github.com/go-chassis/go-chassis/v2/security/cipher/plugins/plain"
	"github.com/go-chassis/sc-client"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
func TestServicecenter_RegisterServiceAndInstance(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	t.Log("Test servercenter.go")
	config.Init()
	runtime.Init()
	t.Log(os.Getenv("CHASSIS_HOME"))
	registry.Enable()
	registry.DoRegister()

	testRegisterServiceAndInstance(t, registry.DefaultRegistrator, registry.DefaultServiceDiscoveryService)

}

func testRegisterServiceAndInstance(t *testing.T, scc registry.Registrator, sd registry.ServiceDiscovery) {
	microservice := &registry.MicroService{
		AppID:       "CSE",
		ServiceName: "DSFtestAppThree",
		Version:     "2.0.3",
		Status:      sc.MicorserviceUp,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]*registry.Endpoint{
			"rest": {
				false,
				"10.146.207.197:8080",
			},
		},
		HostName: "default",
		Status:   sc.MSInstanceUP,
	}
	sid, insID, err := scc.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)
	t.Log("test update")
	scc.UpdateMicroServiceProperties(sid, map[string]string{"test": "test"})
	microservice2, err := sd.GetMicroService(sid)
	assert.NoError(t, err)
	assert.Equal(t, "test", microservice2.Metadata["test"])

	success, err := scc.Heartbeat(sid, insID)
	assert.Equal(t, success, true)
	assert.NoError(t, err)

	_, err = scc.Heartbeat("jdfhbh", insID)
	assert.Error(t, err)

	err = scc.UpdateMicroServiceInstanceStatus(sid, insID, "UP")
	assert.NoError(t, err)

	err = scc.UpdateMicroServiceInstanceProperties(sid, insID, map[string]string{"test": "test"})
	assert.NoError(t, err)
	err = registry.DefaultRegistrator.AddSchemas(sid, "dsfapp.HelloHuawei", "Testschemainfo")
	assert.NoError(t, err)

	heartBeatSvc := registry.HeartbeatService{}
	heartBeatSvc.Stop()
	scc.Close()
}

func TestInstanceWSHeartbeat(t *testing.T) {
	var serviceID, instanceID string
	t.Run("register service & instance, should success", func(t *testing.T) {
		microservice := &registry.MicroService{
			AppID:       "CSE",
			ServiceName: "testService",
			Version:     "2.0.3",
			Status:      sc.MicorserviceUp,
			Level:       "FRONT",
			Schemas:     []string{"dsfapp.HelloHuawei"},
		}
		microServiceInstance := &registry.MicroServiceInstance{
			EndpointsMap: map[string]*registry.Endpoint{
				"rest": {
					false,
					"11.11.11.11:8080",
				},
			},
			HostName: "default",
			Status:   sc.MSInstanceUP,
		}
		sid, insID, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
		assert.NoError(t, err)
		serviceID, instanceID = sid, insID
	})
	t.Run("send heartbeat,should success", func(t *testing.T) {
		success, err := registry.DefaultRegistrator.Heartbeat(serviceID, instanceID)
		assert.Equal(t, success, true)
		assert.NoError(t, err)

		success, err = registry.DefaultRegistrator.WSHeartbeat(serviceID, instanceID, func() {})
		assert.Equal(t, success, true)
		assert.NoError(t, err)
	})
	t.Run("unregister service & instance", func(t *testing.T) {
		err := registry.DefaultRegistrator.UnRegisterMicroServiceInstance(serviceID, instanceID)
		assert.NoError(t, err)
	})
}

func TestRegroupInstances(t *testing.T) {
	keys := []*scregistry.FindService{
		{
			Service: &scregistry.MicroServiceKey{
				ServiceName: "Service1",
			},
		},
		{
			Service: &scregistry.MicroServiceKey{
				ServiceName: "Service2",
			},
		},
		{
			Service: &scregistry.MicroServiceKey{
				ServiceName: "Service3",
			},
		},
	}
	resp := &scregistry.BatchFindInstancesResponse{
		Services: &scregistry.BatchFindResult{
			Updated: []*scregistry.FindResult{
				{Index: 2,
					Instances: []*scregistry.MicroServiceInstance{{
						InstanceId: "1",
					}}},
			},
		},
	}
	m := servicecenter.RegroupInstances(keys, resp)
	t.Log(m)
	assert.Equal(t, 1, len(m["Service3"]))
	assert.Equal(t, 0, len(m["Service1"]))
	assert.Equal(t, 0, len(m["Service2"]))
}
