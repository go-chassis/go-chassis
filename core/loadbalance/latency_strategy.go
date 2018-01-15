package loadbalance

import (
	"github.com/ServiceComb/go-chassis/core/registry"
	"sort"
	"strings"
	"sync"
	"time"
)

// ByDuration is for calculating the duration
type ByDuration []time.Duration

func (a ByDuration) Len() int           { return len(a) }
func (a ByDuration) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDuration) Less(i, j int) bool { return a[i] < a[j] }

// variables for latency map, rest and highway requests count
var (
	//LatencyMap key is the combination of instance addr and microservice name separated by "/"
	LatencyMap      map[string][]time.Duration
	RestReqCount    = 0
	HighwayReqCount = 0
	//maintain different locks since multiple goroutine access the map
	mutexLatencyMap   sync.RWMutex
	avgmtx            sync.RWMutex
	weightedRespMutex sync.Mutex
)

// SetLatency for each requests
func SetLatency(duration time.Duration, addr, microServiceName string) {
	if strings.Contains(microServiceName, "rest") {
		RestReqCount++
	} else {
		HighwayReqCount++
	}
	//ReqCount++
	key := addr + "/" + microServiceName

	mutexLatencyMap.RLock()
	val, ok := LatencyMap[key]
	mutexLatencyMap.RUnlock()

	if !ok {
		var durationQueue []time.Duration
		durationQueue = append(durationQueue, duration)
		mutexLatencyMap.Lock()
		LatencyMap[key] = durationQueue
		mutexLatencyMap.Unlock()
	} else {
		mutexLatencyMap.Lock()
		if len(val) < 10 {
			val = append(val, duration)
			LatencyMap[key] = val
		} else {
			val = val[1:]
			val = append(val, duration)
			LatencyMap[key] = val
		}
		mutexLatencyMap.Unlock()
	}
}

// WeightedResponse is a strategy
func WeightedResponse(instances []*registry.MicroServiceInstance, metadata interface{}) Next {
	if strings.Contains(metadata.(string), "rest") {
		return selectWeightedInstance(RestReqCount, instances, metadata)
	}

	return selectWeightedInstance(HighwayReqCount, instances, metadata)

}

// SortingLatencyDuration sorting the average latencies recored for each instance
// and returning the instance addr which has the least latency
func SortingLatencyDuration(metadata string, avgLatencyMap map[string]time.Duration) string {
	var mtx sync.Mutex
	var tempLatencyMap = make(map[string]time.Duration)
	for k, v := range avgLatencyMap {
		epMs := strings.Split(k, "/")
		//comparing the microservice name
		if epMs[1] == metadata {
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

// FindingAvgLatency Calculating the average latency for each instance using the statistics collected
func FindingAvgLatency(metadata string) map[string]time.Duration {
	var avgLatencyMap = make(map[string]time.Duration)
	mutexLatencyMap.RLock()
	defer mutexLatencyMap.RUnlock()
	for k, v := range LatencyMap {
		epMs := strings.Split(k, "/")
		//comparing the microservice name
		if epMs[1] == metadata {
			var t2, tempDuration time.Duration
			for i := 0; i < len(v); i++ {
				t2 = t2 + v[i]
			}

			tempDuration = t2 / time.Duration(int64(len(v)))
			avgmtx.Lock()
			avgLatencyMap[k] = tempDuration
			avgmtx.Unlock()
		}
	}

	return avgLatencyMap

}

// selectWeightedInstance select instance based on protocol and less latency
func selectWeightedInstance(reqCount int, instances []*registry.MicroServiceInstance, metadata interface{}) Next {
	//Checking for length of instances * 10 to get 10 sample latencies for each instance to
	//get average latency for each instance.
	if reqCount <= len(instances)*10 {
		return func() (*registry.MicroServiceInstance, error) {
			if len(instances) == 0 {
				return nil, ErrNoneAvailable
			}
			weightedRespMutex.Lock()
			node := instances[i%len(instances)]
			i++
			weightedRespMutex.Unlock()
			return node, nil
		}
	}
	var instanceAddr string
	if len(LatencyMap) != 0 {
		avgLatencyMap := FindingAvgLatency(metadata.(string))
		instanceAddr = SortingLatencyDuration(metadata.(string), avgLatencyMap)
	}

	return func() (*registry.MicroServiceInstance, error) {
		if len(instances) == 0 {
			return nil, ErrNoneAvailable
		}

		for _, node := range instances {
			weightedRespMutex.Lock()
			if instanceAddr == node.EndpointsMap["rest"] {
				weightedRespMutex.Unlock()
				return node, nil
			} else if instanceAddr == node.EndpointsMap["highway"] {
				weightedRespMutex.Unlock()
				return node, nil
			}
			weightedRespMutex.Unlock()
		}

		//if no instances are selected round robin will be done
		weightedRespMutex.Lock()
		node := instances[i%len(instances)]
		i++
		weightedRespMutex.Unlock()
		return node, nil
	}
}
