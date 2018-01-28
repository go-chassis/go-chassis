package loadbalance_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	_ "github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/registry"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestLatencyStrategyNoInstanceMapZero(t *testing.T) {
	config.Init()
	testData := []*registry.MicroServiceInstance{}

	for _, strategy := range map[string]selector.Strategy{"weightedresponse": loadbalance.WeightedResponse} {

		next := strategy(testData, "")

		_, err := next()
		assert.Error(t, err)
	}
}

func TestLatencyStrategies(t *testing.T) {
	config.Init()
	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
	}

	for name, strategy := range map[string]selector.Strategy{"weightedresponse": loadbalance.WeightedResponse} {

		next := strategy(testData, "")
		counts := make(map[string]int)

		for i := 0; i < 100; i++ {
			node, err := next()
			if err != nil {
				t.Fatal(err)
			}
			counts[node.InstanceID]++
		}

		t.Logf("%s: %+v", name, counts)
	}
}

func TestLatencyStrategyNoInstance(t *testing.T) {
	config.Init()
	testData := []*registry.MicroServiceInstance{}

	for _, strategy := range map[string]selector.Strategy{"weightedresponse": loadbalance.WeightedResponse} {
		next := strategy(testData, "")

		_, err := next()
		assert.Error(t, err)
	}
}

func TestLatencyFunc(t *testing.T) {
	loadbalance.SetLatency(time.Second, "127.0.0.1:8080", "Server/"+common.ProtocolHighway)
	var avgLatency = make(map[string]time.Duration)
	avgLatency["127.0.0.1:8080/Server/"+common.ProtocolHighway] = time.Second
	avgLatency["127.0.0.1:8081/Server/"+common.ProtocolHighway] = time.Second
	_ = loadbalance.SortingLatencyDuration("Server/"+common.ProtocolHighway, avgLatency)
	_, _ = loadbalance.FindingAvgLatency("Server/" + common.ProtocolHighway)
}
func TestLatencyStrategy(t *testing.T) {
	config.Init()

	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8081", "highway": "10.0.0.3:8080"},
		},
	}

	for name, strategy := range map[string]selector.Strategy{"weightedresponse": loadbalance.WeightedResponse} {

		next := strategy(testData, "Server/"+common.ProtocolHighway)

		for i := 0; i < 100; i++ {
			_, err := next()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
		next1 := strategy(testData, "Server/"+common.ProtocolHighway)

		for i := 0; i < 100; i++ {
			_, err := next1()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
	}
	for name, strategy := range map[string]selector.Strategy{"weightedresponse": loadbalance.WeightedResponse} {
		LBstr := make(map[string]string)
		LBstr["name"] = "WeightedResponse"
		config.GlobalDefinition = &model.GlobalCfg{}
		config.GlobalDefinition.Cse.Loadbalance.Strategy = LBstr
		loadbalance.SetLatency(time.Second, "127.0.0.1:8080", "Server/"+common.ProtocolRest)
		var tempDur []time.Duration
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		loadbalance.SetLatency(time.Second, "127.0.0.1:8080", "Server/"+common.ProtocolRest)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)

		loadbalance.LatencyMap["127.0.0.1:8080/Server/"+common.ProtocolRest] = tempDur
		loadbalance.LatencyMap["127.0.0.1:8081/Server/"+common.ProtocolRest] = tempDur
		loadbalance.SetLatency(time.Second, "127.0.0.1:8080", "Server/"+common.ProtocolRest)
		next := strategy(testData, "Server/"+common.ProtocolRest)

		for i := 0; i < 100; i++ {
			_, err := next()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)

		next1 := strategy(testData, "Server/"+common.ProtocolRest)

		for i := 0; i < 100; i++ {
			_, err := next1()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
	}
}
func TestFindingAvgLatency(t *testing.T) {
	loadbalance.LatencyMap = make(map[string][]time.Duration)
	loadbalance.SetLatency(1*time.Second, "127.0.0.1:3000", "Server/"+common.ProtocolRest)
	loadbalance.SetLatency(3*time.Second, "10.1.1.1.1:3000", "Server/"+common.ProtocolRest)
	loadbalance.SetLatency(1*time.Second, "127.0.0.1:5000", "Server/"+common.ProtocolHighway)
	loadbalance.SetLatency(5*time.Second, "127.0.0.1:5000", "Server/"+common.ProtocolHighway)
	loadbalance.SetLatency(1*time.Second, "10.0.0.1:5000", "Server/"+common.ProtocolHighway)
	loadbalance.SetLatency(9*time.Second, "10.0.0.1:5000", "Server/"+common.ProtocolHighway)
	avgLatencyMap, p := loadbalance.FindingAvgLatency("Server/" + common.ProtocolHighway)
	assert.Equal(t, common.ProtocolHighway, p)
	for k, v := range avgLatencyMap {
		if k == "127.0.0.1/Server/highway" {
			assert.Equal(t, time.Duration(3*time.Second), v)
		}
		if k == "10.0.0.1/Server/highway" {
			assert.Equal(t, time.Duration(5*time.Second), v)
		}
	}
	addr := loadbalance.SortingLatencyDuration("Server/"+common.ProtocolHighway, avgLatencyMap)
	assert.Equal(t, "127.0.0.1:5000", addr)
	t.Log(addr)

	avgLatencyMap, p = loadbalance.FindingAvgLatency("Server/" + common.ProtocolRest)
	assert.Equal(t, common.ProtocolRest, p)

	addr = loadbalance.SortingLatencyDuration("Server/"+common.ProtocolRest, avgLatencyMap)
	assert.Equal(t, "127.0.0.1:3000", addr)
	t.Log(addr)
}
