package endpoint_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/endpoint"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	_ "github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.registry.address", "http://127.0.0.1:30100")
	archaius.Set("servicecomb.service.name", "Client")
	config.ReadGlobalConfigFromArchaius()
	config.MicroserviceDefinition = &config.GlobalDefinition.ServiceComb.ServiceDescription

}
func TestGetEndpointFromServiceCenterInvalidScenario(t *testing.T) {
	t.Log("Testing GetEndpoint function")
	registry.Enable()
	_, err := endpoint.GetEndpoint("default", "test", "0.1")
	assert.NotNil(t, err)
}

func TestGetEndpointFromServiceCenterForZeroInstance(t *testing.T) {
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
	_, err = endpoint.GetEndpoint(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.NotNil(t, err)
}

func TestGetEndpointFromServiceCenterValidScenario(t *testing.T) {
	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "FtestAppThree",
		Version:     "2.0.4",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]*registry.Endpoint{
			"rest": {Address: "10.146.207.197:8080", SSLEnabled: false},
		},
		HostName: "default",
		Status:   common.DefaultStatus,
	}

	_, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
	url, err := endpoint.GetEndpoint(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.Nil(t, err)
	assert.Contains(t, url, "http://")
	t.Logf("url %s", url)
}

func TestGetEndpointFromServiceCenterValidScenarioForEnabled(t *testing.T) {
	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "FtestAppTwo",
		Version:     "2.0.5",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]*registry.Endpoint{
			"rest": {Address: "10.146.207.197:8080", SSLEnabled: true},
		},
		HostName: "default",
		Status:   common.DefaultStatus,
	}

	_, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
	url, err := endpoint.GetEndpoint(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.Nil(t, err)
	assert.Contains(t, url, "https://")
	t.Logf("url %s", url)
}

func TestGetEndpointFromServiceCenterValidScenarioForDisabled(t *testing.T) {
	microservice := &registry.MicroService{
		AppID:       "default",
		ServiceName: "FtestAppOne",
		Version:     "2.0.6",
		Status:      common.DefaultStatus,
		Level:       "FRONT",
		Schemas:     []string{"dsfapp.HelloHuawei"},
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]*registry.Endpoint{
			"rest": {Address: "10.146.207.197:8080", SSLEnabled: false},
		},
		HostName: "default",
		Status:   common.DefaultStatus,
	}

	_, _, err := registry.DefaultRegistrator.RegisterServiceAndInstance(microservice, microServiceInstance)
	time.Sleep(1 * time.Second)
	assert.NoError(t, err)
	url, err := endpoint.GetEndpoint(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.Nil(t, err)
	assert.Contains(t, url, "http://")
	t.Logf("url %s", url)
}
