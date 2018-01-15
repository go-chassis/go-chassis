package handler_test

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func initEnv() {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	config.Init()
	archaius.Init()
}

func TestProviderRateLimiterDisable(t *testing.T) {
	t.Log("testing providerratelimiter handler with qps enabled as false")
	initEnv()

	c := handler.Chain{}
	c.AddHandler(&handler.ProviderRateLimiterHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.FlowControl.Provider.QPS.Enabled = false
	i := &invocation.Invocation{
		SourceMicroService: "service1",
		SchemaID:           "schema1",
		OperationID:        "SayHello",
		Args:               &helloworld.HelloRequest{Name: "peter"},
	}
	c.Next(i, func(r *invocation.InvocationResponse) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})

}

func TestProviderRateLimiterHandler_Handle(t *testing.T) {
	t.Log("testing providerratelimiter handler with qps enabled as true")

	initEnv()
	c := handler.Chain{}
	c.AddHandler(&handler.ProviderRateLimiterHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.FlowControl.Provider.QPS.Enabled = true
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

func TestProviderRateLimiterHandler_Handle_SourceMicroService(t *testing.T) {
	t.Log("testing providerratelimiter handler with source microservice and qps enabled as true")

	initEnv()
	c := handler.Chain{}
	c.AddHandler(&handler.ProviderRateLimiterHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.FlowControl.Provider.QPS.Enabled = true
	i := &invocation.Invocation{
		SourceMicroService: "service1",
		SchemaID:           "schema1",
		OperationID:        "SayHello",
		Args:               &helloworld.HelloRequest{Name: "peter"},
	}
	c.Next(i, func(r *invocation.InvocationResponse) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})
}

func TestProviderRateLimiterHandler_Name(t *testing.T) {
	r1 := &handler.ProviderRateLimiterHandler{}
	name := r1.Name()
	assert.Equal(t, "providerratelimiter", name)
}
