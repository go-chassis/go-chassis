package loadbalance_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	_ "github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/registry"

	"github.com/stretchr/testify/assert"
	"time"
)

func TestLatencyStrategyNoInstanceMapZero(t *testing.T) {
	config.Init()
	testData := []*registry.MicroServiceInstance{}

	for _, strategy := range map[string]loadbalance.Strategy{"weightedresponse": loadbalance.WeightedResponse} {

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

	for name, strategy := range map[string]loadbalance.Strategy{"weightedresponse": loadbalance.WeightedResponse} {

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

	for _, strategy := range map[string]loadbalance.Strategy{"weightedresponse": loadbalance.WeightedResponse} {
		next := strategy(testData, "")

		_, err := next()
		assert.Error(t, err)
	}
}

func TestLatencyFunc(t *testing.T) {
	loadbalance.SetLatency(time.Second, "127.0.0.1:8080", "Server")
	var avgLatency = make(map[string]time.Duration)
	avgLatency["127.0.0.1:8080/Server"] = time.Second
	avgLatency["127.0.0.1:8081/Server"] = time.Second
	_ = loadbalance.SortingLatencyDuration("Server", avgLatency)
	_ = loadbalance.FindingAvgLatency("Server")
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

	for name, strategy := range map[string]loadbalance.Strategy{"weightedresponse": loadbalance.WeightedResponse} {

		next := strategy(testData, "Server")

		for i := 0; i < 100; i++ {
			_, err := next()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
		next1 := strategy(testData, "Server")

		for i := 0; i < 100; i++ {
			_, err := next1()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
	}
	for name, strategy := range map[string]loadbalance.Strategy{"weightedresponse": loadbalance.WeightedResponse} {
		LBstr := make(map[string]string)
		LBstr["name"] = "WeightedResponse"
		config.GlobalDefinition = &model.GlobalCfg{}
		config.GlobalDefinition.Cse.Loadbalance.Strategy = LBstr
		loadbalance.SetLatency(time.Second, "127.0.0.1:8080", "Serverrest")
		var tempDur []time.Duration
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		loadbalance.SetLatency(time.Second, "127.0.0.1:8080", "Serverrest")
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)
		tempDur = append(tempDur, time.Second)

		loadbalance.LatencyMap["127.0.0.1:8080/Serverrest"] = tempDur
		loadbalance.LatencyMap["127.0.0.1:8081/Serverrest"] = tempDur
		loadbalance.SetLatency(time.Second, "127.0.0.1:8080", "Serverrest")
		loadbalance.RestReqCount = 22
		next := strategy(testData, "Server")

		for i := 0; i < 100; i++ {
			_, err := next()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)

		next1 := strategy(testData, "Serverrest")

		for i := 0; i < 100; i++ {
			_, err := next1()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
	}
}
