package client

import (
	"github.com/cenkalti/backoff"
	"github.com/go-chassis/go-chassis/pkg/scclient/proto"
	"log"
	"net/url"
	"time"
)

func getBackOff(backoffType string) backoff.BackOff {
	switch backoffType {
	case "Exponential":
		return &backoff.ExponentialBackOff{
			InitialInterval:     1000 * time.Millisecond,
			RandomizationFactor: backoff.DefaultRandomizationFactor,
			Multiplier:          backoff.DefaultMultiplier,
			MaxInterval:         30000 * time.Millisecond,
			MaxElapsedTime:      10000 * time.Millisecond,
			Clock:               backoff.SystemClock,
		}
	case "Constant":
		return backoff.NewConstantBackOff(DefaultRetryTimeout * time.Millisecond)
	case "Zero":
		return &backoff.ZeroBackOff{}
	default:
		return backoff.NewConstantBackOff(DefaultRetryTimeout * time.Millisecond)
	}
}

func getProtocolMap(eps []string) map[string]string {
	m := make(map[string]string)
	for _, ep := range eps {
		u, err := url.Parse(ep)
		if err != nil {
			log.Println("URL error" + err.Error())
			continue
		}
		m[u.Scheme] = u.Host
	}
	return m
}

//RegroupInstances organize raw data to better format
func RegroupInstances(keys []*proto.FindService, response proto.BatchFindInstancesResponse) map[string][]*proto.MicroServiceInstance {
	instanceMap := make(map[string][]*proto.MicroServiceInstance, 0)
	if response.Services != nil {
		for _, result := range response.Services.Updated {
			if len(result.Instances) == 0 {
				continue
			}
			for _, instance := range result.Instances {
				instance.ServiceName = keys[result.Index].Service.ServiceName
				instance.App = keys[result.Index].Service.AppId
				instances, ok := instanceMap[instance.ServiceName]
				if !ok {
					instances = make([]*proto.MicroServiceInstance, 0)
					instanceMap[instance.ServiceName] = instances
				}
				instanceMap[instance.ServiceName] = append(instances, instance)
			}

		}
	}
	return instanceMap
}
