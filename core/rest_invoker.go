package core

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"golang.org/x/net/context"
)

// RestInvoker is rest invoker
// one invoker for one microservice
// thread safe
type RestInvoker struct {
	opts Options
}

// NewRestInvoker is gives the object of rest invoker
func NewRestInvoker(opt ...Option) *RestInvoker {
	opts := newOptions(opt...)

	ri := &RestInvoker{
		opts: opts,
	}
	return ri
}

// ContextDo is for requesting the API
func (ri *RestInvoker) ContextDo(ctx context.Context, req *rest.Request, options ...InvocationOption) (*rest.Response, error) {
	opts := getOpts(string(req.GetRequest().Host()), options...)
	opts.Protocol = common.ProtocolRest
	if len(opts.Filters) == 0 {
		opts.Filters = ri.opts.Filters
	}
	if string(req.GetRequest().URI().Scheme()) != "cse" {
		return nil, fmt.Errorf("Scheme invalid: %s, only support cse://", req.GetRequest().URI().Scheme())
	}
	if req.GetHeader("Content-Type") == "" {
		req.SetHeader("Content-Type", "application/json")
	}
	newReq := req.Copy()
	defer newReq.Close()
	resp := rest.NewResponse()
	newReq.SetHeader(common.HeaderSourceName, config.SelfServiceName)
	inv := invocation.CreateInvocation()
	wrapInvocationWithOpts(inv, opts)
	inv.AppID = config.GlobalDefinition.AppID
	inv.MicroServiceName = string(req.GetRequest().Host())
	inv.Args = newReq
	inv.Reply = resp
	inv.Ctx = ctx
	inv.URLPathFormat = req.Req.URI().String()
	inv.MethodType = req.GetMethod()
	c, err := handler.GetChain(common.Consumer, ri.opts.ChainName)
	if err != nil {
		lager.Logger.Errorf(err, "Handler chain init err.")
		return nil, err
	}
	c.Next(inv, func(ir *invocation.InvocationResponse) error {
		err = ir.Err
		return err
	})
	return resp, err
}
