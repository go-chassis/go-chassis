package handler

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"strings"
)

// constant for bizkeeper-consumer
const (
	Name = "bizkeeper-consumer"
)

// BizKeeperConsumerHandler bizkeeper consumer handler
type BizKeeperConsumerHandler struct{}

// NewHystrixCmd new hystrix command
func NewHystrixCmd(sourceName, protype, servicename, schemaID, OperationID string) string {
	var cmd string
	//for mesher to govern circuit breaker
	if sourceName != "" {
		cmd = strings.Join([]string{sourceName, protype}, ".")
	} else {
		cmd = protype
	}
	if servicename != "" {
		cmd = strings.Join([]string{cmd, servicename}, ".")
	}
	if schemaID != "" {
		cmd = strings.Join([]string{cmd, schemaID}, ".")
	}
	if OperationID != "" {
		cmd = strings.Join([]string{cmd, OperationID}, ".")
	}
	return cmd

}

// Handle function is for to handle the chain
func (bk *BizKeeperConsumerHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	// command 支撑到SchemaID，后面根据情况进行测试
	command := NewHystrixCmd(i.SourceMicroService, common.Consumer, i.MicroServiceName, i.SchemaID, i.OperationID)
	cmdConfig := GetHystrixConfig(command, common.Consumer)
	hystrix.ConfigureCommand(command, cmdConfig)
	err := hystrix.Do(command, func() error {
		var err error
		chain.Next(i, func(resp *invocation.InvocationResponse) error {
			err = cb(resp)
			return err
		})
		return err
	}, GetFallbackFun(command, common.Consumer, i, cb, cmdConfig.ForceFallback))
	//if err is not nil, means fallback is nil, return original err
	if err != nil {
		writeErr(err, cb)
	}
}

// GetHystrixConfig get hystrix config
func GetHystrixConfig(command, t string) hystrix.CommandConfig {
	cmdConfig := hystrix.CommandConfig{}
	cmdConfig.ForceFallback = archaius.GetForceFallback(command, t)
	cmdConfig.TimeoutEnabled = archaius.GetTimeoutEnabled(command, t)
	cmdConfig.Timeout = archaius.GetTimeout(command, t)
	cmdConfig.MaxConcurrentRequests = archaius.GetMaxConcurrentRequests(command, t)
	cmdConfig.ErrorPercentThreshold = archaius.GetErrorPercentThreshold(command, t)
	cmdConfig.RequestVolumeThreshold = archaius.GetRequestVolumeThreshold(command, t)
	cmdConfig.SleepWindow = archaius.GetSleepWindow(command, t)
	cmdConfig.ForceClose = archaius.GetForceClose(command, t)
	cmdConfig.ForceOpen = archaius.GetForceOpen(command, t)
	cmdConfig.CircuitBreakerEnabled = archaius.GetCircuitBreakerEnabled(command, t)
	return cmdConfig
}

// GetFallbackFun get fallback function
func GetFallbackFun(cmd, t string, i *invocation.Invocation, cb invocation.ResponseCallBack, isForce bool) func(error) error {
	enabled := archaius.GetFallbackEnabled(cmd, t)
	if enabled || isForce {
		return func(err error) error {
			if err.Error() == hystrix.ErrForceFallback.Error() || err.Error() == hystrix.ErrCircuitOpen.Error() || err.Error() == hystrix.ErrMaxConcurrency.Error() || err.Error() == hystrix.ErrTimeout.Error() {
				//进入这里说明断路了，run函数没执行，那么必须callback
				lager.Logger.Error("fallback: "+cmd, err)
				resp := &invocation.InvocationResponse{}
				if archaius.PolicyNull == archaius.GetPolicy(cmd, t) {
					resp.Err = hystrix.FallbackNullError{Message: "return null"}
				} else {
					resp.Err = hystrix.CircuitError{Message: i.MicroServiceName + " is isolated because of error: " + err.Error()}
				}
				cb(resp)
				return nil //没有返回错误的必要
			}
			// 回调函数目前默认都是执行成功的
			return nil
		}
	}
	return nil
}

// newBizKeeperConsumerHandler new bizkeeper consumer handler
func newBizKeeperConsumerHandler() Handler {
	return &BizKeeperConsumerHandler{}
}

// Name is for to represent the name of bizkeeper handler
func (bk *BizKeeperConsumerHandler) Name() string {
	return Name
}
