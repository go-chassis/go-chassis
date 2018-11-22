package iputil_test

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLocapIp(t *testing.T) {
	localip := iputil.GetLocalIP()
	assert.NotNil(t, localip)
}

func TestLocalhost(t *testing.T) {
	assert.NotEmpty(t, iputil.Localhost())
}

func TestDefaultEndpoint4Protocol(t *testing.T) {
	assert.NotEmpty(t, iputil.DefaultEndpoint4Protocol(common.ProtocolRest))
}

func TestDefaultPort4ProtocolRest(t *testing.T) {
	assert.NotEmpty(t, iputil.DefaultPort4Protocol(common.ProtocolRest))
}

func TestDefaultPort4ProtocolHighway(t *testing.T) {
	assert.NotEmpty(t, iputil.DefaultPort4Protocol(common.ProtocolHighway))
}

func TestDefaultPort4ProtocolNone(t *testing.T) {
	assert.NotEmpty(t, iputil.DefaultPort4Protocol("http"))
}

//
//func TestLocalhost2(t *testing.T) {
//	_, err := net.Dial("tcp", "[fe80::7f28:7160:56cd:3ec9%enp0s3]:5001")
//	assert.NoError(t, err)
//	_,err:= url.Parse("http://[fe80::7f28:7160:56cd:3ec9]:5001")
//	//assert.Equal(t,"http://[fe80::7f28:7160:56cd:3ec9%25enp0s3]:5001",s)
//	_,err =http.DefaultClient.Get("http://[fe80::7f28:7160:56cd:3ec9]:5001")
//	assert.NoError(t,err)
//}
