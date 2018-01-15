package transport_test

import (
	"errors"
	"net"
	"testing"

	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"github.com/stretchr/testify/assert"
)

func TestListenTransport(t *testing.T) {

	fn := func(string) (net.Listener, error) {
		var n net.Listener

		return n, nil
	}
	n, err := transport.Listen("0.0.0.0", fn)
	assert.Nil(t, n)
	assert.Nil(t, err)

	n, err = transport.Listen("0.0.0.0:8080", fn)
	assert.Nil(t, n)
	assert.Nil(t, err)

	n, err = transport.Listen("0.0.0.0:8080-8000", fn)
	assert.Nil(t, n)
	assert.Error(t, err)

	n, err = transport.Listen("0.0.0.0:8080-8090", fn)
	assert.Nil(t, n)
	assert.Nil(t, err)

	n, err = transport.Listen("0.0.0.0:abcd-1234", fn)
	assert.Nil(t, n)
	assert.Error(t, err)

	n, err = transport.Listen("0.0.0.0:8080-abcd", fn)
	assert.Nil(t, n)
	assert.Error(t, err)

	fn = func(string) (net.Listener, error) {
		return nil, errors.New("Invalid")
	}
	n, err = transport.Listen("0.0.0.0:8080-8090", fn)
	assert.Nil(t, n)
	assert.Error(t, err)
}
