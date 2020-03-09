package registry_test

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	_ "github.com/go-chassis/go-chassis/security/plugins/plain"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
func TestServicecenter_Heartbeat(t *testing.T) {
	goModuleValue := os.Getenv("GO111MODULE")
	rootDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis")
	if goModuleValue == "on" || goModuleValue == "auto" {
		rootDir, _ = os.Getwd()
		rootDir = filepath.Join(rootDir, "..", "..")
	}

	os.Setenv("CHASSIS_HOME", filepath.Join(rootDir, "examples", "discovery", "server"))
	t.Log("Test servercenter.go")
	err := config.Init()
	if err != nil {
		t.Error(err.Error())
	}
	runtime.Init()
	registry.Enable()
	registry.DoRegister()

	microservice := &registry.MicroService{
		AppID:       "CSE",
		ServiceName: "DSFtestAppThree",
		Version:     "2.0.3",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]*registry.Endpoint{"rest": {
			Address:    "10.146.207.197:8080",
			SSLEnabled: false,
		},
			"cse": {
				Address:    "10.146.207.197:8080",
				SSLEnabled: false,
			},
		},
		HostName: "default",
		Status:   common.DefaultStatus,
	}

	sid, insID, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)

	heartBeatService := registry.HeartbeatService{}
	heartBeatService.DoHeartBeat(sid, insID)
	heartBeatService.RetryRegister(sid, insID)
	err = heartBeatService.ReRegisterSelfMSandMSI()
	assert.NoError(t, err)

}

func TestServicecenter_HeartbeatUpdatProperties(t *testing.T) {
	goModuleValue := os.Getenv("GO111MODULE")
	rootDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis")
	if goModuleValue == "on" || goModuleValue == "auto" {
		rootDir, _ = os.Getwd()
		rootDir = filepath.Join(rootDir, "..", "..")
	}
	os.Setenv("CHASSIS_HOME", filepath.Join(rootDir, "examples", "discovery", "server"))
	t.Log("Test servercenter.go")
	config.Init()
	var ins = map[string]string{"type": "test"}
	config.MicroserviceDefinition.ServiceDescription.InstanceProperties = ins
	registry.Enable()
	registry.DoRegister()

	microservice := &registry.MicroService{
		AppID:       "CSE",
		ServiceName: "DSFtestAppThree",
		Version:     "2.0.3",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]*registry.Endpoint{"rest": {
			Address:    "10.146.207.197:8080",
			SSLEnabled: false,
		}},
		HostName: "default",
		Status:   common.DefaultStatus,
	}

	_, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)

}
