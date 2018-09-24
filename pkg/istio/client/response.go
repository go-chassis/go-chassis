package client

import (
	"errors"

	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/go-chassis/go-chassis/pkg/istio/util"
)

// GetRouteConfiguration returns routeconfiguration from discovery response
func GetRouteConfiguration(res *xdsapi.DiscoveryResponse) (*xdsapi.RouteConfiguration, error) {
	if res.TypeUrl != util.RouteType || res.Resources[0].TypeUrl != util.RouteType {
		return nil, errors.New("Invalid typeURL" + res.TypeUrl)
	}

	cla := &xdsapi.RouteConfiguration{}
	err := cla.Unmarshal(res.Resources[0].Value)
	if err != nil {
		return nil, err
	}
	return cla, nil
}

// GetClusterConfiguration returns cluster information from discovery response
func GetClusterConfiguration(res *xdsapi.DiscoveryResponse) ([]xdsapi.Cluster, error) {
	if res.TypeUrl != util.ClusterType {
		return nil, errors.New("Invalid typeURL" + res.TypeUrl)
	}

	var cluster []xdsapi.Cluster
	for _, value := range res.GetResources() {
		cla := &xdsapi.Cluster{}
		err := cla.Unmarshal(value.Value)
		if err != nil {
			return nil, errors.New("unmarshall error")

		}
		cluster = append(cluster, *cla)

	}
	return cluster, nil
}
