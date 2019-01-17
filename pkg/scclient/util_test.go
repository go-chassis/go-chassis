package client_test

import (
	"github.com/go-chassis/go-chassis/pkg/scclient"
	"github.com/go-chassis/go-chassis/pkg/scclient/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegroupInstances(t *testing.T) {
	keys := []*proto.FindService{
		{
			Service: &proto.MicroServiceKey{
				ServiceName: "Service1",
			},
		},
		{
			Service: &proto.MicroServiceKey{
				ServiceName: "Service2",
			},
		},
		{
			Service: &proto.MicroServiceKey{
				ServiceName: "Service3",
			},
		},
	}
	resp := proto.BatchFindInstancesResponse{
		Services: &proto.BatchFindResult{
			Updated: []*proto.FindResult{
				{Index: 2,
					Instances: []*proto.MicroServiceInstance{{
						InstanceId: "1",
					}}},
			},
		},
	}
	m := client.RegroupInstances(keys, resp)
	t.Log(m)
	assert.Equal(t, 1, len(m["Service3"]))
	assert.Equal(t, 0, len(m["Service1"]))
	assert.Equal(t, 0, len(m["Service2"]))
}
