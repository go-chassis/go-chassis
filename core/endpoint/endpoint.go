package endpoint

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-mesh/openlogging"
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
			if strings.Contains(value, "?") {
				separation := strings.Split(value, "?")
				if separation[1] == "sslEnabled=true" {
					endpoint = "https://" + separation[0]
				} else {
					endpoint = "http://" + separation[0]
				}
			} else {
				endpoint = "https://" + value
			}
		}
	}

	return endpoint, nil
}
