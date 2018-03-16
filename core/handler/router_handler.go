package handler

import (
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/router"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/metadata"
)

// RouterHandler router handler
type RouterHandler struct{}

// Handle is to handle the router related things
func (ph *RouterHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {

	tags := map[string]string{}
	for k, v := range i.Metadata {
		tags[k] = v.(string)
	}
	tags[common.BuildinTagApp] = config.GlobalDefinition.AppID

	var h map[string]string
	if i.Protocol == "rest" {
		req, _ := i.Args.(*rest.Request)
		//h = req.GetRequest().Header.HeaderMap()
		h = map[string]string{}
		req.GetRequest().Header.VisitAll(func(key, value []byte) {
			h[string(key)] = string(value)
		})

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
