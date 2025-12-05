package iputil_test

import (
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/pkg/util/iputil"
	"github.com/stretchr/testify/assert"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLocalIp(t *testing.T) {
	localip := iputil.GetLocalIP()
	t.Log(localip)
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

	r.Header = make(http.Header, 1)
	r.Header.Add("X-Forwarded-For", "127.0.0.1")

	assert.EqualValues(t, "127.0.0.1", iputil.ClientIP(r))
}

func TestNormalizeAddrWithNetwork(t *testing.T) {
	tests := []struct {
		name               string
		addr               string
		wantNormalizedAddr string
		wantNetwork        string
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name:               "normal ipv4 addr",
			addr:               "127.0.0.1:30110",
			wantNormalizedAddr: "127.0.0.1:30110",
			wantNetwork:        "tcp4",
			wantErr:            assert.NoError,
		},
		{
			name:               "invalid ipv4 addr",
			addr:               "127.0.0.256:30110",
			wantNormalizedAddr: "",
			wantNetwork:        "",
			wantErr:            assert.Error,
		},
		{
			name:               "normal ipv6 addr",
			addr:               "::1:30100",
			wantNormalizedAddr: "[::1]:30100",
			wantNetwork:        "tcp6",
			wantErr:            assert.NoError,
		},
		{
			name:               "invalid ipv6 addr",
			addr:               "::1:301000",
			wantNormalizedAddr: "",
			wantNetwork:        "",
			wantErr:            assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNormalizedAddr, gotNetwork, err := iputil.NormalizeAddrWithNetwork(tt.addr)
			if !tt.wantErr(t, err, fmt.Sprintf("NormalizeAddrWithNetwork(%v)", tt.addr)) {
				return
			}
			assert.Equalf(t, tt.wantNormalizedAddr, gotNormalizedAddr, "NormalizeAddrWithNetwork(%v)", tt.addr)
			assert.Equalf(t, tt.wantNetwork, gotNetwork, "NormalizeAddrWithNetwork(%v)", tt.addr)
		})
	}
}
