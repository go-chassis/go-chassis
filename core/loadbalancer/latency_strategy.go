package loadbalancer

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"math/rand"
)

// ByDuration is for calculating the duration
type ByDuration []*ProtocolStats

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i].AvgLatency < a[j].AvgLatency }

// variables for latency map, rest and highway requests count
var (
	//ProtocolStatsMap saves all stats for all service's protocol, one protocol has a lot of instances
	ProtocolStatsMap = make(map[string][]*ProtocolStats)
	//maintain different locks since multiple goroutine access the map
	LatencyMapRWMutex sync.RWMutex
	weightedRespMutex sync.Mutex
)

//BuildKey return key of stats map
func BuildKey(microServiceName, version, app, protocol string) string {
	//TODO add more data
	return strings.Join([]string{microServiceName, protocol}, "/")
}

// SetLatency for a instance ,it only save latest 10 stats for instance's protocol
func SetLatency(latency time.Duration, addr, microServiceName, version, app, protocol string) {
	key := BuildKey(microServiceName, version, app, protocol)

	LatencyMapRWMutex.RLock()
	stats, ok := ProtocolStatsMap[key]
	LatencyMapRWMutex.RUnlock()
	if !ok {
		stats = make([]*ProtocolStats, 0)
	}
	exist := false
	for _, v := range stats {
		if v.Addr == addr {
			v.SaveLatency(latency)
			exist = true
		}
	}
	if !exist {
		ps := &ProtocolStats{
			Addr: addr,
		}

		ps.SaveLatency(latency)
		stats = append(stats, ps)
	}
	LatencyMapRWMutex.Lock()
	ProtocolStatsMap[key] = stats
	LatencyMapRWMutex.Unlock()
}

// SortLatency sort instance based on  the average latencies
func SortLatency() {
	LatencyMapRWMutex.RLock()
	for _, v := range ProtocolStatsMap {
		sort.Sort(ByDuration(v))
	}
	LatencyMapRWMutex.RUnlock()

}

// CalculateAvgLatency Calculating the average latency for each instance using the statistics collected,
// key is addr/service/protocol
func CalculateAvgLatency() {
	LatencyMapRWMutex.RLock()
	for _, v := range ProtocolStatsMap {
		for _, stats := range v {
			stats.CalculateAverageLatency()
		}
	}
	LatencyMapRWMutex.RUnlock()
}

// WeightedResponseStrategy is a strategy plugin
type WeightedResponseStrategy struct {
	instances        []*registry.MicroServiceInstance
	mtx              sync.Mutex
	serviceName      string
	protocol         string
	checkValuesExist bool
	avgLatencyMap    map[string]time.Duration
}

func init() {
	ticker := time.NewTicker((30 * time.Second))
	//run routine to prepare data
	go func() {
		for range ticker.C {
			if config.GetLoadBalancing() != nil {
				useLatencyAware := false
				for _, v := range config.GetLoadBalancing().AnyService {
					if v.Strategy["name"] == StrategyLatency {
						useLatencyAware = true
						break
					}
				}
				if config.GetLoadBalancing().Strategy["name"] == StrategyLatency {
					useLatencyAware = true
				}
				if useLatencyAware {
					CalculateAvgLatency()
					SortLatency()
					lager.Logger.Info("Preparing data for Weighted Response Strategy")
				}
			}

		}
	}()
}
func newWeightedResponseStrategy() Strategy {

	return &WeightedResponseStrategy{}
}

// ReceiveData receive data
func (r *WeightedResponseStrategy) ReceiveData(instances []*registry.MicroServiceInstance, serviceName, protocol, sessionID string) {
	r.instances = instances
	r.serviceName = serviceName
	r.protocol = protocol
}

// Pick return instance
func (r *WeightedResponseStrategy) Pick() (*registry.MicroServiceInstance, error) {
	if rand.Intn(100) < 70 {
		var instanceAddr string
		LatencyMapRWMutex.RLock()
		if len(ProtocolStatsMap[BuildKey(r.serviceName, "", "", r.protocol)]) != 0 {
			instanceAddr = ProtocolStatsMap[BuildKey(r.serviceName, "", "", r.protocol)][0].Addr
		}
		LatencyMapRWMutex.RUnlock()
		for _, instance := range r.instances {
			if instanceAddr == instance.EndpointsMap[r.protocol] {
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
