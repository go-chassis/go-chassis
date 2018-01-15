package handler

import (
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
)

// constants of bizkeeper fake and loadbalance fake
const (
	//BIZKEEPERFAKE & LOADBALANCEFAKE are variables of type string
	BIZKEEPERFAKE   = "bizkeeper-fake"
	LOADBALANCEFAKE = "loadbalance-fake"
)

// BizkeeperFakeHandler fake handler for bizkeeper
type BizkeeperFakeHandler struct{}

// Name 方法 实现
func (bizkeeperfhandler *BizkeeperFakeHandler) Name() string {
	return BIZKEEPERFAKE
}

// Handle 方法 实现
func (bizkeeperfhandler *BizkeeperFakeHandler) Handle(c *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	// 调用Chain.Next(i *invocation.Invocation, f invocation.ResponseCallBack)方法，
	// 执行Chain的下一个handler
	c.Next(i, func(r *invocation.InvocationResponse) error {
		return cb(r)
	})
}

// func() Handler方法 实现
func createBizkeeperFakeHandler() handler.Handler {
	return &BizkeeperFakeHandler{}
}

// LoadbalanceHandlerFake fake handler for loadbalancer
type LoadbalanceHandlerFake struct{}

// Name 方法 实现
func (lbfakehandler *LoadbalanceHandlerFake) Name() string {
	return LOADBALANCEFAKE
}

// Handle 方法 实现
func (lbfakehandler *LoadbalanceHandlerFake) Handle(c *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	// 调用Chain.Next(i *invocation.Invocation, f invocation.ResponseCallBack)方法，
	// 执行Chain的下一个handler
	c.Next(i, func(r *invocation.InvocationResponse) error {
		return cb(r)
	})
}

// func() Handler方法 实现
func createLoadbalanceHandler() handler.Handler {
	return &LoadbalanceHandlerFake{}
}

func init() {
	// 注册handler name和对应的func() Handler方法
	handler.RegisterHandler(BIZKEEPERFAKE, createBizkeeperFakeHandler)
	handler.RegisterHandler(LOADBALANCEFAKE, createLoadbalanceHandler)
}
