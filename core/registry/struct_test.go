package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMicroServiceInstance_Equal(t *testing.T) {
	ins1 := &MicroServiceInstance{
		InstanceID: "1",
		ServiceID:  "bill",
	}
	ins2 := &MicroServiceInstance{
		InstanceID: "1",
		ServiceID:  "bill",
	}
	assert.True(t, ins1.Equal(ins2))

	ins3 := &MicroServiceInstance{
		InstanceID: "1",
		ServiceID:  "bill",
		Metadata: map[string]string{
			"a": "b",
			"c": "d",
		},
	}
	ins4 := &MicroServiceInstance{
		InstanceID: "1",
		ServiceID:  "bill",
		Metadata: map[string]string{
			"a": "b",
			"c": "d",
		},
	}
	assert.True(t, ins3.Equal(ins4))

	ins5 := &MicroServiceInstance{
		InstanceID: "2",
		ServiceID:  "bill",
	}
	assert.False(t, ins5.Equal(ins4))

	ins6 := &MicroServiceInstance{
		InstanceID: "2",
		ServiceID:  "text",
	}
	assert.False(t, ins5.Equal(ins6))
}

