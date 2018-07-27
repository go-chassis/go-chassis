package pilot

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoundRobin(t *testing.T) {
	next := RoundRobin([]string{
		"a", "b",
	})
	v1, e := next()
	assert.NoError(t, e)
	assert.NotEmpty(t, v1)
	v2, e := next()
	assert.NoError(t, e)
	assert.NotEqual(t, v1, v2)

	next = RoundRobin(nil)
	v, e := next()
	assert.Error(t, e)
	assert.Empty(t, v)
}
