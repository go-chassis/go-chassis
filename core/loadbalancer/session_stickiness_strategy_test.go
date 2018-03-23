package loadbalancer_test

import (
	"testing"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/session"
	"github.com/stretchr/testify/assert"
)

func TestGetSuccessiveFailureCount(t *testing.T) {
	c := loadbalancer.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 0, c)
	loadbalancer.IncreaseSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	c = loadbalancer.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 1, c)
	loadbalancer.IncreaseSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	c = loadbalancer.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 2, c)
	loadbalancer.DeleteSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	c = loadbalancer.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 0, c)
	loadbalancer.ResetSuccessiveFailureMap()
	c = loadbalancer.GetSuccessiveFailureCount("0807040b-0f08-4609-4608-010c00050e03")
	assert.Equal(t, 0, c)
}
func TestSessionStickinessStrategy_Pick(t *testing.T) {
	config.Init()
	instances := []*registry.MicroServiceInstance{
		{
			EndpointsMap: map[string]string{"rest": "1", "highway": "10.0.0.3:8080"},
		},
		{
			EndpointsMap: map[string]string{"rest": "2", "highway": "10.0.0.3:8080"},
		},
	}

	s := &loadbalancer.SessionStickinessStrategy{}
	s.ReceiveData(instances, "", "", "dummy")
	var last string = "none"
	for i := 0; i < 100; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		assert.NotEqual(t, last, instance.EndpointsMap["rest"])
		last = instance.EndpointsMap["rest"]
	}
	session.Save("dummy", "1", 0)
	for i := 0; i < 100; i++ {
		instance, err := s.Pick()
		assert.NoError(t, err)
		assert.NotEqual(t, 1, instance.EndpointsMap["rest"])
	}
}
