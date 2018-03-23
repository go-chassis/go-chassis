package loadbalancer_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	_ "github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"
	"github.com/ServiceComb/go-chassis/core/registry"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestLatencyFunc(t *testing.T) {
	loadbalancer.SetLatency(time.Second, "127.0.0.1:8080", "Server/"+common.ProtocolHighway)
	var avgLatency = make(map[string]time.Duration)
	avgLatency["127.0.0.1:8080/Server/"+common.ProtocolHighway] = time.Second
	avgLatency["127.0.0.1:8081/Server/"+common.ProtocolHighway] = time.Second
	_ = loadbalancer.SortingLatencyDuration("Server/"+common.ProtocolHighway, avgLatency)
	_, _ = loadbalancer.FindingAvgLatency("Server/" + common.ProtocolHighway)
}
func TestLatencyStrategy(t *testing.T) {
	config.Init()

	instances := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8081", "highway": "10.0.0.3:8080"},
		},
	}

	LBstr := make(map[string]string)
	LBstr["name"] = "WeightedResponse"
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GetLoadBalancing().Strategy = LBstr
	loadbalancer.SetLatency(time.Second, "127.0.0.1:8080", "Server/"+common.ProtocolRest)
	var tempDur []time.Duration
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	loadbalancer.SetLatency(time.Second, "127.0.0.1:8080", "Server/"+common.ProtocolRest)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)
	tempDur = append(tempDur, time.Second)

	loadbalancer.LatencyMap["127.0.0.1:8080/Server/"+common.ProtocolRest] = tempDur
	loadbalancer.LatencyMap["127.0.0.1:8081/Server/"+common.ProtocolRest] = tempDur
	loadbalancer.SetLatency(time.Second, "127.0.0.1:8080", "Server/"+common.ProtocolRest)
	s := &loadbalancer.WeightedResponseStrategy{}
	s.ReceiveData(instances, "Server", common.ProtocolRest, "")
	var last string = "none"
	for i := 0; i < 100; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		assert.NotEqual(t, last, instance.EndpointsMap["rest"])
		last = instance.EndpointsMap["rest"]
	}
	for i := 0; i < 100; i++ {
		_, err := s.Pick()
		if err != nil {
			t.Fatal(err)
		}
	}
}
func TestFindingAvgLatency(t *testing.T) {
	loadbalancer.LatencyMap = make(map[string][]time.Duration)
	loadbalancer.SetLatency(1*time.Second, "127.0.0.1:3000", "Server/"+common.ProtocolRest)
	loadbalancer.SetLatency(3*time.Second, "10.1.1.1.1:3000", "Server/"+common.ProtocolRest)
	loadbalancer.SetLatency(1*time.Second, "127.0.0.1:5000", "Server/"+common.ProtocolHighway)
	loadbalancer.SetLatency(5*time.Second, "127.0.0.1:5000", "Server/"+common.ProtocolHighway)
	loadbalancer.SetLatency(1*time.Second, "10.0.0.1:5000", "Server/"+common.ProtocolHighway)
	loadbalancer.SetLatency(9*time.Second, "10.0.0.1:5000", "Server/"+common.ProtocolHighway)
	avgLatencyMap, p := loadbalancer.FindingAvgLatency("Server/" + common.ProtocolHighway)
	assert.Equal(t, common.ProtocolHighway, p)
	for k, v := range avgLatencyMap {
		if k == "127.0.0.1/Server/highway" {
			assert.Equal(t, time.Duration(3*time.Second), v)
		}
		if k == "10.0.0.1/Server/highway" {
			assert.Equal(t, time.Duration(5*time.Second), v)
		}
	}
	addr := loadbalancer.SortingLatencyDuration("Server/"+common.ProtocolHighway, avgLatencyMap)
	assert.Equal(t, "127.0.0.1:5000", addr)
	t.Log(addr)

	avgLatencyMap, p = loadbalancer.FindingAvgLatency("Server/" + common.ProtocolRest)
	assert.Equal(t, common.ProtocolRest, p)

	addr = loadbalancer.SortingLatencyDuration("Server/"+common.ProtocolRest, avgLatencyMap)
	assert.Equal(t, "127.0.0.1:3000", addr)
	t.Log(addr)
}
