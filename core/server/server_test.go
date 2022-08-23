package server_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/lager"

	"errors"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/go-chassis/v2/core/registry/mock"
	"github.com/go-chassis/go-chassis/v2/core/server"
	_ "github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.noRefreshSchema", true)
	config.ReadGlobalConfigFromArchaius()
}

func TestWithOptions(t *testing.T) {
	var o = new(server.Options)
	o.ChainName = "fakechain"
	var md = make(map[string]string)
	md["abc"] = "abc"
	var rego = new(server.RegisterOptions)
	c2 := server.WithSchemaID("schemaid")
	c2(rego)
	assert.Equal(t, "schemaid", rego.SchemaID)

}

const MockError = "movk error"

// TestSrcMgr Test for server_manager.go
func TestSrcMgr(t *testing.T) {

	testRegistryObj := new(mock.RegistratorMock)
	registry.DefaultRegistrator = testRegistryObj
	testRegistryObj.On("UnRegisterMicroServiceInstance", "microServiceID", "microServiceInstanceID").Return(nil)

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""
	defaultChain["default1"] = ""

	var mp model.Protocol
	mp.Listen = "127.0.0.1:0"
	mp.Advertise = "127.0.0.1:8080"
	mp.WorkerNumber = 10
	mp.Transport = "tcp"

	var cseproto map[string]model.Protocol = make(map[string]model.Protocol)

	cseproto["rest"] = mp

	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.ServiceComb.Protocols = cseproto

	server.Init()

	srv := server.GetServers()
	assert.NotNil(t, srv)
	err := server.UnRegistrySelfInstances()
	assert.NoError(t, err)
	err = server.StartServer(server.WithServerMask("fake"))
	assert.NoError(t, err)

	sr, err := server.GetServer("rest")
	assert.NotNil(t, sr)
	assert.NoError(t, err)

	srv1, err := server.GetServerFunc("abc")
	assert.Nil(t, srv1)
	assert.Error(t, err)

	testRegistryObj.On("UnregisterMicroServiceInstance", "microServiceID", "microServiceInstanceID").Return(errors.New(MockError))
	err = server.UnRegistrySelfInstances()

}
func TestSrcMgrErr(t *testing.T) {
	testRegistryObj := new(mock.RegistratorMock)
	registry.DefaultRegistrator = testRegistryObj
	testRegistryObj.On("Unregister instance", "microServiceID", "microServiceInstanceID").Return(errors.New(MockError))

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	var mp model.Protocol
	//mp.Listen="127.0.0.1:30101"
	mp.Advertise = "127.0.0.1:8091"
	mp.WorkerNumber = 10
	mp.Transport = "abc"

	var cseproto = make(map[string]model.Protocol)

	cseproto["rest"] = mp

	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.ServiceComb.Protocols = cseproto

	server.Init()

	srv := server.GetServers()
	assert.NotNil(t, srv)
	err := server.UnRegistrySelfInstances()
	assert.NoError(t, err)
	err = server.StartServer()
	assert.NoError(t, err)

	sr, err := server.GetServer("rest")
	assert.NotNil(t, sr)
	assert.NoError(t, err)

	srv1, err := server.GetServerFunc("abc")
	assert.Nil(t, srv1)
	assert.Error(t, err)

	srv2, err := server.GetServer("abc")
	assert.Nil(t, srv2)
	assert.Error(t, err)

}
