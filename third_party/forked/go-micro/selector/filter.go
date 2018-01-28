package selector

import (
	"github.com/ServiceComb/go-chassis/core/registry"
)

// constant string for zoneaware
const (
	ZoneAware = "zoneaware"
)

// FilterEndpoint is an endpoint based Select Filter which will
// only return services with the endpoint specified.
func FilterEndpoint(target string) Filter {
	return func(old []*registry.MicroServiceInstance) []*registry.MicroServiceInstance {
		var instances []*registry.MicroServiceInstance
		for _, ins := range old {
			for _, ep := range ins.EndpointsMap {
				if ep == target {
					instances = append(instances, ins)
					break
				}
			}
		}
		return instances
	}
}

// FilterMD is a filtering instances based meta data
func FilterMD(key, val string) Filter {
	return func(old []*registry.MicroServiceInstance) []*registry.MicroServiceInstance {
		var instances []*registry.MicroServiceInstance

		for _, ins := range old {
			if ins.Metadata == nil {
				continue
			}
			if ins.Metadata[key] == val {
				instances = append(instances, ins)
			}
		}

		return instances
	}
}

// FilterProtocol is for filtering the instances based on protocol
func FilterProtocol(protocol string) Filter {
	return func(old []*registry.MicroServiceInstance) []*registry.MicroServiceInstance {
		var instances []*registry.MicroServiceInstance

		for _, ins := range old {
			for p := range ins.EndpointsMap {
				if p == protocol {
					instances = append(instances, ins)
					break
				}

			}

		}

		return instances
	}
}
