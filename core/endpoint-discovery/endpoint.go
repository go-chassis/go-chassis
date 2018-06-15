package endpoint

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
)

// GetEndpointFromServiceCenter is used to get the endpoint based on appID, microservice and version
func GetEndpointFromServiceCenter(appID, microService, version string) (string, error) {
	var endPoint string

	tags := registry.NewDefaultTag(version, appID)
	instances, err := registry.DefaultServiceDiscoveryService.FindMicroServiceInstances(config.SelfServiceID, microService, tags)
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

	return endPoint, nil
}
