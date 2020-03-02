package endpoint

import (
	"errors"
	"fmt"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-mesh/openlogging"
	"net"
)

//GetEndpoint is an API used to get the endpoint of a service in discovery service
//it will only return endpoints of a service
func GetEndpoint(appID, microService, version string) (string, error) {
	var endpoint string
	tags := utiltags.NewDefaultTag(version, appID)
	instances, err := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(runtime.ServiceID, microService, tags)
	if err != nil {
		openlogging.GetLogger().Warnf("Get service instance failed, for key: %s:%s:%s",
			appID, microService, version)
		return "", err
	}

	if len(instances) == 0 {
		instanceError := fmt.Sprintf("No available instance, key: %s:%s:%s",
			appID, microService, version)
		return "", errors.New(instanceError)
	}

	for _, instance := range instances {
		for _, value := range instance.EndpointsMap {
			if value.IsSSLEnable() {
				endpoint = "https://" + net.JoinHostPort(value.HostOrIP, value.Port)
			} else {
				endpoint = "http://" + net.JoinHostPort(value.HostOrIP, value.Port)
			}
		}
	}

	return endpoint, nil
}
