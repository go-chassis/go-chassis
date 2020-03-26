package restful

import (
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/server"
	"net/http"
)

// ResourceHandler wraps go-chassis restful function
type ResourceHandler struct {
	handleFunc func(ctx *Context)
	rc         *Context
	opts       server.Options
}

// Handle is to handle the router related things
func (h *ResourceHandler) Handle(chain *handler.Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {
	Invocation2HTTPRequest(inv, h.rc.Req)

	// check body size
	if h.opts.BodyLimit > 0 {
		h.rc.Req.Request.Body = http.MaxBytesReader(h.rc.Resp, h.rc.Req.Request.Body, h.opts.BodyLimit)
	}

	h.rc.Ctx = inv.Ctx
	// call real route func
	h.handleFunc(h.rc)
	ir := &invocation.Response{}
	ir.Status = h.rc.Resp.StatusCode()
	ir.Result = h.rc.Resp
	//call next chain
	cb(ir)
}

func newHandler(f func(ctx *Context), rc *Context, opts server.Options) handler.Handler {
	return &ResourceHandler{
		handleFunc: f,
		rc:         rc,
		opts:       opts,
	}
}

// Name returns the name string
func (h *ResourceHandler) Name() string {
	return "restful"
}
