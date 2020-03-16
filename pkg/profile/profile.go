package profile

import (
	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-mesh/openlogging"
)

// const
const (
	msgWriteError = "write to response err: "
)

// Profile contains route rule and discovery
type Profile struct {
	RouteRule map[string][]*config.RouteRule              `json:"routeRule"`
	Discovery map[string][]*registry.MicroServiceInstance `json:"discovery"`
}

// HTTPHandleRouteRuleFunc is a go-restful handler which can expose profile of route rule in http server
func HTTPHandleRouteRuleFunc(req *restful.Request, rep *restful.Response) {
	if err := rep.WriteAsJson(listRouteRule()); err != nil {
		openlogging.Error(msgWriteError + err.Error())
	}
}

// HTTPHandleDiscoveryFunc is a go-restful handler which can expose profile of discovery in http server
func HTTPHandleDiscoveryFunc(req *restful.Request, rep *restful.Response) {
	if err := rep.WriteAsJson(listMicroServiceInstance()); err != nil {
		openlogging.Error(msgWriteError + err.Error())
	}
}

// HTTPHandleProfileFunc is a go-restful handler which can expose all profiles in http server
func HTTPHandleProfileFunc(req *restful.Request, rep *restful.Response) {
	if err := rep.WriteAsJson(newProfile()); err != nil {
		openlogging.Error(msgWriteError + err.Error())
	}
}

func newProfile() Profile {
	return Profile{
		RouteRule: listRouteRule(),
		Discovery: listMicroServiceInstance(),
	}
}

func listRouteRule() map[string][]*config.RouteRule {
	return router.DefaultRouter.ListRouteRule()
}

func listMicroServiceInstance() map[string][]*registry.MicroServiceInstance {
	items := registry.MicroserviceInstanceIndex.FullCache().Items()
	m := make(map[string][]*registry.MicroServiceInstance)
	for k, v := range items {
		m[k] = v.Object.([]*registry.MicroServiceInstance)
	}
	return m
}
