package circuit_test

import (
	"errors"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/circuit"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFallbackErr(t *testing.T) {
	inv := &invocation.Invocation{}
	finish := make(chan *invocation.Response)
	f := circuit.FallbackErr(inv, finish)

	err := f(errors.New("internal error"))
	assert.NoError(t, err)
}
