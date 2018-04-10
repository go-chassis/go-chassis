package pilot

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPilotHandler struct {
}

func (m *mockPilotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var b []byte
	hs1 := Hosts{
		Hosts: []*Host{
			{
				"1.1.1.1", 80,
				&Tags{
					"az", false, 80,
				},
			},
			{
				"1.1.1.2", 80,
				&Tags{
					"az", false, 20,
				},
			},
		},
	}
	hs2 := Hosts{
		Hosts: []*Host{
			{
				"1.1.2.1", 80,
				&Tags{
					"az", false, 100,
				},
			},
		},
	}
	if r.URL.Path == BaseRoot {
		s := []*Service{
			{ServiceKey: "a", Hosts: hs1.Hosts},
			{ServiceKey: "b", Hosts: hs2.Hosts},
		}
		b, _ = json.Marshal(s)
	} else {
		b, _ = json.Marshal(hs1)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

type mockErrPilotHandler struct {
}

func (m *mockErrPilotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestPilotClient_GetServiceHosts(t *testing.T) {
	client := EnvoyDSClient{}
	err := client.Initialize(Options{})
	assert.NoError(t, err)
	hosts, err := client.GetServiceHosts("a")
	assert.Nil(t, hosts)
	assert.Error(t, err)
	defer client.Close()

	s := httptest.NewServer(http.NotFoundHandler())
	client.Options.Addrs = []string{s.Listener.Addr().String()}
	hosts, err = client.GetServiceHosts("a")
	assert.Nil(t, hosts)
	assert.Error(t, err)
	s.Close()

	s = httptest.NewServer(&mockErrPilotHandler{})
	client.Options.Addrs = []string{s.Listener.Addr().String()}
	hosts, err = client.GetServiceHosts("a")
	assert.Nil(t, hosts)
	assert.Error(t, err)
	s.Close()

	s = httptest.NewServer(&mockPilotHandler{})
	client.Options.Addrs = []string{s.Listener.Addr().String()}
	hosts, err = client.GetServiceHosts("a")
	assert.NotNil(t, hosts)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(hosts.Hosts))
	assert.Equal(t, "1.1.1.1", hosts.Hosts[0].Address)
	assert.Equal(t, 80, hosts.Hosts[0].Port)
	s.Close()
}

func TestPilotClient_GetAllServices(t *testing.T) {
	client := EnvoyDSClient{}
	err := client.Initialize(Options{})
	assert.NoError(t, err)
	svcs, err := client.GetAllServices()
	assert.Nil(t, svcs)
	assert.Error(t, err)
	defer client.Close()

	s := httptest.NewServer(http.NotFoundHandler())
	client.Options.Addrs = []string{s.Listener.Addr().String()}
	svcs, err = client.GetAllServices()
	assert.Nil(t, svcs)
	assert.Error(t, err)
	s.Close()

	s = httptest.NewServer(&mockErrPilotHandler{})
	client.Options.Addrs = []string{s.Listener.Addr().String()}
	svcs, err = client.GetAllServices()
	assert.Nil(t, svcs)
	assert.Error(t, err)
	s.Close()

	s = httptest.NewServer(&mockPilotHandler{})
	client.Options.Addrs = []string{s.Listener.Addr().String()}
	svcs, err = client.GetAllServices()
	assert.NotNil(t, svcs)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(svcs))
	assert.Equal(t, "a", svcs[0].ServiceKey)
	s.Close()
}
