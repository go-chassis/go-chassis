package rest_test

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	_ "github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/examples/schemas"
	_ "github.com/go-chassis/go-chassis/server/restful"

	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/stretchr/testify/assert"
)

var addrRest = "127.0.0.1:8039"

func initEnv() {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("CHASSIS_HOME"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition = &model.GlobalCfg{}
}

func TestNewRestClient_Call(t *testing.T) {
	initEnv()
	runtime.ServiceName = "Server"
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

	c, err := rest.NewRestClient(client.Options{})
	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	reply := rest.NewResponse()
	arg, _ := rest.NewRequest("GET", "http://Server/instances", nil)
	req := &invocation.Invocation{
		MicroServiceName: "Server",
		Args:             arg,
		Metadata:         nil,
	}

	log.Println("protocol name:", "rest")
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
	expectedError := rest.ErrCanceled
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

func TestNewRestClient_ParseDurationFailed(t *testing.T) {
	t.Log("Testing NewRestClient function for parse duration failed scenario")
	initEnv()
	runtime.ServiceName = "Server1"
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

	c, err := rest.NewRestClient(client.Options{})
	if err != nil {
		t.Errorf("Unexpected dial err: %v", err)
	}

	reply := rest.NewResponse()
	arg, _ := rest.NewRequest("GET", "http://Server1/instances", nil)
	req := &invocation.Invocation{
		MicroServiceName: "Server1",
		Args:             arg,
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
	runtime.ServiceName = "Server2"
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
	c, _ := rest.NewRestClient(client.Options{
		Failure:  fail,
		PoolSize: 3,
	})
	reply := rest.NewResponse()
	req := &invocation.Invocation{
		MicroServiceName: "Server",
		Args:             "",
		Metadata:         nil,
	}

	name := c.String()
	log.Println("protocol name:", name)
	err = c.Call(context.TODO(), "127.0.0.1:8092", req, reply)
	log.Println("hellp reply", reply)
	assert.Error(t, err)
}
