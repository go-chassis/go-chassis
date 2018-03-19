package metrics

import (
	//"github.com/ServiceComb/auth"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	_ "github.com/ServiceComb/go-chassis/core/registry/servicecenter"
	"github.com/stretchr/testify/assert"
	//"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func initialize() {
	p := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server")
	os.Setenv("CHASSIS_HOME", p)
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
}

func TestInitEmptyServerURI(t *testing.T) {
	//t.Log("Testing metric init function with empty serverURI")
	initialize()
	time.Sleep(1 * time.Second)
	registry.Enable()
	config.GlobalDefinition = &model.GlobalCfg{}
	baseURL := config.GlobalDefinition.Cse.Monitor.Client.ServerURI
	err := Init()
	if baseURL == "" && err != nil {
		t.Error("Expected failure if Server URI is not present")
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
