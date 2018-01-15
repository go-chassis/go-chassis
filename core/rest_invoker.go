package core

import (
	"fmt"
	"net/url"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"golang.org/x/net/context"
)

// ProtocolName is constant variable for rest
const ProtocolName = "rest"

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
	reqURL, err := url.Parse(req.GetURI())
	if err != nil {
		return nil, err
	}
	opts := getOpts(reqURL.Host, options...)

	if opts.Protocol == "" {
		opts.Protocol = common.ProtocolRest
	}
	if len(opts.Filters) == 0 {
		opts.Filters = ri.opts.Filters
	}
	if reqURL.Scheme != "cse" {
		return nil, fmt.Errorf("Scheme invalid: %s, only support cse://", reqURL.Scheme)
	}
	if req.GetHeader("Content-Type") == "" {
		req.SetHeader("Content-Type", "application/json")
	}
	newReq := req.Copy()
	newReq.SetURI(reqURL.String())
	defer newReq.Close()
	resp := rest.NewResponse()
	newReq.SetHeader(common.HeaderSourceName, config.SelfServiceName)
	inv := invocation.CreateInvocation()
	wrapInvocationWithOpts(inv, opts)
	inv.AppID = config.GlobalDefinition.AppID
	inv.MicroServiceName = reqURL.Host
	inv.Args = newReq
	inv.Reply = resp
	inv.Ctx = ctx
	inv.URLPathFormat = reqURL.RequestURI()
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
