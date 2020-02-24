package handler

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/go-chassis/core/status"
	"github.com/go-chassis/go-chassis/pkg/runtime"
)

// RouterHandler router handler
type RouterHandler struct{}

// Handle is to handle the router related things
func (ph *RouterHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	tags := map[string]string{}
	for k, v := range i.Metadata {
		tags[k] = v.(string)
	}

	if i.RouteTags.KV != nil {
		for k, v := range i.RouteTags.KV {
			tags[k] = v
		}
	}

	if _, ok := tags[common.BuildinTagApp]; !ok {
		tags[common.BuildinTagApp] = runtime.App
	}

	h := make(map[string]string)
	if i.Ctx != nil {
		at, ok := i.Ctx.Value(common.ContextHeaderKey{}).(map[string]string)
		if ok {
			h = map[string]string(at)
		}
	}

	err := router.Route(h, &registry.SourceInfo{Name: i.SourceMicroService, Tags: tags}, i)
	if err != nil {
		WriteBackErr(err, status.Status(i.Protocol, status.ServiceUnavailable), cb)
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
