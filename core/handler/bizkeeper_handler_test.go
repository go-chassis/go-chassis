package handler_test

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func initialize() {
	os.Setenv("CHASSIS_HOME", "/tmp/")
	chassisConf := filepath.Join("/tmp/", "conf")
	os.MkdirAll(chassisConf, 0600)
	os.Create(filepath.Join(chassisConf, "chassis.yaml"))
	os.Create(filepath.Join(chassisConf, "microservice.yaml"))
}

func TestNewHystrixCmd(t *testing.T) {
	t.Log("testing hystrix command with various parameter")
	cmd := handler.NewHystrixCmd("vmall", common.Consumer, "Carts", "cartService", "get")
	assert.Equal(t, "vmall.Consumer.Carts.cartService.get", cmd)
	cmd = handler.NewHystrixCmd("", common.Consumer, "Carts", "cartService", "get")
	assert.Equal(t, "Consumer.Carts.cartService.get", cmd)
	cmd = handler.NewHystrixCmd("", common.Consumer, "Carts", "cartService", "")
	assert.Equal(t, "Consumer.Carts.cartService", cmd)
	cmd = handler.NewHystrixCmd("", common.Consumer, "Carts", "", "")
	assert.Equal(t, "Consumer.Carts", cmd)
	cmd = handler.NewHystrixCmd("", common.Consumer, "", "", "")
	assert.Equal(t, "Consumer", cmd)
}

func TestBizKeeperConsumerHandler_Handle(t *testing.T) {
	t.Log("testing bizkeeper consumer handler")
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	initialize()
	config.Init()
	archaius.Init()

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
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	config.Init()
	archaius.Init()

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
