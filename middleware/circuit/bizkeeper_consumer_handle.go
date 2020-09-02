package circuit

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/control"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/status"
	"github.com/go-chassis/go-chassis/v2/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-chassis/openlog"
)

// constant for bizkeeper-consumer
const (
	Name = "bizkeeper-consumer"
)

// BizKeeperConsumerHandler bizkeeper consumer handler
type BizKeeperConsumerHandler struct{}

// Handle function is for to handle the chain
func (bk *BizKeeperConsumerHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	command, cmdConfig := control.DefaultPanel.GetCircuitBreaker(*i, common.Consumer)

	cmdConfig.MetricsConsumerNum = archaius.GetInt("servicecomb.metrics.circuitMetricsConsumerNum", hystrix.DefaultMetricsConsumerNum)
	hystrix.ConfigureCommand(command, cmdConfig)

	finish := make(chan *invocation.Response, 1)
	f, err := GetFallbackFun(command, common.Consumer, i, finish, cmdConfig.ForceFallback)
	if err != nil {
		handler.WriteBackErr(err, status.Status(i.Protocol, status.InternalServerError), cb)
		return
	}
	err = hystrix.Do(command, func() (err error) {
		chain.Next(i, func(resp *invocation.Response) {
			err = resp.Err
			select {
			case finish <- resp:
			default:
				// means hystrix error occurred
			}
		})
		return
	}, f)

	// err is not nil in conditions:
	// 1 fallback is nil
	//   1.1 chain.Next() fail
	//   1.2 hystrix mechanism, retur error as ErrMaxConcurrency / ErrCircuitOpen / ErrForceFallback
	// 2 fallback is not nil
	//   2.1 fallback failed no matter chain.Next() is executed or not
	if err != nil {
		handler.WriteBackErr(err, status.Status(i.Protocol, status.ServiceUnavailable), cb)
		return
	}

	cb(<-finish)
}

// GetFallbackFun get fallback function
func GetFallbackFun(cmd, t string, i *invocation.Invocation, finish chan *invocation.Response, isForce bool) (func(error) error, error) {
	enabled := config.GetFallbackEnabled(cmd, t)
	if enabled || isForce {
		p := config.GetPolicy(i.MicroServiceName, t)
		if p == "" {
			p = ReturnErr
		}
		f, err := GetFallback(p)
		if err != nil {
			return nil, err
		}
		return f(i, finish), nil
	}
	return nil, nil
}

// newBizKeeperConsumerHandler new bizkeeper consumer handler
func newBizKeeperConsumerHandler() handler.Handler {
	return &BizKeeperConsumerHandler{}
}

// Name is for to represent the name of bizkeeper handler
func (bk *BizKeeperConsumerHandler) Name() string {
	return Name
}

func init() {
	err := handler.RegisterHandler(Name, newBizKeeperConsumerHandler)
	if err != nil {
		openlog.Error(err.Error())
	}
	err = handler.RegisterHandler("bizkeeper-provider", newBizKeeperProviderHandler)
	if err != nil {
		openlog.Error(err.Error())
	}
	Init()
	go hystrix.StartReporter()
}
