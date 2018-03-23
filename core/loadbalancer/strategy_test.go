package loadbalancer_test

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
)

func TestRandomStrategy_Pick(t *testing.T) {
	config.Init()
	instances := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "1", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "2", "highway": "10.0.0.3:8080"},
		},
	}
	s := &loadbalancer.RandomStrategy{}
	s.ReceiveData(instances, "", "", "")
	var last string = "none"
	var count int
	for i := 0; i < 100; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		if last == instance.EndpointsMap["rest"] {
			count++
		}
		last = instance.EndpointsMap["rest"]
	}
	t.Log(count)

}

func TestRoundRobinStrategy_Pick(t *testing.T) {
	config.Init()
	instances := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "1", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "2", "highway": "10.0.0.3:8080"},
		},
	}
	s := &loadbalancer.RoundRobinStrategy{}
	s.ReceiveData(instances, "", "", "")
	var last string = "none"
	for i := 0; i < 100000; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		assert.NotEqual(t, last, instance.EndpointsMap["rest"])
		last = instance.EndpointsMap["rest"]
	}

}
