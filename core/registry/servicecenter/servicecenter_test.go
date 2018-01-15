package servicecenter_test

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	_ "github.com/ServiceComb/go-chassis/core/registry/servicecenter"
	_ "github.com/ServiceComb/go-chassis/security/plugins/plain"
	"github.com/ServiceComb/go-sc-client/model"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestServicecenter_RegisterServiceAndInstance(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	t.Log("Test servercenter.go")
	config.Init()
	t.Log(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	registry.Enable()
	registry.DoRegister()
	testRegisterServiceAndInstance(t, registry.RegistryService)
	sid := testGetMicroServiceID(t, "CSE", "DSFtestAppThree", "2.0.3", registry.RegistryService)
	registry.Enable()
	registry.DoRegister()
	t.Log("获取依赖的实例")
	instances, err := registry.RegistryService.FindMicroServiceInstances(sid, "CSE", "DSFtestAppThree", "2.0.3")
	assert.NoError(t, err)
	assert.NotZero(t, len(instances))

	err = registry.RegistryService.AddSchemas(sid, "dsfapp.HelloHuawei", "Testschemainfo")
	assert.NoError(t, err)

	microservices, err := registry.RegistryService.GetAllMicroServices()
	assert.NoError(t, err)
	assert.NotZero(t, len(microservices))
}

func testRegisterServiceAndInstance(t *testing.T, scc registry.Registry) {
	microservice := &registry.MicroService{
		AppID:       "CSE",
		ServiceName: "DSFtestAppThree",
		Version:     "2.0.3",
		Status:      model.MicorserviceUp,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]string{"rest": "10.146.207.197:8080"},
		HostName:     "default",
		Status:       model.MSInstanceUP,
		Environment:  common.EnvValueProd,
	}
	sid, insID, err := scc.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)
	t.Log("test update")
	scc.UpdateMicroServiceProperties(sid, map[string]string{"test": "test"})
	microservice2, err := scc.GetMicroService(sid)
	assert.NoError(t, err)
	assert.Equal(t, "test", microservice2.Metadata["test"])

	success, err := scc.Heartbeat(sid, insID)
	assert.Equal(t, success, true)
	assert.NoError(t, err)

	_, err = scc.Heartbeat("jdfhbh", insID)
	assert.Error(t, err)

	ins, err := scc.GetMicroServiceInstances(sid, sid)
	assert.NotZero(t, len(ins))
	assert.NoError(t, err)

	err = scc.UpdateMicroServiceInstanceStatus(sid, insID, "UP")
	assert.NoError(t, err)

	err = scc.UpdateMicroServiceInstanceProperties(sid, insID, map[string]string{"test": "test"})
	assert.NoError(t, err)

	name := scc.String()
	assert.NotEmpty(t, name)

	msdep := &registry.MicroServiceDependency{
		Consumer:  &registry.MicroService{AppID: "CSE", ServiceName: "DSFtestAppThree", Version: "2.0.3"},
		Providers: []*registry.MicroService{{AppID: "CSE", ServiceName: "DSFtestAppThree", Version: "2.0.3"}},
	}
	scc.AddDependencies(msdep)

	heartBeatSvc := registry.HeartbeatService{}
	heartBeatSvc.RefreshTask(sid, insID)
	heartBeatSvc.RemoveTask(sid, insID)
	heartBeatSvc.Stop()
	scc.Close()
}

func testGetMicroServiceID(t *testing.T, appID, microServiceName, version string, scc registry.Registry) string {
	sid, err := scc.GetMicroServiceID(appID, microServiceName, version)
	assert.Nil(t, err)
	//sCenter := servicecenter.Servicecenter{}
	//instances, err := sCenter.GetDependentMicroServiceInstances(appID, microServiceName, version)
	//assert.NotZero(t, len(instances))
	//assert.NoError(t, err)
	//instances, err = sCenter.GetDependentMicroServiceInstances("fakeid", "fakename", "fakeversion")
	//assert.Error(t, err)
	return sid
}
