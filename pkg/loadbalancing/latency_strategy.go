package loadbalancing

import (
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-mesh/openlogging"
)

var i int
var weightedRespMutex sync.Mutex

func init() {
	loadbalancer.InstallStrategy(loadbalancer.StrategyLatency, newWeightedResponseStrategy)
}

// ByDuration is for calculating the duration
type ByDuration []*loadbalancer.ProtocolStats

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i].AvgLatency < a[j].AvgLatency }

// SortLatency sort instance based on  the average latencies
func SortLatency() {
	loadbalancer.LatencyMapRWMutex.RLock()
	for _, v := range loadbalancer.ProtocolStatsMap {
		sort.Sort(ByDuration(v))
	}
	loadbalancer.LatencyMapRWMutex.RUnlock()

}

// CalculateAvgLatency Calculating the average latency for each instance using the statistics collected,
// key is addr/service/protocol
func CalculateAvgLatency() {
	loadbalancer.LatencyMapRWMutex.RLock()
	for _, v := range loadbalancer.ProtocolStatsMap {
		for _, stats := range v {
			stats.CalculateAverageLatency()
		}
	}
	loadbalancer.LatencyMapRWMutex.RUnlock()
}

// WeightedResponseStrategy is a strategy plugin
type WeightedResponseStrategy struct {
	instances        []*registry.MicroServiceInstance
	mtx              sync.Mutex
	serviceName      string
	protocol         string
	checkValuesExist bool
	avgLatencyMap    map[string]time.Duration
	tags             string
}

func init() {
	ticker := time.NewTicker(30 * time.Second)
	//run routine to prepare data
	go func() {
		for range ticker.C {
			if config.GetLoadBalancing() != nil {
				useLatencyAware := false
				for _, v := range config.GetLoadBalancing().AnyService {
					if v.Strategy["name"] == loadbalancer.StrategyLatency {
						useLatencyAware = true
						break
					}
				}
				if config.GetLoadBalancing().Strategy["name"] == loadbalancer.StrategyLatency {
					useLatencyAware = true
				}
				if useLatencyAware {
					CalculateAvgLatency()
					SortLatency()
					openlogging.Info("Preparing data for Weighted Response Strategy")
				}
			}

		}
	}()
}
func newWeightedResponseStrategy() loadbalancer.Strategy {

	return &WeightedResponseStrategy{}
}

// ReceiveData receive data
func (r *WeightedResponseStrategy) ReceiveData(inv *invocation.Invocation, instances []*registry.MicroServiceInstance, serviceKey string) {
	r.instances = instances
	keys := strings.SplitN(serviceKey, "|", 2)
	switch len(keys) {
	case 1:
		r.serviceName = keys[0]
	case 2:
		r.serviceName = keys[0]
		r.tags = keys[1]

	}
	r.protocol = inv.Protocol
}

// Pick return instance
func (r *WeightedResponseStrategy) Pick() (*registry.MicroServiceInstance, error) {
	if rand.Intn(100) < 70 {
		var instanceAddr string
		loadbalancer.LatencyMapRWMutex.RLock()
		if len(loadbalancer.ProtocolStatsMap[loadbalancer.BuildKey(r.serviceName, r.tags, r.protocol)]) != 0 {
			instanceAddr = loadbalancer.ProtocolStatsMap[loadbalancer.BuildKey(r.serviceName, r.tags, r.protocol)][0].Addr
		}
		loadbalancer.LatencyMapRWMutex.RUnlock()
		for _, instance := range r.instances {
			if len(instanceAddr) != 0 && strings.Contains(instance.EndpointsMap[r.protocol], instanceAddr) {
				return instance, nil
			}
		}
	}

	//if no instances are selected round robin will be done
	weightedRespMutex.Lock()
	node := r.instances[i%len(r.instances)]
	i++
	weightedRespMutex.Unlock()
	return node, nil

}
