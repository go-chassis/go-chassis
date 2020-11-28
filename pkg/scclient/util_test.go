package client_test

import (
	scregistry "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/go-chassis/v2/core/registry/servicecenter"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegroupInstances(t *testing.T) {
	keys := []*scregistry.FindService{
		{
			Service: &scregistry.MicroServiceKey{
				ServiceName: "Service1",
			},
		},
		{
			Service: &scregistry.MicroServiceKey{
				ServiceName: "Service2",
			},
		},
		{
			Service: &scregistry.MicroServiceKey{
				ServiceName: "Service3",
			},
		},
	}
	resp := &scregistry.BatchFindInstancesResponse{
		Services: &scregistry.BatchFindResult{
			Updated: []*scregistry.FindResult{
				{Index: 2,
					Instances: []*scregistry.MicroServiceInstance{{
						InstanceId: "1",
					}}},
			},
		},
	}
	m := servicecenter.RegroupInstances(keys, resp)
	t.Log(m)
	assert.Equal(t, 1, len(m["Service3"]))
	assert.Equal(t, 0, len(m["Service1"]))
	assert.Equal(t, 0, len(m["Service2"]))
}
