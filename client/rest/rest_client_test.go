package rest_test

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	_ "github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	_ "github.com/ServiceComb/go-chassis/server/restful"
	microClient "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	serverOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

var addrRest = "127.0.0.1:8039"

func initEnv() {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "ServiceComb", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition = &model.GlobalCfg{}
}

func TestNewRestClient_Call(t *testing.T) {
	initEnv()
	msName := "Server"
	schema := "schema2"

	trServer := tcp.NewTransport()
	trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain
	strategyRule := make(map[string]string)
	strategyRule["sessionTimeoutInSeconds"] = "30"
	config.GlobalDefinition.Cse.Loadbalance.Strategy = strategyRule

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(
		serverOption.Transport(trServer),
		serverOption.Address(addrRest),
		serverOption.ChainName("default"))
	_, err = s.Register(&schemas.RestFulHello{},
		serverOption.WithMicroServiceName(msName),
		serverOption.WithSchemaID(schema))
	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)

	c := rest.NewRestClient(
		microClient.Transport(trClient),
		microClient.ContentType("application/protobuf"))
	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	reply := rest.NewResponse()
	arg, _ := rest.NewRequest("GET", "cse://Server/instances")
	req := &microClient.Request{
		ID:               1,
		MicroServiceName: "Server",
		Struct:           "",
		Method:           "instances",
		Arg:              arg,
		Metadata:         nil,
	}

	name := c.String()
	log.Println("protocol name:", name)
	options := c.Options()
	log.Println("options are :", options)
	err = c.Call(context.TODO(), addrRest, req, reply)
	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}
	log.Println("hellp reply", &reply)

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	err = c.Call(ctx, addrRest, req, reply)
	expectedError := errors.New("Request Cancelled")
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

func TestNewRestClient_ParseDurationFailed(t *testing.T) {
	t.Log("Testing NewRestClient function for parse duration failed scenario")
	initEnv()
	msName := "Server1"
	schema := "schema2"

	trServer := tcp.NewTransport()
	trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(
		serverOption.Transport(trServer),
		serverOption.Address("127.0.0.1:8040"),
		serverOption.ChainName("default"))
	_, err = s.Register(&schemas.RestFulHello{},
		serverOption.WithMicroServiceName(msName),
		serverOption.WithSchemaID(schema))
	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)

	c := rest.NewRestClient(
		microClient.Transport(trClient),
		microClient.ContentType("application/protobuf"))
	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	reply := rest.NewResponse()
	arg, _ := rest.NewRequest("GET", "cse://Server1/instances")
	req := &microClient.Request{
		ID:               1,
		MicroServiceName: "Server1",
		Struct:           "",
		Method:           "instances",
		Arg:              arg,
		Metadata:         nil,
	}

	name := c.String()
	log.Println("protocol name:", name)
	options := c.Options()
	log.Println("options are :", options)
	err = c.Call(context.TODO(), "127.0.0.1:8040", req, reply)
	log.Println("hellp reply", reply)
	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
	}

}

func TestNewRestClient_Call_Error_Scenarios(t *testing.T) {
	t.Log("Testing NewRestClient call function for error scenarios")
	initEnv()
	msName := "Server2"
	schema := "schema2"

	trServer := tcp.NewTransport()
	trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(
		serverOption.Transport(trServer),
		serverOption.Address("127.0.0.1:8092"),
		serverOption.ChainName("default"))
	_, err = s.Register(&schemas.RestFulHello{},
		serverOption.WithMicroServiceName(msName),
		serverOption.WithSchemaID(schema))
	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)
	fail := make(map[string]bool)
	fail["something"] = false
	c := rest.NewRestClient(
		microClient.Transport(trClient),
		microClient.WithFailure(fail),
		microClient.PoolSize(3))
	c.Init(microClient.ContentType("application/json"))
	reply := rest.NewResponse()
	req := &microClient.Request{
		ID:               1,
		MicroServiceName: "Server2",
		Struct:           "",
		Method:           "instances",
		Arg:              "",
		Metadata:         nil,
	}

	name := c.String()
	log.Println("protocol name:", name)
	options := c.Options()
	log.Println("options are :", options)
	err = c.Call(context.TODO(), "127.0.0.1:8092", req, reply, microClient.WithContentType("application/protobuf"))
	log.Println("hellp reply", reply)
	assert.Error(t, err)
}
func TestNewRequest(t *testing.T) {
	var cl *rest.Client = new(rest.Client)
	i := cl.NewRequest("service", "schemaid", "operationID", "arg")
	assert.Equal(t, i.Method, "operationID")
}
