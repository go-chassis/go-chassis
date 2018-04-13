package pilot

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var (
	s  = httptest.NewServer(&mockPilotHandler{})
	r  = newRegistrator(registry.Options{Addrs: []string{s.Listener.Addr().String()}})
	sd = newDiscoveryService(registry.Options{Addrs: []string{s.Listener.Addr().String()}})
)

func init() {
	lager.Initialize("stdout", "", "", "",
		true, 0, 0, 0)
}

func TestPilot_RegisterServiceAndInstance(t *testing.T) {
	microservice := &registry.MicroService{
		ServiceName: "a",
	}
	microServiceInstance := &registry.MicroServiceInstance{
		EndpointsMap: map[string]string{"rest": "1.1.1.1:80"},
	}
	serviceID, instanceID, err := r.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)
	assert.Equal(t, "a", serviceID)
	assert.Equal(t, "1.1.1.1_80", instanceID)

	microservice, err = sd.GetMicroService("a")
	assert.NoError(t, err)
	assert.Equal(t, "a", microservice.ServiceName)

	serviceID, err = sd.GetMicroServiceID("", "a", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "a", serviceID)

	services, err := sd.GetAllMicroServices()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(services))

	instances, err := sd.GetMicroServiceInstances("", "a")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(instances))
	assert.Equal(t, instanceID, instances[0].InstanceID)
	assert.Equal(t, microServiceInstance.EndpointsMap["rest"], instances[0].EndpointsMap["rest"])

	instances, err = sd.FindMicroServiceInstances("", "", "a", "", "")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(instances))
	assert.Equal(t, instanceID, instances[0].InstanceID)
	assert.Equal(t, microServiceInstance.EndpointsMap["rest"], instances[0].EndpointsMap["rest"])
}
