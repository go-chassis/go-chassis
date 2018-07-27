package core

import (
	"context"
	"fmt"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
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
// by default if http status is 5XX, then it will return error
func (ri *RestInvoker) ContextDo(ctx context.Context, req *rest.Request, options ...InvocationOption) (*rest.Response, error) {
	if string(req.GetRequest().URL.Scheme) != "cse" {
		return nil, fmt.Errorf("scheme invalid: %s, only support cse://", req.GetRequest().URL.Scheme)
	}

	opts := getOpts(req.GetRequest().Host, options...)
	opts.Protocol = common.ProtocolRest

	resp := rest.NewResponse()

	inv := invocation.New(ctx)
	wrapInvocationWithOpts(inv, opts)
	inv.MicroServiceName = req.GetRequest().Host
	// TODO load from openAPI schema
	// inv.SchemaID = schemaID
	// inv.OperationID = operationID
	inv.Args = req
	inv.Reply = resp
	inv.URLPathFormat = req.Req.URL.Path

	inv.SetMetadata(common.RestMethod, req.GetMethod())

	err := ri.invoke(inv)
	return resp, err
}
