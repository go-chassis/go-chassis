package configcenter_test

import (
	_ "github.com/go-chassis/go-chassis/initiator"

	"github.com/go-chassis/go-chassis/configcenter"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/registry"
	_ "github.com/go-chassis/go-chassis/core/registry/servicecenter"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetConfigCenterEndpoint(t *testing.T) {
	config.GlobalDefinition = &model.GlobalCfg{
		Cse: model.CseStruct{
			Config: model.Config{
				Client: model.ConfigClient{},
			},
		},
	}
	uri, err := configcenter.GetConfigCenterEndpoint()
	assert.NoError(t, err)
	t.Log(uri)
}
func TestInitConfigCenter(t *testing.T) {
	t.Log("Testing InitConfigCenter function")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	err := config.Init()
	registry.Enable()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = ""
	err = configcenter.InitConfigCenter()
	t.Log("HEllo", err)

}

func TestInitConfigCenterWithTenantEmpty(t *testing.T) {
	t.Log("Testing InitConfigCenter function with autodiscovery true and tenant name empty")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	err := config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.Autodiscovery = true
	config.GlobalDefinition.Cse.Config.Client.TenantName = ""
	err = configcenter.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithEmptyURI(t *testing.T) {
	t.Log("Testing InitConfigCenter function with empty ServerURI")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	err := config.Init()

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = ""
	err = configcenter.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithEmptyMicroservice(t *testing.T) {
	t.Log("Testing InitConfigCenter function with empty microservice definition")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	err := config.Init()

	config.MicroserviceDefinition = &model.MicroserviceCfg{}
	err = configcenter.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithEnableSSl(t *testing.T) {
	t.Log("Testing InitConfigCenter function without initializing any parameter")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	err := config.Init()

	err = configcenter.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithInvalidURI(t *testing.T) {
	t.Log("Testing InitConfigCenter function with Invalid URI")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	err := config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = "hdhhhd:njdj"
	config.GlobalDefinition.Cse.Config.Client.Type = "config_center"
	err = configcenter.InitConfigCenter()
	t.Log("HEllo", err)
}

func TestInitConfigCenterWithSSL(t *testing.T) {
	t.Log("Testing InitConfigCenter function with ServerURI https://127.0.0.1:8787")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/examples/discovery/server/")
	err := config.Init()
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Config.Client.ServerURI = "https://127.0.0.1:8787"
	config.GlobalDefinition.Cse.Config.Client.Type = "config_center"
	err = configcenter.InitConfigCenter()
	t.Log("HEllo", err)

	t.Run("with empty DI,should error", func(t *testing.T) {
		config.MicroserviceDefinition.ServiceDescription.Name = ""
		config.MicroserviceDefinition.ServiceDescription.Version = ""
		runtime.App = ""
		err = configcenter.InitConfigCenter()
		assert.Error(t, err)
	})
}
