package ratelimiter_test

import (
	"log"
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
	"github.com/go-chassis/go-chassis/middleware/ratelimiter"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}
func initEnv() {

	config.Init()
}

func TestProviderRateLimiterDisable(t *testing.T) {
	t.Log("testing providerratelimiter handler with qps enabled as false")
	initEnv()

	c := handler.Chain{}
	c.AddHandler(&ratelimiter.ProviderRateLimiterHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.FlowControl.Provider.QPS.Enabled = false
	i := &invocation.Invocation{
		SourceMicroService: "service1",
		SchemaID:           "schema1",
		OperationID:        "SayHello",
		Args:               &helloworld.HelloRequest{Name: "peter"},
	}
	c.Next(i, func(r *invocation.Response) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})

}

func TestProviderRateLimiterHandler_Handle(t *testing.T) {
	t.Log("testing providerratelimiter handler with qps enabled as true")

	initEnv()
	c := handler.Chain{}
	c.AddHandler(&ratelimiter.ProviderRateLimiterHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.FlowControl.Provider.QPS.Enabled = true
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}
	c.Next(i, func(r *invocation.Response) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})
}

func TestProviderRateLimiterHandler_Handle_SourceMicroService(t *testing.T) {
	t.Log("testing providerratelimiter handler with source microservice and qps enabled as true")

	initEnv()
	c := handler.Chain{}
	c.AddHandler(&ratelimiter.ProviderRateLimiterHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.FlowControl.Provider.QPS.Enabled = true
	i := &invocation.Invocation{
		SourceMicroService: "service1",
		SchemaID:           "schema1",
		OperationID:        "SayHello",
		Args:               &helloworld.HelloRequest{Name: "peter"},
	}
	c.Next(i, func(r *invocation.Response) error {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
		return r.Err
	})
}
