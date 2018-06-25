package pilot

import (
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/pkg/istio/util"
	envoy_api_v2_route "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
)

// VirtualHostsToRouteRule translate virtual hosts to route rule
func VirtualHostsToRouteRule(vh *envoy_api_v2_route.VirtualHost) []*model.RouteRule {
	routes := make([]*model.RouteRule, 0, len(vh.Routes))
	for i, v := range vh.Routes {
		if action := v.GetRoute(); action != nil {
			if wc := action.GetWeightedClusters(); wc != nil {
				routes = append(routes, WeightedClustersToRouteRule(wc, i))
				continue
			}
		}
	}
	return routes
}

// WeightedClustersToRouteRule translate weighted clusters to route rule
func WeightedClustersToRouteRule(w *envoy_api_v2_route.WeightedCluster, i int) *model.RouteRule {
	tags := make([]*model.RouteTag, len(w.Clusters))
	for i, c := range w.Clusters {
		tags[i] = &model.RouteTag{
			Tags:   map[string]string{"version": util.ServiceKeyToLabel(c.Name)},
			Weight: int(c.Weight.GetValue()),
		}
	}

	return &model.RouteRule{
		Routes:     tags,
		Precedence: i,
	}
}
