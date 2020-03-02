package endpoint_test

import (
	"github.com/go-chassis/go-chassis/core/endpoint"
	_ "github.com/go-chassis/go-chassis/initiator"
	"path/filepath"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/registry"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/stretchr/testify/assert"

	"github.com/go-chassis/go-chassis/core/common"
	"os"
	"testing"
	"time"
)

func TestGetEndpointFromServiceCenterInvalidScenario(t *testing.T) {
	t.Log("Testing GetEndpoint function")
	goModuleValue := os.Getenv("GO111MODULE")
	rootDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis")
	if goModuleValue == "on" || goModuleValue == "auto" {
		rootDir, _ = os.Getwd()
		rootDir = filepath.Join(rootDir, "..", "..")
	}
	os.Setenv("CHASSIS_HOME", filepath.Join(rootDir, "examples", "discovery", "server"))
	chassis.Init()
	registry.Enable()
	_, err := endpoint.GetEndpoint("default", "test", "0.1")
	assert.NotNil(t, err)
}

func TestGetEndpointFromServiceCenterForZeroInstance(t *testing.T) {
	goModuleValue := os.Getenv("GO111MODULE")
	rootDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis")
	if goModuleValue == "on" || goModuleValue == "auto" {
		rootDir, _ = os.Getwd()
		rootDir = filepath.Join(rootDir, "..", "..")
	}
	os.Setenv("CHASSIS_HOME", filepath.Join(rootDir, "examples", "discovery", "server"))
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
	_, err = endpoint.GetEndpoint(microservice.AppID, microservice.ServiceName, microservice.Version)
	assert.NotNil(t, err)
}

func TestGetEndpointFromServiceCenterValidScenario(t *testing.T) {
	goModuleValue := os.Getenv("GO111MODULE")
	rootDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis")
	if goModuleValue == "on" || goModuleValue == "auto" {
		rootDir, _ = os.Getwd()
		rootDir = filepath.Join(rootDir, "..", "..")
	}
	os.Setenv("CHASSIS_HOME", filepath.Join(rootDir, "examples", "discovery", "server"))
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
		EndpointsMap: map[string]*registry.EndPoint{
			"rest": {HostOrIP: "10.146.207.197", Port: "8088", SslEnabled: false},
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
	goModuleValue := os.Getenv("GO111MODULE")
	rootDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis")
	if goModuleValue == "on" || goModuleValue == "auto" {
		rootDir, _ = os.Getwd()
		rootDir = filepath.Join(rootDir, "..", "..")
	}
	os.Setenv("CHASSIS_HOME", filepath.Join(rootDir, "examples", "discovery", "server"))
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
		EndpointsMap: map[string]*registry.EndPoint{
			"rest": {HostOrIP: "10.146.207.197", Port: "8080", SslEnabled: true},
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
	goModuleValue := os.Getenv("GO111MODULE")
	rootDir := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "go-chassis", "go-chassis")
	if goModuleValue == "on" || goModuleValue == "auto" {
		rootDir, _ = os.Getwd()
		rootDir = filepath.Join(rootDir, "..", "..")
	}
	os.Setenv("CHASSIS_HOME", filepath.Join(rootDir, "examples", "discovery", "server"))
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
		EndpointsMap: map[string]*registry.EndPoint{
			"rest": {HostOrIP: "10.146.207.197", Port: "8080", SslEnabled: false},
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
