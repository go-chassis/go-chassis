package loadbalance

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/registry"
)

// constant string for zoneaware
const (
	ZoneAware = "zoneaware"
)

// Filters is a map of string and array of *registry.MicroServiceInstance
var Filters map[string]func([]*registry.MicroServiceInstance) []*registry.MicroServiceInstance = make(map[string]func([]*registry.MicroServiceInstance) []*registry.MicroServiceInstance)

// InstallFilter install filter
func InstallFilter(name string, f Filter) {
	Filters[name] = f
}

func init() {
	InstallFilter(ZoneAware, FilterAvailableZoneAffinity)
}

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

//FilterAvailableZoneAffinity is a region and zone based Select Filter which will Do the selection of instance in the same region and zone, if not Do the selection of instance in any zone in same region , if not Do the selection of instance in any zone of any region
func FilterAvailableZoneAffinity(old []*registry.MicroServiceInstance) []*registry.MicroServiceInstance {
	var instances []*registry.MicroServiceInstance
	if config.GlobalDefinition.DataCenter == nil {
		return old
	}
	if config.GlobalDefinition.DataCenter.Name == "" || config.GlobalDefinition.DataCenter.AvailableZone == "" {
		return old // Either no information or partial data center information specified, return all instances
	}

	availableZone := config.GlobalDefinition.DataCenter.AvailableZone
	regionName := config.GlobalDefinition.DataCenter.Name
	instances = getInstancesZoneWise(old, regionName, availableZone)
	if len(instances) == 0 {
		instances = getAvailableInstancesInSameRegion(old, regionName)
		if len(instances) == 0 {
			return old //out of region (multi region) case
		}

		return instances //same region but any available zone case
	}

	return instances //same region and same zone case
}

// getInstancesZoneWise check for the same zone and region
func getInstancesZoneWise(providerInstances []*registry.MicroServiceInstance, region, availableZone string) []*registry.MicroServiceInstance {
	instances := make([]*registry.MicroServiceInstance, 0)
	for _, ins := range providerInstances {
		if ins.DataCenterInfo == nil {
			continue
		}

		if ins.DataCenterInfo.Region == region && ins.DataCenterInfo.AvailableZone == availableZone {
			instances = append(instances, ins)
		}
	}

	return instances
}

// getAvailableInstancesInSameRegion check for available instances in same region
func getAvailableInstancesInSameRegion(providerInstances []*registry.MicroServiceInstance, region string) []*registry.MicroServiceInstance {
	instances := make([]*registry.MicroServiceInstance, 0)
	for _, ins := range providerInstances {
		if ins.DataCenterInfo == nil || ins.DataCenterInfo.Region != region {
			continue
		}

		instances = append(instances, ins)
	}

	return instances
}
