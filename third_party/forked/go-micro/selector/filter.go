package selector

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"github.com/ServiceComb/go-chassis/core/registry"
)

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
