package handler

import (
	"fmt"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"io"
	"io/ioutil"
	"net/http"
)

// constant for bizkeeper-consumer
const (
	Name = "bizkeeper-consumer"
)

// BizKeeperConsumerHandler bizkeeper consumer handler
type BizKeeperConsumerHandler struct{}

// Handle function is for to handle the chain
func (bk *BizKeeperConsumerHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	command, cmdConfig := control.DefaultPanel.GetCircuitBreaker(*i, common.Consumer)
	cmdConfig.MetricsConsumerNum = archaius.GetInt("cse.metrics.circuitMetricsConsumerNum", hystrix.DefaultMetricsConsumerNum)
	hystrix.ConfigureCommand(command, cmdConfig)

	finish := make(chan *invocation.Response, 1)
	err := hystrix.Do(command, func() (err error) {
		chain.Next(i, func(resp *invocation.Response) error {
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
func GetFallbackFun(cmd, t string, i *invocation.Invocation, finish chan *invocation.Response, isForce bool) func(error) error {
	enabled := config.GetFallbackEnabled(cmd, t)
	if enabled || isForce {
		return func(err error) error {
			// if err is type of hystrix error, return a new response
			if err.Error() == hystrix.ErrForceFallback.Error() || err.Error() == hystrix.ErrCircuitOpen.Error() ||
				err.Error() == hystrix.ErrMaxConcurrency.Error() || err.Error() == hystrix.ErrTimeout.Error() {
				// isolation happened, so lead to callback
				lager.Logger.Errorf(fmt.Sprintf("fallback for %v, error [%s]", cmd, err.Error()))
				resp := &invocation.Response{}

				var code = http.StatusOK
				if config.PolicyNull == config.GetPolicy(i.MicroServiceName, t) {
					resp.Err = hystrix.FallbackNullError{Message: "return null"}
				} else {
					resp.Err = hystrix.CircuitError{Message: i.MicroServiceName + " is isolated because of error: " + err.Error()}
					code = http.StatusRequestTimeout
				}
				switch i.Reply.(type) {
				case *http.Response:
					resp := i.Reply.(*http.Response)
					resp.StatusCode = code
					//make sure body is empty
					if resp.Body != nil {
						io.Copy(ioutil.Discard, resp.Body)
						resp.Body.Close()
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
