package core

import (
	"context"
	"fmt"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
)

// RestInvoker is rest invoker
// one invoker for one microservice
// thread safe
type RestInvoker struct {
	*abstractInvoker
}

// NewRestInvoker is gives the object of rest invoker
func NewRestInvoker(opt ...Option) *RestInvoker {
	opts := newOptions(opt...)

	ri := &RestInvoker{
		abstractInvoker: &abstractInvoker{
			opts: opts,
		},
	}
	return ri
}

// ContextDo is for requesting the API
func (ri *RestInvoker) ContextDo(ctx context.Context, req *rest.Request, options ...InvocationOption) (*rest.Response, error) {
	if string(req.GetRequest().URL.Scheme) != "cse" {
		return nil, fmt.Errorf("Scheme invalid: %s, only support cse://", req.GetRequest().URL.Scheme)
	}

	opts := getOpts(req.GetRequest().Host, options...)
	opts.Protocol = common.ProtocolRest

	resp := rest.NewResponse()

	inv := invocation.CreateInvocation()
	wrapInvocationWithOpts(inv, opts)
	inv.MicroServiceName = req.GetRequest().Host
	// TODO load from openAPI schema
	// inv.SchemaID = schemaID
	// inv.OperationID = operationID
	inv.Args = req
	inv.Reply = resp
	inv.Ctx = ctx
	inv.URLPathFormat = req.Req.URL.Path

	if inv.Metadata == nil {
		inv.Metadata = make(map[string]interface{})
	}
	inv.Metadata[common.RestMethod] = req.GetMethod()

	err := ri.invoke(inv)
	return resp, err
}
