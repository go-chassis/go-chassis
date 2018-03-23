package loadbalancer

import (
	"github.com/ServiceComb/go-chassis/core/registry"
	"sync"
)

// RoundRobinStrategy is strategy
type RoundRobinStrategy struct {
	instances []*registry.MicroServiceInstance
	mtx       sync.Mutex
}

func newRoundRobinStrategy() Strategy {
	return &RoundRobinStrategy{}
}

//ReceiveData receive data
func (r *RoundRobinStrategy) ReceiveData(instances []*registry.MicroServiceInstance, serviceName, protocol, sessionID string) {
	r.instances = instances
}

//Pick return instance
func (r *RoundRobinStrategy) Pick() (*registry.MicroServiceInstance, error) {
	if len(r.instances) == 0 {
		return nil, ErrNoneAvailableInstance
	}
	r.mtx.Lock()
	instance := r.instances[i%len(r.instances)]
	i++
	r.mtx.Unlock()

	return instance, nil
}
