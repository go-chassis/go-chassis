package servicecenter_test

import (
	"github.com/go-chassis/go-chassis/core/lager"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/registry"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/scclient"
	_ "github.com/go-chassis/go-chassis/security/plugins/plain"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
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
		Status:      client.MicorserviceUp,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]*registry.Endpoint{
			"rest": {
				false,
				"10.146.207.197",
				"8080",
			},
		},
		HostName: "default",
		Status:   client.MSInstanceUP,
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
