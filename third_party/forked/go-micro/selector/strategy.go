package selector

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"math/rand"
	"sync"
	"time"

	"github.com/ServiceComb/go-chassis/core/registry"
)

var i int

func init() {
	rand.Seed(time.Now().UnixNano())
	i = rand.Int()
}

// Random is a random strategy algorithm for node selection
func Random(instances []*registry.MicroServiceInstance, metadata interface{}) Next {
	return func() (*registry.MicroServiceInstance, error) {
		if len(instances) == 0 {
			return nil, ErrNoneAvailable
		}

		k := rand.Int() % len(instances)
		return instances[k], nil
	}
}

// RoundRobin is a roundrobin strategy algorithm for node selection
func RoundRobin(instances []*registry.MicroServiceInstance, metadata interface{}) Next {
	var mtx sync.Mutex
	return func() (*registry.MicroServiceInstance, error) {
		if len(instances) == 0 {
			return nil, ErrNoneAvailable
		}

		mtx.Lock()
		node := instances[i%len(instances)]
		i++
		mtx.Unlock()

		return node, nil
	}
}
