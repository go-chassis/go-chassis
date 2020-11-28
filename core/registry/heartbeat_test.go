package registry_test

import (
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	_ "github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	_ "github.com/go-chassis/go-chassis/v2/security/cipher/plugins/plain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
func TestServicecenter_Heartbeat(t *testing.T) {
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
	config.Init()
	var ins = map[string]string{"type": "test"}
	config.MicroserviceDefinition.InstanceProperties = ins
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
