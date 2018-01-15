package handler

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
)

// BizKeeperProviderHandler bizkeeper provider handler
type BizKeeperProviderHandler struct{}

// Handle handler for bizkeeper provider
func (bk *BizKeeperProviderHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	// command 支撑到SchemaID，后面根据情况进行测试
	command := NewHystrixCmd(i.SourceMicroService, common.Provider, i.MicroServiceName, i.SchemaID, i.OperationID)
	cmdConfig := GetHystrixConfig(command, common.Provider)
	hystrix.ConfigureCommand(command, cmdConfig)
	err := hystrix.Do(command, func() error {
		var err error
		chain.Next(i, func(resp *invocation.InvocationResponse) error {
			err = cb(resp)
			return err
		})
		return err
	}, GetFallbackFun(command, common.Provider, i, cb, cmdConfig.ForceFallback))
	//if err is not nil, means fallback is nil, return original err
	if err != nil {
		writeErr(err, cb)
	}
}

// Name returns bizkeeper-provider string
func (bk *BizKeeperProviderHandler) Name() string {
	return "bizkeeper-provider"
}

func newBizKeeperProviderHandler() Handler {
	return &BizKeeperProviderHandler{}
}
