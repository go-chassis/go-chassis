package client_test

import (
	"github.com/go-chassis/go-chassis/v2/pkg/scclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWithGlobal(t *testing.T) {
	o := client.WithGlobal()
	opts := &client.CallOptions{}
	o(opts)
	assert.True(t, opts.WithGlobal)
}
