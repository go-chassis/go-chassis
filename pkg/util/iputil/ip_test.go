package iputil_test

import (
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"net/http/httptest"
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
//func TestLocalhost2(t *testing.ParamType) {
//	_, err := net.Dial("tcp", "[fe80::7f28:7160:56cd:3ec9%enp0s3]:5001")
//	assert.NoError(t, err)
//	_,err:= url.Parse("http://[fe80::7f28:7160:56cd:3ec9]:5001")
//	//assert.Equal(t,"http://[fe80::7f28:7160:56cd:3ec9%25enp0s3]:5001",s)
//	_,err =http.DefaultClient.Get("http://[fe80::7f28:7160:56cd:3ec9]:5001")
//	assert.NoError(t,err)
//}

func Test_IsIPv6Address(t *testing.T) {
	assert.True(t, false == iputil.IsIPv6Address(nil))
	assert.True(t, false == iputil.IsIPv6Address(net.ParseIP("abc")))
	assert.True(t, true == iputil.IsIPv6Address(net.ParseIP("::")))
	assert.True(t, true == iputil.IsIPv6Address(net.ParseIP("::")))
	assert.True(t, true == iputil.IsIPv6Address(net.ParseIP("fe80::c706:e006:d53e:f9fb")))
	assert.True(t, true == iputil.IsIPv6Address(net.ParseIP("fe80::10.25.21.2")))
	assert.True(t, false == iputil.IsIPv6Address(net.ParseIP("10.25.21.2")))
	assert.True(t, false == iputil.IsIPv6Address(net.ParseIP("0.0.0.0")))

	t.Run("convert uri to hosts with in consistent schema, it should return err",
		func(t *testing.T) {
			_, _, err := iputil.URIs2Hosts([]string{"http://127.0.0.1:8080/", "https://127.0.0.1:8080/"})
			assert.Error(t, err)
		})
	t.Run("convert uri to hosts with consistent schema, it should not return err",
		func(t *testing.T) {
			hosts, s, err := iputil.URIs2Hosts([]string{"http://127.0.0.1:8080/", "http://127.0.0.1:8080/"})
			assert.NoError(t, err)
			assert.Equal(t, "http", s)
			assert.Equal(t, "127.0.0.1:8080", hosts[0])
		})
}

func TestRemoteIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:80/hello", nil)

	r.RemoteAddr = "127.0.0.1:49152"
	assert.EqualValues(t, "127.0.0.1", iputil.RemoteIP(r))
}

func TestForwardedIPs(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:80/hello", nil)

	r.Header.Add("X-Forwarded-For", "127.0.0.1")
	assert.EqualValues(t, []string{"127.0.0.1"}, iputil.ForwardedIPs(r))
}

func TestRealIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:80/hello", nil)

	r.Header.Add("X-Real-Ip", "127.0.0.1")
	assert.EqualValues(t, "127.0.0.1", iputil.RealIP(r))
}

func TestClientIP(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1:80/hello", nil)

	r.Header.Add("X-Real-Ip", "127.0.0.1")
	assert.EqualValues(t, "127.0.0.1", iputil.ClientIP(r))
}
