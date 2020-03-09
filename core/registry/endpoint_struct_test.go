package registry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndPoint(t *testing.T) {
	sslEndpoint := "10.0.1.9:1234?" + ssLEnabledQuery
	ep, err := NewEndPoint(sslEndpoint)
	assert.Nil(t, err, "parse "+sslEndpoint+" error")
	assert.True(t, ep.GenEndpoint() == sslEndpoint)
	assert.True(t, ep.IsSSLEnable())

	commonEndpoint := "10.0.1.9:8080"
	ep, err = NewEndPoint(commonEndpoint)
	assert.Nil(t, err, "parse "+commonEndpoint+" error")
	assert.True(t, ep.GenEndpoint() == commonEndpoint)
	assert.False(t, ep.IsSSLEnable())

	noPortSslEndpoint := "10.0.1.9?" + ssLEnabledQuery
	ep, err = NewEndPoint(noPortSslEndpoint)
	assert.Nil(t, err, "parse "+noPortSslEndpoint+" error")
	assert.True(t, ep.IsSSLEnable())

	sslFalseEndpoint := "10.0.1.9?sslEnabled=false"
	ep, err = NewEndPoint(sslFalseEndpoint)
	assert.Nil(t, err, "parse "+sslFalseEndpoint+" error")
	assert.False(t, ep.IsSSLEnable())
}
