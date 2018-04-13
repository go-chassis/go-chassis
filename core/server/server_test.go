package server_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/registry/mock"
	"github.com/ServiceComb/go-chassis/core/server"
	_ "github.com/ServiceComb/go-chassis/server/restful"
	"github.com/stretchr/testify/assert"
)

func TestWithOptions(t *testing.T) {
	t.Log("setting various parameter to Server Option ")

	var o = new(server.Options)
	o.ChainName = "fakechain"

	var md = make(map[string]string)
	md["abc"] = "abc"

	t.Log("setting various parameter to server register Option")
	var rego *server.RegisterOptions = new(server.RegisterOptions)

	c2 := server.WithSchemaID("schemaid")
	c2(rego)
	assert.Equal(t, "schemaid", rego.SchemaID)

}

const MockError = "movk error"

func TestSrcMgr(t *testing.T) {

	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	//config.Init()

	err := config.Init()
	assert.NoError(t, err)

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	var arr = []string{"127.0.0.1", "hgfghfff"}

	registry.SelfInstancesCache.Set("abc", arr, time.Second*30)
	/*a:=func(...transport.Option) transport.Transport{
		//var t transport.Transport
		tp :=tcp.NewTransport()
		return tp
	}

	transport.InstallPlugin("rest",a)*/

	/*f:=func(...server.Option) server.Server{
		var s server.Server

		return s
	}
	server.InstallPlugin("rest",f)*/

	testRegistryObj := new(mock.RegistratorMock)
	registry.DefaultRegistrator = testRegistryObj
	testRegistryObj.On("UnRegisterMicroServiceInstance", "microServiceID", "microServiceInstanceID").Return(nil)

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""
	defaultChain["default1"] = ""

	var mp model.Protocol
	mp.Listen = "127.0.0.1:30100"
	mp.Advertise = "127.0.0.1:8080"
	mp.Failure = "127.0.0.1:8080"
	mp.WorkerNumber = 10
	mp.Transport = "tcp"

	var cseproto map[string]model.Protocol = make(map[string]model.Protocol)

	cseproto["rest"] = mp

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Protocols = cseproto

	server.Init()

	srv := server.GetServers()
	assert.NotNil(t, srv)
	//fmt.Println("SSSSSSSSSss",srv)
	err = server.UnRegistrySelfInstances()
	assert.NoError(t, err)
	//fmt.Println("err",err)
	err = server.StartServer()
	assert.NoError(t, err)
	//fmt.Println("err@@@@@@@",err)

	sr, err := server.GetServer("rest")
	assert.NotNil(t, sr)
	assert.NoError(t, err)
	//fmt.Println("Sr",sr)
	//fmt.Println("err",err)

	srv1, err := server.GetServerFunc("abc")
	assert.Nil(t, srv1)
	assert.Error(t, err)

	testRegistryObj.On("UnregisterMicroServiceInstance", "microServiceID", "microServiceInstanceID").Return(errors.New(MockError))
	err = server.UnRegistrySelfInstances()

}
func TestSrcMgrErr(t *testing.T) {

	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")
	//config.Init()

	err := config.Init()
	assert.NoError(t, err)

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	var arr = []string{"127.0.0.1", "hgfghfff"}

	registry.SelfInstancesCache.Set("abc", arr, time.Second*30)
	registry.SelfInstancesCache.Set("def", "def", time.Second*30)

	/*f:=func(...server.Option) server.Server{
		var s server.Server
		return s
	}
	server.InstallPlugin("protocol",f)*/

	testRegistryObj := new(mock.RegistratorMock)
	registry.DefaultRegistrator = testRegistryObj
	//testRegistryObj.On("UnregisterMicroServiceInstance","microServiceID", "microServiceInstanceID").Return(nil)
	testRegistryObj.On("UnRegisterMicroServiceInstance", "microServiceID", "microServiceInstanceID").Return(errors.New(MockError))

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	var mp model.Protocol
	//mp.Listen="127.0.0.1:30101"
	mp.Advertise = "127.0.0.1:8081"
	mp.Failure = "127.0.0.1:8081"
	mp.WorkerNumber = 10
	mp.Transport = "abc"

	var cseproto map[string]model.Protocol = make(map[string]model.Protocol)

	cseproto["rest"] = mp

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Protocols = cseproto

	server.Init()

	srv := server.GetServers()
	assert.NotNil(t, srv)
	//fmt.Println("SSSSSSSSSss",srv)
	err = server.UnRegistrySelfInstances()
	assert.NoError(t, err)
	//fmt.Println("err",err)
	err = server.StartServer()
	assert.NoError(t, err)

	//fmt.Println("err@@@@@@@",err)

	sr, err := server.GetServer("rest")
	assert.NotNil(t, sr)
	assert.NoError(t, err)
	//fmt.Println("Sr",sr)
	//fmt.Println("err",err)

	srv1, err := server.GetServerFunc("abc")
	assert.Nil(t, srv1)
	assert.Error(t, err)

	srv2, err := server.GetServer("abc")
	assert.Nil(t, srv2)
	assert.Error(t, err)

}
