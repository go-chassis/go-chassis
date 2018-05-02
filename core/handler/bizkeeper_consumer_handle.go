package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"io"
	"io/ioutil"
)

// constant for bizkeeper-consumer
const (
	Name = "bizkeeper-consumer"
)

// BizKeeperConsumerHandler bizkeeper consumer handler
type BizKeeperConsumerHandler struct{}

// GetHystrixConfig get hystrix config
func GetHystrixConfig(service, protype string) (string, hystrix.CommandConfig) {
	command := protype
	if service != "" {
		command = strings.Join([]string{protype, service}, ".")
	}
	return command, hystrix.CommandConfig{
		ForceFallback:          config.GetForceFallback(service, protype),
		TimeoutEnabled:         config.GetTimeoutEnabled(service, protype),
		Timeout:                config.GetTimeout(command, protype),
		MaxConcurrentRequests:  config.GetMaxConcurrentRequests(command, protype),
		ErrorPercentThreshold:  config.GetErrorPercentThreshold(command, protype),
		RequestVolumeThreshold: config.GetRequestVolumeThreshold(command, protype),
		SleepWindow:            config.GetSleepWindow(command, protype),
		ForceClose:             config.GetForceClose(service, protype),
		ForceOpen:              config.GetForceOpen(service, protype),
		CircuitBreakerEnabled:  config.GetCircuitBreakerEnabled(command, protype),
	}
}

// Handle function is for to handle the chain
func (bk *BizKeeperConsumerHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	command, cmdConfig := GetHystrixConfig(i.MicroServiceName, common.Consumer)
	hystrix.ConfigureCommand(command, cmdConfig)

	finish := make(chan *invocation.InvocationResponse, 1)
	err := hystrix.Do(command, func() (err error) {
		chain.Next(i, func(resp *invocation.InvocationResponse) error {
			err = resp.Err
			select {
			case finish <- resp:
			default:
				// means hystrix error occurred
			}
			return err
		})
		return
	}, GetFallbackFun(command, common.Consumer, i, finish, cmdConfig.ForceFallback))

	//if err is not nil, means fallback is nil, return original err
	if err != nil {
		writeErr(err, cb)
		return
	}

	cb(<-finish)
}

// GetFallbackFun get fallback function
func GetFallbackFun(cmd, t string, i *invocation.Invocation, finish chan *invocation.InvocationResponse, isForce bool) func(error) error {
	enabled := config.GetFallbackEnabled(cmd, t)
	if enabled || isForce {
		return func(err error) error {
			// if err is type of hystrix error, return a new response
			if err.Error() == hystrix.ErrForceFallback.Error() || err.Error() == hystrix.ErrCircuitOpen.Error() ||
				err.Error() == hystrix.ErrMaxConcurrency.Error() || err.Error() == hystrix.ErrTimeout.Error() {
				// isolation happened, so lead to callback
				lager.Logger.Errorf(err, fmt.Sprintf("fallback for %v", cmd))
				resp := &invocation.InvocationResponse{}

				var code = http.StatusOK
				if config.PolicyNull == config.GetPolicy(i.MicroServiceName, t) {
					resp.Err = hystrix.FallbackNullError{Message: "return null"}
				} else {
					resp.Err = hystrix.CircuitError{Message: i.MicroServiceName + " is isolated because of error: " + err.Error()}
					code = http.StatusRequestTimeout
				}
				switch i.Reply.(type) {
				case *rest.Response:
					resp := i.Reply.(*rest.Response)
					resp.SetStatusCode(code)
					//make sure body is empty
					if resp.GetResponse().Body != nil {
						io.Copy(ioutil.Discard, resp.GetResponse().Body)
						resp.GetResponse().Body.Close()
					}
				}
				select {
				case finish <- resp:
				default:
				}
				return nil //no need to return error
			}
			// call back success
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
