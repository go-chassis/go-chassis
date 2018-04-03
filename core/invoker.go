package core

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
)

// newOptions is for updating options
func newOptions(options ...Option) Options {
	opts := DefaultOptions

	for _, o := range options {
		o(&opts)
	}
	if opts.ChainName == "" {
		opts.ChainName = common.DefaultChainName
	}
	return opts
}

// abstract invoker is a common invoker for RPC
type abstractInvoker struct {
	opts Options
}

func (ri *abstractInvoker) invoke(i *invocation.Invocation) error {
	if len(i.Filters) == 0 {
		i.Filters = ri.opts.Filters
	}

	// add self service name into remote context, this value used in provider rate limiter
	i.Ctx = common.WithContext(i.Ctx, common.HeaderSourceName, config.SelfServiceName)

	c, err := handler.GetChain(common.Consumer, ri.opts.ChainName)
	if err != nil {
		lager.Logger.Errorf(err, "Handler chain init err.")
		return err
	}

	c.Next(i, func(ir *invocation.InvocationResponse) error {
		err = ir.Err
		return err
	})
	return err
}
