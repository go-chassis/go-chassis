package util_test

import (
	"github.com/ServiceComb/go-chassis/pkg/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsePortName(t *testing.T) {
	p, n, err := util.ParsePortName("http-admin")
	assert.Equal(t, "http", p)
	assert.Equal(t, "admin", n)
	assert.NoError(t, err)

	_, _, err = util.ParsePortName("http")
	assert.NoError(t, err)

	_, _, err = util.ParsePortName("")
	assert.Error(t, err)

	_, _, err = util.ParsePortName("http-admin-1")
	assert.Error(t, err)
}
