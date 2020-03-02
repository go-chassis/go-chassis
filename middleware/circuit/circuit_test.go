package circuit_test

import (
	"errors"
	"testing"

	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/middleware/circuit"
	"github.com/stretchr/testify/assert"
)

func TestFallbackErr(t *testing.T) {
	inv := &invocation.Invocation{}
	finish := make(chan *invocation.Response)
	f := circuit.FallbackErr(inv, finish)

	err := f(errors.New("internal error"))
	assert.NoError(t, err)
}
