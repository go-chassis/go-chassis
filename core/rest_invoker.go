package core

import (
	"context"
	"fmt"

	"net/http"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util"
)

const (
	//HTTP is url schema name
	HTTP = "http"
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
	if string(req.URL.Scheme) != "cse" && string(req.URL.Scheme) != HTTP {
		return nil, fmt.Errorf("scheme invalid: %s, only support {cse|http}://", req.URL.Scheme)
	}

	// set headers to Ctx
	if len(req.Header) > 0 {
		m := make(map[string]string, 0)
		ctx = context.WithValue(ctx, common.ContextHeaderKey{}, m)
		for k := range req.Header {
			m[k] = req.Header.Get(k)
		}
	}

	opts := getOpts(req.Host, options...)
	service, port, _ := util.ParseServiceAndPort(req.Host)
	opts.Protocol = common.ProtocolRest
	opts.Port = port

	resp := rest.NewResponse()

	inv := invocation.New(ctx)

	wrapInvocationWithOpts(inv, opts)
	inv.MicroServiceName = service
	//TODO load from openAPI schema
	inv.SchemaID = port
	if inv.SchemaID == "" {
		inv.SchemaID = "rest"
	}
	inv.OperationID = req.URL.Path
	inv.Args = req
	inv.Reply = resp
	inv.URLPathFormat = req.URL.Path

	inv.SetMetadata(common.RestMethod, req.Method)

	err := ri.invoke(inv)

	if err == nil {
		setCookieToCache(*inv, getNamespaceFromMetadata(opts.Metadata))
	}
	return resp, err
}
