package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildNodeID(t *testing.T) {
	cases := []struct {
		NS   string
		Name string
		IP   string
	}{
		{"", "istio", "10.120.10.22"},
		{"pilot", "bob", "10.230.12.22"},
		{"default", "tina", "127.0.0.1"},
	}

	result := []string{
		"sidecar~10.120.10.22~istio.default~default.svc.cluster.local",
		"sidecar~10.230.12.22~bob.pilot~pilot.svc.cluster.local",
		"sidecar~127.0.0.1~tina.default~default.svc.cluster.local",
	}

	for i, n := range cases {
		os.Setenv(PODNAMESPACE, n.NS)
		os.Setenv(PODNAME, n.Name)
		os.Setenv(PODIP, n.IP)

		nodeID := BuildNodeID()
		assert.Equal(t, nodeID, result[i])
	}
}

func TestServiceAndLabel(t *testing.T) {
	cases := []string{"istioserver", "istioclient"}
	result := []string{
		"istioserver.default.svc.cluster.local",
		"istioclient.default.svc.cluster.local",
	}
	for i, n := range cases {
		assert.Equal(t, ServiceKey(n), result[i])
	}

	cases = []string{
		"outbound|9080|v1|ratings.default.svc.cluster.local",
		"outbound|9080|v2|reviews.default.svc.cluster.local",
	}
	result = []string{"v1", "v2"}
	for i, n := range cases {
		assert.Equal(t, ServiceKeyToLabel(n), result[i])
	}

	cases = []string{
		"ratings.default.svc.cluster.local:8090",
		"reviews.default.svc.cluster.local:15109",
	}
	service := []string{"ratings", "reviews"}
	port := []string{"8090", "15109"}
	for i, n := range cases {
		ss, sp := ServiceAndPort(n)
		assert.Equal(t, ss, service[i])
		assert.Equal(t, sp, port[i])
	}
}
