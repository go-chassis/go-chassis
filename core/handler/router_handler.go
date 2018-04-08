package handler

import (
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/router"
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

	h := make(map[string]string)
	if i.Protocol == "rest" {
		req, _ := i.Args.(*rest.Request)
		for k := range req.GetRequest().Header {
			h[k] = req.Req.Header.Get(k)
		}
	} else if i.Ctx != nil {
		at, ok := i.Ctx.Value(common.ContextValueKey{}).(map[string]string)
		if ok {
			h = map[string]string(at)
		}
	}

	err := router.Route(h, &registry.SourceInfo{Name: i.SourceMicroService, Tags: tags}, i)
	if err != nil {
		writeErr(err, cb)
	}

	//call next chain
	chain.Next(i, cb)
}

func newRouterHandler() Handler {
	return &RouterHandler{}
}

// Name returns the router string
func (ph *RouterHandler) Name() string {
	return "router"
}
