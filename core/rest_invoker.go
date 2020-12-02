package core

import (
	"context"
	"fmt"

	"net/http"

	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/pkg/util"
)

//schemas
const (
	HTTP  = "http"
	HTTPS = "https"
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
func (ri *RestInvoker) ContextDo(ctx context.Context, req *http.Request, options ...InvocationOption) (*http.Response, error) {
	if req.URL.Scheme != HTTP && req.URL.Scheme != HTTPS {
		return nil, fmt.Errorf("scheme invalid: %s, only support http(s)://", req.URL.Scheme)
	}
	common.SetXCSEContext(map[string]string{common.HeaderSourceName: runtime.ServiceName}, req)
	// set headers to Ctx
	if len(req.Header) > 0 {
		m, ok := ctx.Value(common.ContextHeaderKey{}).(map[string]string)
		if !ok {
			m = make(map[string]string)
		}
		ctx = context.WithValue(ctx, common.ContextHeaderKey{}, m)
		for k := range req.Header {
			m[k] = req.Header.Get(k)
		}
	}

	opts := getOpts(options...)
	service, port, err := util.ParseServiceAndPort(req.Host)
	if err != nil {
		return nil, err
	}
	opts.Protocol = common.ProtocolRest
	opts.Port = port

	resp := rest.NewResponse()

	inv := invocation.New(ctx)
	inv.MicroServiceName = service
	wrapInvocationWithOpts(inv, opts)

	//TODO load from openAPI schema
	inv.SchemaID = port
	if inv.SchemaID == "" {
		inv.SchemaID = "rest"
	}
	inv.OperationID = req.URL.Path
	inv.Args = req
	inv.Reply = resp
	inv.URLPath = req.URL.Path

	inv.SetMetadata(common.RestMethod, req.Method)

	err = ri.invoke(inv)
	if err == nil {
		setCookieToCache(*inv, getNamespaceFromMetadata(opts.Metadata))
	}
	return resp, err
}
