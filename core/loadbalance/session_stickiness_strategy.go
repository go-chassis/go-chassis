package loadbalance

import (
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	cache "github.com/patrickmn/go-cache"
	"sync"
	"time"
)

var (
	// SessionCache session cache variable
	SessionCache *cache.Cache
	// successiveFailureCount success and failure count
	successiveFailureCount      map[string]int
	successiveFailureCountMutex sync.RWMutex
)

func init() {
	SessionCache = initCache()
	successiveFailureCount = make(map[string]int)
}

//ResetSuccessiveFailureCount reset failure count
func ResetSuccessiveFailureCount(ep string) {
	successiveFailureCountMutex.Lock()
	successiveFailureCount[ep] = 0
	successiveFailureCountMutex.Unlock()
}

//ResetSuccessiveFailureMap make map again
func ResetSuccessiveFailureMap() {
	successiveFailureCountMutex.Lock()
	successiveFailureCount = make(map[string]int)
	successiveFailureCountMutex.Unlock()
}

//IncreaseSuccessiveFailureCount increase failure count
func IncreaseSuccessiveFailureCount(ep string) {
	successiveFailureCountMutex.Lock()
	c, ok := successiveFailureCount[ep]
	if ok {
		successiveFailureCount[ep] = c + 1
		successiveFailureCountMutex.Unlock()
		return
	}
	successiveFailureCount[ep] = 1
	successiveFailureCountMutex.Unlock()
	return
}

//GetSuccessiveFailureCount get failure count
func GetSuccessiveFailureCount(ep string) int {
	successiveFailureCountMutex.RLock()
	defer successiveFailureCountMutex.RUnlock()
	return successiveFailureCount[ep]
}
func initCache() *cache.Cache {
	var value *cache.Cache

	value = cache.New(3e+10, time.Second*30)
	return value
}

// SessionStickiness is a SessionStickiness strategy algorithm for node selection
func SessionStickiness(instances []*registry.MicroServiceInstance, metadata interface{}) selector.Next {
	var mtx sync.Mutex
	strategyRoundRobinClosur := func() (*registry.MicroServiceInstance, error) {
		if len(instances) == 0 {
			return nil, selector.ErrNoneAvailable
		}

		mtx.Lock()
		node := instances[i%len(instances)]
		i++
		mtx.Unlock()

		return node, nil
	}
	if metadata == nil {
		return strategyRoundRobinClosur
	}

	instanceAddr, ok := SessionCache.Get(metadata.(string))
	if ok {
		return func() (*registry.MicroServiceInstance, error) {
			if len(instances) == 0 {
				return nil, selector.ErrNoneAvailable
			}

			for _, node := range instances {
				mtx.Lock()
				if instanceAddr == node.EndpointsMap["rest"] {
					return node, nil
				}

				mtx.Unlock()
			}
			// if micro service instance goes down then related entry in endpoint map will be deleted,
			//so instead of sending nil, a new instance can be selected using roundrobin
			//
			mtx.Lock()
			nodes := instances[i%len(instances)]
			i++
			mtx.Unlock()
			return nodes, nil
		}
	}

	return strategyRoundRobinClosur
}
