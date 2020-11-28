package loadbalancing

import (
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/loadbalancer"
	"github.com/go-chassis/go-chassis/v2/core/registry"
)

func init() {
	loadbalancer.InstallFilter(loadbalancer.ZoneAware, FilterAvailableZoneAffinity)
}

//FilterAvailableZoneAffinity is a region and zone based Select Filter which will Do the selection of instance in the same region and zone, if not Do the selection of instance in any zone in same region , if not Do the selection of instance in any zone of any region
func FilterAvailableZoneAffinity(old []*registry.MicroServiceInstance, c []*loadbalancer.Criteria) []*registry.MicroServiceInstance {
	var instances []*registry.MicroServiceInstance
	if config.GetDataCenter() == nil {
		return old
	}
	if config.GetDataCenter().Name == "" || config.GetDataCenter().AvailableZone == "" {
		return old // Either no information or partial data center information specified, return all instances
	}

	availableZone := config.GetDataCenter().AvailableZone
	regionName := config.GetDataCenter().Name
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

// FilterByMetadata filter instances based meta data
func FilterByMetadata(old []*registry.MicroServiceInstance, c []*loadbalancer.Criteria) []*registry.MicroServiceInstance {
	var instances []*registry.MicroServiceInstance

	for _, ins := range old {
		if ins.Metadata == nil {
			continue
		}
		//TODO read tags in router.yaml and filter instances based on properties and tags
	}

	return instances
}
