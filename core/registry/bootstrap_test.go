package registry_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/config/schema"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/registry/mock"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.registry.address", "http://127.0.0.1:30100")
	archaius.Set("servicecomb.service.name", "Client")
	runtime.HostName = "localhost"
	config.MicroserviceDefinition = &model.ServiceSpec{}
	archaius.UnmarshalConfig(config.MicroserviceDefinition)
	config.ReadGlobalConfigFromArchaius()
}

func TestRegisterService(t *testing.T) {
	runtime.Init()

	config.MicroserviceDefinition.Schemas = []string{"schemaId2", "schemaId3", "schemaId4"}

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
	err := schema.SetSchemaInfoByMap(m)
	assert.NoError(t, err)

	// 	case schemaIDs is empty
	registry.RegisterService()
	registry.RegisterServiceInstances()
	assert.NoError(t, err)

}
