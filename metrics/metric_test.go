package metrics

import (
	//"github.com/ServiceComb/auth"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	//"net/http"
	"os"
	"path/filepath"
	"testing"
)

func initialize() {
	p := filepath.Join(os.Getenv("GOPATH"), "src", "code.huawei.com", "cse", "go-chassis", "examples", "discovery", "server")
	os.Setenv("CHASSIS_HOME", p)
	chassisConf := filepath.Join(p, "conf")
	os.MkdirAll(chassisConf, 0600)

	chassisFilePath := filepath.Join(chassisConf, "chassis.yaml")
	microserviceFilePath := filepath.Join(chassisConf, "microservice.yaml")
	os.Create(chassisFilePath)
	os.Create(microserviceFilePath)

	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
}

func TestInitEmptyServerURI(t *testing.T) {
	//t.Log("Testing metric init function with empty serverURI")
	initialize()
	config.GlobalDefinition = &model.GlobalCfg{}
	baseURL := config.GlobalDefinition.Cse.Monitor.Client.ServerURI
	err := Init()
	if baseURL == "" && err != nil {
		t.Error("Expected failure if Server URI is not present")
	}
}

func TestInitServerURItlsError(t *testing.T) {
	//t.Log("Testing metric init function with serverURI https://csemonitor.com")
	initialize()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Monitor.Client.ServerURI = "https://csemonitor.com"
	_, tlsError := getTLSForClient()
	err := Init()
	if tlsError != nil && err == nil {
		t.Error("Expected failure if failed in GetTlsForClient")
	}
}

func TestInitServerURItlsConfig(t *testing.T) {
	//t.Log("Testing Init function with tlsConfig")
	initialize()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Monitor.Client.ServerURI = "http://csemonitor.com"
	tlsConfig, _ := getTLSForClient()
	Init()
	if tlsConfig != nil {
		t.Error("Expected tlsConfig to be nil if scheme is http")
	}
}

func TestInitServerUriEmptyString(t *testing.T) {
	//t.Log("Testing Init function with ServerURI")
	initialize()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Monitor.Client.ServerURI = ""

	err := Init()
	assert.NoError(t, err)
}

func TestInitUsernameEmpty(t *testing.T) {
	//t.Log("Testing Init function with empty Username")
	initialize()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Monitor.Client.UserName = ""
	err := Init()
	assert.NoError(t, err)
}

func TestInitDomainNameEmpty(t *testing.T) {
	//t.Log("Testing Init function with empty Domain name")
	initialize()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Monitor.Client.DomainName = ""
	err := Init()
	assert.NoError(t, err)
}

func TestInitGenAuthHeader(t *testing.T) {
	//t.Log("Testing Init function with array returned from GenAuthHeader")
	initialize()

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Monitor.Client.ServerURI = "http://csemonitor.com"
	tlsConfig, _ := getTLSForClient()
	Init()
	if tlsConfig != nil {
		t.Error("Expected tlsConfig to be nil if scheme is http")
	}
}
