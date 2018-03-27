package rest_test

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	_ "github.com/ServiceComb/go-chassis/core/loadbalancer"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	_ "github.com/ServiceComb/go-chassis/server/restful"

	"github.com/stretchr/testify/assert"
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
	config.SelfServiceName = "Server"
	schema := "schema2"

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain
	strategyRule := make(map[string]string)
	strategyRule["sessionTimeoutInSeconds"] = "30"
	config.GetLoadBalancing().Strategy = strategyRule

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(
		server.Options{
			Address:   addrRest,
			ChainName: "default",
		})
	_, err = s.Register(&schemas.RestFulHello{},
		server.WithSchemaID(schema))
	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)

	c := rest.NewRestClient(client.Options{})
	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	reply := rest.NewResponse()
	arg, _ := rest.NewRequest("GET", "cse://Server/instances")
	req := &client.Request{
		ID:               1,
		MicroServiceName: "Server",
		Schema:           "",
		Operation:        "instances",
		Arg:              arg,
		Metadata:         nil,
	}

	name := c.String()
	log.Println("protocol name:", name)
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
	config.SelfServiceName = "Server1"
	schema := "schema2"

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(server.Options{
		Address:   "127.0.0.1:8040",
		ChainName: "default",
	})
	_, err = s.Register(&schemas.RestFulHello{},
		server.WithSchemaID(schema))
	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)

	c := rest.NewRestClient(client.Options{})
	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	reply := rest.NewResponse()
	arg, _ := rest.NewRequest("GET", "cse://Server1/instances")
	req := &client.Request{
		ID:               1,
		MicroServiceName: "Server1",
		Schema:           "",
		Operation:        "instances",
		Arg:              arg,
		Metadata:         nil,
	}

	name := c.String()
	log.Println("protocol name:", name)
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
	config.SelfServiceName = "Server2"
	schema := "schema2"

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(server.Options{
		Address:   "127.0.0.1:8092",
		ChainName: "default",
	})
	_, err = s.Register(&schemas.RestFulHello{},
		server.WithSchemaID(schema))
	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)
	fail := make(map[string]bool)
	fail["something"] = false
	c := rest.NewRestClient(client.Options{
		Failure:  fail,
		PoolSize: 3,
	})
	reply := rest.NewResponse()
	req := &client.Request{
		ID:               1,
		MicroServiceName: "Server2",
		Operation:        "instances",
		Arg:              "",
		Metadata:         nil,
	}

	name := c.String()
	log.Println("protocol name:", name)
	err = c.Call(context.TODO(), "127.0.0.1:8092", req, reply)
	log.Println("hellp reply", reply)
	assert.Error(t, err)
}
func TestNewRequest(t *testing.T) {
	i := client.NewRequest("service", "schemaid", "operationID", "arg")
	assert.Equal(t, i.Operation, "operationID")
}
