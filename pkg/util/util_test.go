package util_test

import (
	"github.com/go-chassis/go-chassis/v2/pkg/util"
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
func TestParseServiceAndPort(t *testing.T) {
	s, p, err := util.ParseServiceAndPort("Service1:legacy")
	assert.Equal(t, "Service1", s)
	assert.Equal(t, "legacy", p)
	assert.NoError(t, err)

	s, p, err = util.ParseServiceAndPort("Service1")
	assert.Equal(t, "Service1", s)
	assert.Equal(t, "", p)
	assert.NoError(t, err)

	s, p, err = util.ParseServiceAndPort("http://Service1:admin")
	assert.Equal(t, util.ErrInvalidURL, err)

	s, p, err = util.ParseServiceAndPort("")
	assert.Equal(t, util.ErrInvalidURL, err)
}

func TestGenProtoEndPoint(t *testing.T) {
	ep := util.GenProtoEndPoint("rest", "admin")
	assert.Equal(t, "rest-admin", ep)

	ep = util.GenProtoEndPoint("rest", "")
	assert.Equal(t, "rest", ep)
}
