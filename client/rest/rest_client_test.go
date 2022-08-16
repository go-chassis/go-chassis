package rest_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	_ "github.com/go-chassis/go-chassis/v2/core/loadbalancer"
	_ "github.com/go-chassis/go-chassis/v2/server/restful"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core/client"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/go-chassis/go-chassis/v2/examples/schemas"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/stretchr/testify/assert"
)

var addrRest = "127.0.0.1:8039"

func initEnv() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.noRefreshSchema", true)
	config.ReadGlobalConfigFromArchaius()
	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

}

func TestNewRestClient_Close(t *testing.T) {
	initEnv()
	//runtime.ServiceName = "Server"
	addr := "127.0.0.1:8041"
	schema := "schema2"

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = defaultChain
	strategyRule := make(map[string]string)
	strategyRule["sessionTimeoutInSeconds"] = "30"

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(
		server.Options{
			Address:   addr,
			ChainName: "default",
		})
	_, err = s.Register(&schemas.RestFulHello{},
		server.WithSchemaID(schema))
	assert.NoError(t, err)
	err = s.Start()
	assert.NoError(t, err)

	c, err := rest.NewRestClient(client.Options{})
	assert.Nil(t, err)

	reply := rest.NewResponse()
	arg, _ := rest.NewRequest("POST", "http://Server2/sayhi", nil)
	req := &invocation.Invocation{
		MicroServiceName: "Server2",
		Args:             arg,
		Metadata:         nil,
	}

	t.Log("protocol name:", "rest")
	err = c.Call(context.TODO(), addr, req, reply)
	assert.Nil(t, err)
	t.Logf("help reply: %v", reply)

	err = c.Close()
	assert.Nil(t, err)
	log.Println("close client")
}

func TestNewRestClient_Call(t *testing.T) {
	initEnv()
	runtime.ServiceName = "Server"
	schema := "schema2"

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = defaultChain
	strategyRule := make(map[string]string)
	strategyRule["sessionTimeoutInSeconds"] = "30"

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
	assert.Nil(t, err)
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

	t.Log("protocol name:", "rest")
	err = c.Call(context.TODO(), addrRest, req, reply)
	assert.Nil(t, err)
	t.Logf("help reply: %v", reply)

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	err = c.Call(ctx, addrRest, req, reply)
	expectedError := client.ErrCanceled
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
}

func TestNewRestClient_ParseDurationFailed(t *testing.T) {
	t.Log("Testing NewRestClient function for parse duration failed scenario")
	initEnv()
	runtime.ServiceName = "Server1"
	schema := "schema2"

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = defaultChain

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
	assert.Nil(t, err)

	reply := rest.NewResponse()
	arg, _ := rest.NewRequest("GET", "http://Server1/instances", nil)
	req := &invocation.Invocation{
		MicroServiceName: "Server1",
		Args:             arg,
		Metadata:         nil,
	}

	name := c.String()
	t.Log("protocol name: ", name)
	err = c.Call(context.TODO(), "127.0.0.1:8040", req, reply)
	t.Logf("hellp reply: %v", reply)
	assert.Nil(t, err)
}

func TestNewRestClient_Call_Error_Scenarios(t *testing.T) {
	initEnv()

	t.Run("prepare http server and schema", func(t *testing.T) {
		runtime.ServiceName = "Server2"
		schema := "schema2"
		defaultChain := make(map[string]string)
		defaultChain["default"] = ""
		config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = defaultChain
		config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = defaultChain
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
	})

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
		assert.Equal(t, 200, reply.StatusCode)
		assert.NotEmpty(t, httputil.ReadBody(reply))
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

// URLPatterns helps to respond for corresponding API calls
func (r *TestSchema) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/", ResourceFunc: r.Root,
			Returns: []*restful.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/error", ResourceFunc: r.Error,
			Returns: []*restful.Returns{{Code: 500}}},
	}
}
