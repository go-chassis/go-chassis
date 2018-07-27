package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWrapInstance(t *testing.T) {
	wi := WrapInstance{
		AppID:       "1",
		ServiceName: "2",
		Version:     "3",
		Instance:    &MicroServiceInstance{InstanceID: "4"},
	}
	assert.Equal(t, "2:3:1:4", wi.String())
	assert.Equal(t, "2:3:1", wi.ServiceKey())
}
