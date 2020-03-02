package circuit

import (
	"errors"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/status"

	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
)

var errNextChainNoResponse = errors.New("hystrix next chain not respond")

// BizKeeperProviderHandler bizkeeper provider handler
type BizKeeperProviderHandler struct{}

// Handle handler for bizkeeper provider
func (bk *BizKeeperProviderHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	command, cmdConfig := control.DefaultPanel.GetCircuitBreaker(*i, common.Provider)
	hystrix.ConfigureCommand(command, cmdConfig)

	var r *invocation.Response
	err := hystrix.Do(command, func() (err error) {
		chain.Next(i, func(resp *invocation.Response) error {
			r = resp
			if resp == nil {
				err = errNextChainNoResponse
			} else {
				err = resp.Err
			}
			return err
		})
		return
	}, nil)

	if err == nil {
		cb(r)
		return
	}

	// when fallback is nil, err not nil only when:
	// 1. chain.Next() is executed and resp.Err is not nil
	// 2. error generated by hystrix mechanism, such as ErrMaxConcurrency / ErrCircuitOpen / ErrForceFallback
	//    in this case chain.Next() is not executed (r == nil)
	if r == nil {
		handler.WriteBackErr(err, status.Status(i.Protocol, status.ServiceUnavailable), cb)
	} else {
		handler.WriteBackErr(r.Err, r.Status, cb)
	}
}

// Name returns bizkeeper-provider string
func (bk *BizKeeperProviderHandler) Name() string {
	return "bizkeeper-provider"
}

func newBizKeeperProviderHandler() handler.Handler {
	return &BizKeeperProviderHandler{}
}
