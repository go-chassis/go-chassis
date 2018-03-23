package handler_test

import (
	"log"
	"os"
	"testing"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
)

func TestCBInit(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/ServiceComb/go-chassis/examples/discovery/server/")

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	archaius.Init()
}

func TestBizKeeperConsumerHandler_Handle(t *testing.T) {
	t.Log("testing bizkeeper consumer handler")

	c := handler.Chain{}
	c.AddHandler(&handler.BizKeeperConsumerHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Consumer["bizkeeperconsumerdefault"] = "bizkeeper-consumer"
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	c.Next(i, func(r *invocation.InvocationResponse) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})
}
func TestBizKeeperProviderHandler_Handle(t *testing.T) {
	t.Log("testing bizkeeper provider handler")

	c := handler.Chain{}
	c.AddHandler(&handler.BizKeeperProviderHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Provider = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Provider["bizkeeperproviderdefault"] = "bizkeeper-provider"
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	c.Next(i, func(r *invocation.InvocationResponse) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})
}

func TestBizKeeperHandler_Names(t *testing.T) {
	bizPro := &handler.BizKeeperProviderHandler{}
	proName := bizPro.Name()
	assert.Equal(t, "bizkeeper-provider", proName)

	bizCon := &handler.BizKeeperConsumerHandler{}
	conName := bizCon.Name()
	assert.Equal(t, "bizkeeper-consumer", conName)

}

func BenchmarkBizKeepConsumerHandler_Handler(b *testing.B) {
	b.Log("benchmark for bizkeeper consumer handler")
	c := handler.Chain{}
	c.AddHandler(&handler.BizKeeperConsumerHandler{})

	inv := &invocation.Invocation{
		MicroServiceName: "fakeService",
		SchemaID:         "schema",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Next(inv, func(r *invocation.InvocationResponse) error {
			assert.NoError(b, r.Err)
			return r.Err
		})
		c.Reset()
	}
}
