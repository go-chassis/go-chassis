package handler

import (
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/route"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
)

// RouterHandler router handler
type RouterHandler struct{}

// Handle is to handle the router related things
func (ph *RouterHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {

	tags := map[string]string{}
	for k, v := range i.Metadata {
		tags[k] = v.(string)
	}

	var h map[string]string
	if i.Protocol == "rest" {
		req, _ := i.Args.(*rest.Request)
		h = req.GetRequest().Header.HeaderMap()
	} else {
		ctx, _ := metadata.FromContext(i.Ctx)
		h = map[string]string(ctx)
	}

	err := router.Route(h, &registry.SourceInfo{Name: i.SourceMicroService, Tags: tags}, i)
	if err != nil {
		writeErr(err, cb)
	}

	//call next chain
	chain.Next(i, func(r *invocation.InvocationResponse) error {
		return cb(r)
	})
}

func newRouterHandler() Handler {
	return &RouterHandler{}
}

// Name returns the router string
func (ph *RouterHandler) Name() string {
	return "router"
}
