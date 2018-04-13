package endpoint_test

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/endpoint-discovery"
	"github.com/ServiceComb/go-chassis/core/registry"
	_ "github.com/ServiceComb/go-chassis/core/registry/servicecenter"
	"github.com/stretchr/testify/assert"

	"github.com/ServiceComb/go-chassis/core/common"
	"os"
	"testing"
	"time"
)

func TestGetEndpointFromServiceCenterInvalidScenario(t *testing.T) {
	t.Log("Testing GetEndpointFromServiceCenter function")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	chassis.Init()
	registry.Enable()
	_, err := endpoint.GetEndpointFromServiceCenter("default", "test", "0.1")
	assert.NotNil(t, err)
}

func TestGetEndpointFromServiceCenterForZeroInstance(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	chassis.Init()
	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "FtestAppThreeZero",
		Version:     "2.0.9",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}

	_, err := registry.DefaultRegistrator.RegisterService(microservice)
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
	_, err = endpoint.GetEndpointFromServiceCenter(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.NotNil(t, err)
}

func TestGetEndpointFromServiceCenterValidScenario(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	chassis.Init()
	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "FtestAppThree",
		Version:     "2.0.4",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]string{"rest": "10.146.207.197:8088"},
		HostName:     "default",
		Status:       common.DefaultStatus,
	}

	_, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
	_, err = endpoint.GetEndpointFromServiceCenter(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.Nil(t, err)
}

func TestGetEndpointFromServiceCenterValidScenarioForEnabled(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	chassis.Init()
	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "FtestAppTwo",
		Version:     "2.0.5",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]string{"rest": "10.146.207.197:8080?sslEnabled=true"},
		HostName:     "default",
		Status:       common.DefaultStatus,
	}

	_, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
	_, err = endpoint.GetEndpointFromServiceCenter(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.Nil(t, err)
}

func TestGetEndpointFromServiceCenterValidScenarioForDisabled(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	chassis.Init()
	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "FtestAppOne",
		Version:     "2.0.6",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]string{"rest": "10.146.207.197:8089?sslEnabled=false"},
		HostName:     "default",
		Status:       common.DefaultStatus,
	}

	_, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
	_, err = endpoint.GetEndpointFromServiceCenter(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.Nil(t, err)
}
