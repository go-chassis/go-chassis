package registry_test

import (
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/schema"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/core/registry/mock"
	"github.com/go-chassis/go-chassis/pkg/runtime"
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

func TestRegisterService(t *testing.T) {
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

	config.MicroserviceDefinition.ServiceDescription.Schemas = []string{"schemaId2", "schemaId3", "schemaId4"}

	testRegistryObj := new(mock.RegistratorMock)
	registry.DefaultRegistrator = testRegistryObj
	testRegistryObj.On("UnRegisterMicroServiceInstance", "microServiceID", "microServiceInstanceID").Return(nil)

	m := make(map[string]string, 0)
	m["id1"] = "schemaInfo1"
	m["id2"] = "schemaInfo2"
	m["id3"] = "schemaInfo3"

	// 	case schemaIDs is empty
	registry.RegisterService()
	registry.RegisterServiceInstances()
	err = schema.SetSchemaInfoByMap(m)
	assert.NoError(t, err)

	// 	case schemaIDs is empty
	registry.RegisterService()
	registry.RegisterServiceInstances()
	assert.NoError(t, err)

}
