package loadbalance_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/registry"
	"time"
)

func TestSessionStickyStrategies(t *testing.T) {
	config.Init()
	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
	}

	for name, strategy := range map[string]loadbalance.Strategy{"sessionstickiness": loadbalance.SessionStickiness} {

		next := strategy(testData, nil)
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
func TestStickySessionStrategy(t *testing.T) {
	config.Init()

	testData := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "127.0.0.1:8080", "highway": "10.0.0.3:8080"},
		},
	}

	for name, strategy := range map[string]loadbalance.Strategy{"sessionstickiness": loadbalance.SessionStickiness} {
		loadbalance.SessionCache.Set("sticky1", "sdhgfa", time.Second*10)
		next := strategy(testData, "sticky1")

		for i := 0; i < 100; i++ {
			_, err := next()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
		next1 := strategy(testData, "sticky1")

		for i := 0; i < 100; i++ {
			_, err := next1()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
	}
	for name, strategy := range map[string]loadbalance.Strategy{"sessionstickiness": loadbalance.SessionStickiness} {
		LBstr := make(map[string]string)
		LBstr["name"] = "SessionStickiness"
		LBstr["sessionTimeoutInSeconds"] = "30"
		config.GlobalDefinition.Cse.Loadbalance.Strategy = LBstr
		next := strategy(testData, "sticky3")

		for i := 0; i < 100; i++ {
			_, err := next()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)

		next1 := strategy(testData, "sticky3")

		for i := 0; i < 100; i++ {
			_, err := next1()
			if err != nil {
				t.Fatal(err)
			}
		}

		t.Logf("%s", name)
	}
}
