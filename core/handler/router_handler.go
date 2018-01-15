package handler

import (
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/route"
)

// RouterHandler router handler
type RouterHandler struct{}

// Handle is to handle the router related things
func (ph *RouterHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	if i.Protocol == "rest" {

		tags := map[string]string{}
		for k, v := range i.Metadata {
			tags[k] = v.(string)
		}
		request, _ := i.Args.(*rest.Request)
		if err := router.Route(request.GetRequest().Header, &registry.SourceInfo{Name: i.SourceMicroService, Tags: tags}, i); err != nil {
			writeErr(err, cb)
		}
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
