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

	_, _, err = util.ParsePortName("httpadmin")
	assert.Error(t, err)
}
