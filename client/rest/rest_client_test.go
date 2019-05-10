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

	"fmt"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/server/restful"
	"github.com/stretchr/testify/assert"
	"net/http"
	"time"
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
	expectedError := client.ErrCanceled
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
	handler.CreateChains(common.Provider, defaultChain)
	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(server.Options{
		Address:   "127.0.0.1:8092",
		ChainName: "default",
	})
	_, err = s.Register(&TestSchema{}, server.WithSchemaID(schema))
	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)
	fail := make(map[string]bool)
	fail["http_500"] = true
	c, _ := rest.NewRestClient(client.Options{
		Failure:  fail,
		PoolSize: 3,
	})
	t.Run("get options, it should success", func(t *testing.T) {
		o := c.GetOptions()
		assert.Equal(t, 3, o.PoolSize)
	})
	t.Run("call API, status should be 200", func(t *testing.T) {
		reply := rest.NewResponse()
		r, err := rest.NewRequest("GET", "http://Server/", nil)
		assert.NoError(t, err)
		req := &invocation.Invocation{
			MicroServiceName: "Server",
			Args:             r,
			Metadata:         nil,
			Ctx: common.NewContext(map[string]string{
				"os": "mac",
			}),
		}
		name := c.String()
		t.Log("protocol plugin name:", name)
		err = c.Call(context.TODO(), "127.0.0.1:8092", req, reply)
		t.Log("hello reply", reply)
		assert.NoError(t, err)
	})
	t.Run("call error API with failure map settings, client should return err,",
		func(t *testing.T) {
			reply := rest.NewResponse()
			r, err := rest.NewRequest("GET", "http://Server/error", nil)
			assert.NoError(t, err)
			req := &invocation.Invocation{
				MicroServiceName: "Server",
				Args:             r,
			}
			err = c.Call(context.TODO(), "127.0.0.1:8092", req, reply)
			t.Log("error reply", reply)
			assert.Error(t, err)
		})
	t.Run("reconfigure client",
		func(t *testing.T) {
			c.ReloadConfigs(client.Options{
				Failure:  fail,
				PoolSize: 3,
				Timeout:  3 * time.Second,
			})
			o := c.GetOptions()
			assert.Equal(t, 3, o.PoolSize)
			assert.Equal(t, 3*time.Second, o.Timeout)
		})

}

type TestSchema struct {
}

func (r *TestSchema) Root(b *restful.Context) {
	b.Write([]byte(fmt.Sprintf("x-forwarded-host %s", b.ReadRequest().Host)))
}

func (r *TestSchema) Error(b *restful.Context) {
	b.WriteHeader(http.StatusInternalServerError)
}

//URLPatterns helps to respond for corresponding API calls
func (r *TestSchema) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/", ResourceFuncName: "Root",
			Returns: []*restful.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/error", ResourceFuncName: "Error",
			Returns: []*restful.Returns{{Code: 500}}},
	}
}
