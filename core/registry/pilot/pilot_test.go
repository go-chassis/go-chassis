package pilot

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var (
	s = httptest.NewServer(&mockPilotHandler{})
	r = newPilotRegistry(registry.Addrs(s.Listener.Addr().String()))
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
	serviceId, instanceId, err := r.RegisterServiceAndInstance(microservice, microServiceInstance)
	assert.NoError(t, err)
	assert.Equal(t, "a", serviceId)
	assert.Equal(t, "1.1.1.1_80", instanceId)

	microservice, err = r.GetMicroService("a")
	assert.NoError(t, err)
	assert.Equal(t, "a", microservice.ServiceName)

	serviceId, err = r.GetMicroServiceID("", "a", "", "")
	assert.NoError(t, err)
	assert.Equal(t, "a", serviceId)

	services, err := r.GetAllMicroServices()
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(services))

	instances, err := r.GetMicroServiceInstances("", "a")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(instances))
	assert.Equal(t, instanceId, instances[0].InstanceID)
	assert.Equal(t, microServiceInstance.EndpointsMap["rest"], instances[0].EndpointsMap["rest"])

	instances, err = r.FindMicroServiceInstances("", "", "a", "", "")
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(instances))
	assert.Equal(t, instanceId, instances[0].InstanceID)
	assert.Equal(t, microServiceInstance.EndpointsMap["rest"], instances[0].EndpointsMap["rest"])
}
