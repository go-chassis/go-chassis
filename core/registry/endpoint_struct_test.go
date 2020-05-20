package registry_test

import (
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEndPoint(t *testing.T) {
	sslEndpoint := "10.0.1.9:1234?" + registry.SSLEnabledQuery
	ep, err := registry.NewEndPoint(sslEndpoint)
	assert.Nil(t, err, "parse "+sslEndpoint+" error")
	assert.True(t, ep.GenEndpoint() == sslEndpoint)
	assert.True(t, ep.IsSSLEnable())

	commonEndpoint := "10.0.1.9:8080"
	ep, err = registry.NewEndPoint(commonEndpoint)
	assert.Nil(t, err, "parse "+commonEndpoint+" error")
	assert.True(t, ep.GenEndpoint() == commonEndpoint)
	assert.False(t, ep.IsSSLEnable())

	noPortSslEndpoint := "10.0.1.9?" + registry.SSLEnabledQuery
	ep, err = registry.NewEndPoint(noPortSslEndpoint)
	assert.Nil(t, err, "parse "+noPortSslEndpoint+" error")
	assert.True(t, ep.IsSSLEnable())

	sslFalseEndpoint := "10.0.1.9?sslEnabled=false"
	ep, err = registry.NewEndPoint(sslFalseEndpoint)
	assert.Nil(t, err, "parse "+sslFalseEndpoint+" error")
	assert.False(t, ep.IsSSLEnable())
}
