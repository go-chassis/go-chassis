package pilot

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
)

// Close : Close all connection.
func close(r *EnvoyDSClient) error {
	err := r.Close()
	if err != nil {
		lager.Logger.Errorf(err, "Conn close failed.")
		return err
	}
	lager.Logger.Debugf("Conn close success.")
	return nil
}

// filterInstances filter instances
func filterInstances(hs []*Host) []*registry.MicroServiceInstance {
	instances := make([]*registry.MicroServiceInstance, 0)
	for _, h := range hs {
		msi := ToMicroServiceInstance(h)
		instances = append(instances, msi)
	}
	return instances
}
