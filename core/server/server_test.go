package server_test

import (
	"crypto/tls"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/registry/mock"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/core/transport"
	_ "github.com/ServiceComb/go-chassis/server/restful"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	serverOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	microTransport "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	_ "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
	"github.com/stretchr/testify/assert"
)

func TestWithOptions(t *testing.T) {
	t.Log("setting various parameter to Server Option ")
	var cmap map[string]codec.Codec = make(map[string]codec.Codec)
	var val codec.Codec

	var o *serverOption.Options = new(serverOption.Options)
	o.ID = "1"
	o.ChainName = "fakechain"

	cmap["firstcodec"] = val

	var t1 microTransport.Transport

	var p provider.Provider

	var md = make(map[string]string)
	md["abc"] = "abc"

	var t2 = time.Second

	var t3 = new(tls.Config)

	c1 := serverOption.WithCodecs(cmap)
	c1(o)
	assert.Equal(t, cmap, o.Codecs)

	c1 = serverOption.Name("abc")
	c1(o)
	assert.Equal(t, "abc", o.Name)

	c1 = serverOption.ID("id")
	c1(o)
	assert.Equal(t, "id", o.ID)

	c1 = serverOption.Version("version")
	c1(o)
	assert.Equal(t, "version", o.Version)

	c1 = serverOption.Address("address")
	c1(o)
	assert.Equal(t, "address", o.Address)

	c1 = serverOption.Advertise("advertise")
	c1(o)
	assert.Equal(t, "advertise", o.Advertise)

	c1 = serverOption.ChainName("chainname")
	c1(o)
	assert.Equal(t, "chainname", o.ChainName)

	c1 = serverOption.Transport(t1)
	c1(o)
	assert.Equal(t, t1, o.Transport)

	c1 = serverOption.Provider(p)
	c1(o)
	assert.Equal(t, p, o.Provider)

	c1 = serverOption.Metadata(md)
	c1(o)
	assert.Equal(t, md, o.Metadata)

	c1 = serverOption.TLSConfig(t3)
	c1(o)
	assert.Equal(t, t3, o.TLSConfig)

	c1 = serverOption.RegisterTTL(t2)
	c1(o)
	assert.Equal(t, t2, o.RegisterTTL)

	t.Log("setting various parameter to server register Option")
	var rego *serverOption.RegisterOptions = new(serverOption.RegisterOptions)

	c2 := serverOption.WithMicroServiceName("ms")
	c2(rego)
	assert.Equal(t, "ms", rego.MicroServiceName)

	c2 = serverOption.WithSchemaID("schemaid")
	c2(rego)
	assert.Equal(t, "schemaid", rego.SchemaID)

	c2 = serverOption.WithGrpcRegister("grpcreg")
	c2(rego)
	assert.Equal(t, "grpcreg", rego.GrpcRegister)

	c2 = serverOption.WithServiceProvider(p)
	c2(rego)
	assert.Equal(t, p, rego.Provider)
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

	testRegistryObj := new(mock.RegistryMock)
	registry.RegistryService = testRegistryObj
	testRegistryObj.On("UnregisterMicroServiceInstance", "microServiceID", "microServiceInstanceID").Return(nil)

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

	a := func(...microTransport.Option) microTransport.Transport {
		var tp microTransport.Transport
		//tp :=tcp.NewTransport()
		return tp
	}

	transport.InstallPlugin("abc", a)

	testRegistryObj := new(mock.RegistryMock)
	registry.RegistryService = testRegistryObj
	//testRegistryObj.On("UnregisterMicroServiceInstance","microServiceID", "microServiceInstanceID").Return(nil)
	testRegistryObj.On("UnregisterMicroServiceInstance", "microServiceID", "microServiceInstanceID").Return(errors.New(MockError))

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
