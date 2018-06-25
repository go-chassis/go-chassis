package client

import (
	"errors"

	"github.com/ServiceComb/go-chassis/pkg/istio/util"
	xdsapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
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
