package loadbalancer_test

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
)

func TestRandomStrategy_Pick(t *testing.T) {
	config.Init()
	instances := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]*registry.Endpoint{
				"rest": {
					false,
					"1",
				},
				"highway": {
					false,
					"10.0.0.3:8080",
				},
			},
		},
		{
			EndpointsMap: map[string]*registry.Endpoint{
				"rest": {
					false,
					"2",
				},
				"highway": {
					false,
					"10.0.0.3:8080",
				},
			},
		},
	}
	s := &loadbalancer.RandomStrategy{}
	s.ReceiveData(nil, nil, "")
	_, err := s.Pick()
	assert.Error(t, err)
	s.ReceiveData(nil, instances, "")
	var last = "none"
	var count int
	for i := 0; i < 100; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		if last == instance.EndpointsMap["rest"].GenEndpoint() {
			count++
		}
		last = instance.EndpointsMap["rest"].GenEndpoint()
	}
	t.Log(count)

}

func TestRoundRobinStrategy_Pick(t *testing.T) {
	config.Init()
	instances := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]*registry.Endpoint{
				"rest": {
					false,
					"1",
				},
				"highway": {
					false,
					"10.0.0.3:8080",
				},
			},
		},
		{
			EndpointsMap: map[string]*registry.Endpoint{
				"rest": {
					false,
					"2",
				},
				"highway": {
					false,
					"10.0.0.3:8080",
				},
			},
		},
	}
	s := &loadbalancer.RoundRobinStrategy{}
	s.ReceiveData(nil, instances, "")
	var last = "none"
	for i := 0; i < 100000; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		assert.NotEqual(t, last, instance.EndpointsMap["rest"].GenEndpoint())
		last = instance.EndpointsMap["rest"].GenEndpoint()
	}

}
func new() loadbalancer.Strategy {
	return nil
}
func TestGetStrategyPlugin(t *testing.T) {
	_, err := loadbalancer.GetStrategyPlugin("test")
	assert.Error(t, err)
	loadbalancer.InstallStrategy(loadbalancer.StrategyRoundRobin, new)
	_, err = loadbalancer.GetStrategyPlugin(loadbalancer.StrategyRoundRobin)
	assert.NoError(t, err)
}
