package endpoint

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
)

// GetEndpointFromServiceCenter is used to get the endpoint based on appID, microservice and version
func GetEndpointFromServiceCenter(appID, microService, version string) (string, error) {
	var endPoint string

	tags := utiltags.NewDefaultTag(version, appID)
	instances, err := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(runtime.ServiceID, microService, tags)
	if err != nil {
		lager.Logger.Warnf("Get service instance failed, for key: %s:%s:%s",
			appID, microService, version)
		return "", err
	}

	if len(instances) == 0 {
		lager.Logger.Errorf(nil, "No available instance, key: %s:%s:%s",
			appID, microService, version)
		instanceError := fmt.Sprintf("No available instance, key: %s:%s:%s",
			appID, microService, version)
		return "", errors.New(instanceError)
	}

	for _, instance := range instances {
		for _, value := range instance.EndpointsMap {
			if strings.Contains(value, "?") {
				separation := strings.Split(value, "?")
				if separation[1] == "sslEnabled=true" {
					endPoint = "https://" + separation[0]
				} else {
					endPoint = "http://" + separation[0]
				}
			} else {
				endPoint = "http://" + value
			}
		}
	}
	registry.RefreshProviderCache(microService, appID)
	return endPoint, nil
}
