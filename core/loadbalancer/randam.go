package loadbalancer

import (
	"github.com/ServiceComb/go-chassis/core/registry"
	"math/rand"
	"sync"
)

// RandomStrategy is strategy
type RandomStrategy struct {
	instances []*registry.MicroServiceInstance
	mtx       sync.Mutex
}

func newRandomStrategy() Strategy {
	return &RandomStrategy{}
}

// ReceiveData receive data
func (r *RandomStrategy) ReceiveData(instances []*registry.MicroServiceInstance, serviceName, protocol, sessionID string) {
	r.instances = instances
}

// Pick return instance
func (r *RandomStrategy) Pick() (*registry.MicroServiceInstance, error) {
	if len(r.instances) == 0 {
		return nil, ErrNoneAvailableInstance
	}

	r.mtx.Lock()
	k := rand.Int() % len(r.instances)
	r.mtx.Unlock()
	return r.instances[k], nil

}
