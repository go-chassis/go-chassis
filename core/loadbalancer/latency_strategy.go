package loadbalancer

import (
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ServiceComb/go-chassis/core/registry"
)

// ByDuration is for calculating the duration
type ByDuration []time.Duration

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i] < a[j] }

// variables for latency map, rest and highway requests count
var (
	//LatencyMap key is the combination of instance addr and microservice name separated by "/"
	LatencyMap map[string][]time.Duration
	//maintain different locks since multiple goroutine access the map
	LatencyMapRWMutex sync.RWMutex
	avgmtx            sync.RWMutex
	weightedRespMutex sync.Mutex
)

// SetLatency for each requests
func SetLatency(duration time.Duration, addr, microServiceNameAndProtocol string) {
	key := addr + "/" + microServiceNameAndProtocol

	LatencyMapRWMutex.RLock()
	val, ok := LatencyMap[key]
	LatencyMapRWMutex.RUnlock()

	if !ok {
		var durationQueue []time.Duration
		durationQueue = append(durationQueue, duration)
		LatencyMapRWMutex.Lock()
		LatencyMap[key] = durationQueue
		LatencyMapRWMutex.Unlock()
	} else {
		LatencyMapRWMutex.Lock()
		if len(val) < 10 {
			val = append(val, duration)
			LatencyMap[key] = val
		} else { // save only latest 10 data for one micro service's protocol endpoint
			val = val[1:]
			val = append(val, duration)
			LatencyMap[key] = val
		}
		LatencyMapRWMutex.Unlock()
	}
}

// SortingLatencyDuration sorting the average latencies recored for each instance
// and returning the instance addr which has the least latency
func SortingLatencyDuration(serviceAndProtocol string, avgLatencyMap map[string]time.Duration) string {
	var mtx sync.Mutex
	var tempLatencyMap = make(map[string]time.Duration)
	for k, v := range avgLatencyMap {
		epMs := strings.Split(k, "/")
		//comparing the microservice and protocol name
		if (epMs[1] + "/" + epMs[2]) == serviceAndProtocol {
			mtx.Lock()
			tempLatencyMap[epMs[0]] = v
			mtx.Unlock()
		}
	}

	//Inverting maps
	invMap := make(map[time.Duration]string, len(tempLatencyMap))
	for k, v := range tempLatencyMap {
		mtx.Lock()
		invMap[v] = k
		mtx.Unlock()
	}

	//Sorting
	sortedKeys := make([]time.Duration, len(invMap))
	var i int
	for k := range invMap {
		sortedKeys[i] = k
		i++
	}
	sort.Sort(ByDuration(sortedKeys))
	return invMap[sortedKeys[0]]

}

// FindingAvgLatency Calculating the average latency for each instance using the statistics collected,
// key is addr/service/protocol
func FindingAvgLatency(metadata string) (avgMap map[string]time.Duration, protocol string) {
	avgMap = make(map[string]time.Duration)
	LatencyMapRWMutex.RLock()
	defer LatencyMapRWMutex.RUnlock()
	for k, v := range LatencyMap {
		epMs := strings.Split(k, "/")
		//comparing the microservice/protocol name
		if (epMs[1] + "/" + epMs[2]) == metadata {
			protocol = epMs[2]
			var sum time.Duration
			for i := 0; i < len(v); i++ {
				sum = sum + v[i]
			}
			avgmtx.Lock()
			avgMap[k] = time.Duration(sum.Nanoseconds() / int64(len(v)))
			avgmtx.Unlock()
		}
	}
	return avgMap, protocol
}

// WeightedResponseStrategy is a strategy plugin
type WeightedResponseStrategy struct {
	instances        []*registry.MicroServiceInstance
	mtx              sync.Mutex
	protocol         string
	instanceAddr     string
	checkValuesExist bool
	avgLatencyMap    map[string]time.Duration
}

func newWeightedResponseStrategy() Strategy {
	return &WeightedResponseStrategy{}
}

// ReceiveData receive data
func (r *WeightedResponseStrategy) ReceiveData(instances []*registry.MicroServiceInstance, serviceName, protocol, sessionID string) {
	r.instances = instances
	serviceAndProtocol := serviceName + "/" + protocol
	r.avgLatencyMap, r.protocol = FindingAvgLatency(serviceAndProtocol)

	r.checkValuesExist = true
	// to check if 10 latency values of every instance is collected or not
	for k := range r.avgLatencyMap {
		LatencyMapRWMutex.RLock()
		lvalues := len(LatencyMap[k])
		LatencyMapRWMutex.RUnlock()
		if lvalues != 10 {
			r.checkValuesExist = false
		}
	}
	r.instanceAddr = SortingLatencyDuration(serviceAndProtocol, r.avgLatencyMap)
}

// Pick return instance
func (r *WeightedResponseStrategy) Pick() (*registry.MicroServiceInstance, error) {
	if (r.checkValuesExist && len(r.avgLatencyMap) == 0) || !r.checkValuesExist {
		if len(r.instances) == 0 {
			return nil, ErrNoneAvailableInstance
		}

		//if no instances are selected round robin will be done
		weightedRespMutex.Lock()
		node := r.instances[i%len(r.instances)]
		i++
		weightedRespMutex.Unlock()
		return node, nil
	}
	if len(r.instances) == 0 {
		return nil, ErrNoneAvailableInstance
	}

	for _, node := range r.instances {
		weightedRespMutex.Lock()
		if r.instanceAddr == node.EndpointsMap[r.protocol] {
			weightedRespMutex.Unlock()
			return node, nil
		}
		weightedRespMutex.Unlock()
	}

	//if no instances are selected round robin will be done
	weightedRespMutex.Lock()
	node := r.instances[i%len(r.instances)]
	i++
	weightedRespMutex.Unlock()
	return node, nil

}
