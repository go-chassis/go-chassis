package registry_test

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

func TestServicecenter_Heartbeat(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	t.Log("Test servercenter.go")
	config.Init()
	t.Log(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	registry.Enable()
	registry.DoRegister()

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

	sid, insID, err := registry.RegistryService.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)

	heartBeatService := registry.HeartbeatService{}
	heartBeatService.DoHeartBeat(sid, insID)
	heartBeatService.RetryRegister(sid)
	err = heartBeatService.ReRegisterSelfMSandMSI()
	assert.NoError(t, err)

}
func TestServicecenter_HeartbeatUpdatProperties(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	t.Log("Test servercenter.go")
	config.Init()
	var ins = make(map[string]string)
	config.MicroserviceDefinition.ServiceDescription.InstanceProperties = ins
	t.Log(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	registry.Enable()
	registry.DoRegister()

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

	_, _, err := registry.RegistryService.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)

}
