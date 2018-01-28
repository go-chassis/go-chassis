package tcp

import (
	"crypto/tls"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ServiceComb/go-chassis/core/lager"
	microTransport "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"github.com/stretchr/testify/assert"
)

func TestTCPTransportPortRange(t *testing.T) {
	t.Log("Testing for tcp Listen function")
	tp := NewTransport()

	ln1, err := tp.Listen(":55555")
	assert.NoError(t, err)

	parts := strings.Split(ln1.Addr(), ":")
	assert.Equal(t, "55555", parts[len(parts)-1])

	ln, err := tp.Listen(":0")
	assert.NoError(t, err)

	ln.Close()
	ln1.Close()
}

func TestTCPTransportCommunication(t *testing.T) {
	t.Log("Testing for tcp communication functions")
	os.Setenv("CHASSIS_HOME", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication", "client"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	tp := NewTransport()
	l, err := tp.Listen("127.0.0.1:9999")
	assert.NoError(t, err)
	defer l.Close()

	fn := func(sock microTransport.Socket) {
		defer sock.Close()

		for {

			responseHeader, responseBody, _, ID, err := sock.Recv()
			if err != nil {
				lager.Logger.Errorf(err, "server receive err")
				return
			}
			if err := sock.Send(responseHeader, responseBody, nil, ID); err != nil {
				lager.Logger.Errorf(err, "server send err")
				return
			}
		}
	}

	finish := make(chan bool)

	go func() {
		if err := l.Accept(fn); err != nil {
			select {
			case <-finish:
			default:
				t.Errorf("Unexpected err: %v", err)
			}
		}
	}()
	c, err := tp.Dial(l.Addr())
	assert.NoError(t, err)
	defer c.Close()

	var requestBody = []byte("ms name")
	metadata := make(map[string]string)
	metadata["requestID"] = "0"

	if err := c.Send(nil, requestBody, nil, 0); err != nil {
		t.Errorf("Unexpected err: %v", err)
		return

	}
	_, respBody, _, _, err := c.Recv()
	assert.NoError(t, err)
	assert.Equal(t, respBody, requestBody)

	close(finish)
	l.Close()
}

func TestTCPTransportError(t *testing.T) {
	t.Log("Testing for transport send function")
	tp := NewTransport()

	l, err := tp.Listen("127.0.0.1:9989")
	assert.NoError(t, err)
	defer l.Close()

	fn := func(sock microTransport.Socket) {
		defer sock.Close()

		for {
			_, _, _, _, err := sock.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				t.Fatal(err)
			}
		}
	}

	finish := make(chan bool)

	go func() {
		if err := l.Accept(fn); err != nil {
			select {
			case <-finish:
			default:
				t.Errorf("Unexpected err: %v", err)
			}
		}
	}()

	c, err := tp.Dial(l.Addr())
	assert.NoError(t, err)
	defer c.Close()

	var requestBody = []byte("ms name")
	//metadata := make(map[string]string)
	//metadata["requestID"] = "0"

	err = c.Send(nil, requestBody, nil, 0)
	assert.NoError(t, err)

	close(finish)
	l.Close()
}

func TestTCPTransportTimeout(t *testing.T) {
	t.Log("Testing for transport send function timeout")
	tr := NewTransport(microTransport.Timeout(time.Millisecond * 100))

	l, err := tr.Listen("127.0.0.1:9979")
	assert.NoError(t, err)
	defer l.Close()

	finish := make(chan bool)

	fn := func(sock microTransport.Socket) {
		defer func() {
			sock.Close()
			close(finish)
		}()

		go func() {
			select {
			case <-finish:
				return
			case <-time.After(1 * time.Second):
				t.Fatal("deadline not executed")
			}
		}()

		for {
			_, _, _, _, err := sock.Recv()
			if err != nil {
				return
			}
		}
	}

	go func() {
		if err := l.Accept(fn); err != nil {
			select {
			case <-finish:
			default:
				t.Errorf("Unexpected err: %v", err)
			}
		}
	}()

	c, err := tr.Dial(l.Addr())
	assert.NoError(t, err)
	defer c.Close()

	var requestBody = []byte("ms name")
	//metadata := make(map[string]string)
	//metadata["requestID"] = "0"

	err = c.Send(nil, requestBody, nil, 0)
	assert.NoError(t, err)

	<-finish
	l.Close()
}

func BenchmarkTcpTransportClient_Send(b *testing.B) {
	os.Setenv("CHASSIS_HOME", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication", "client"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	tp := NewTransport()

	l, err := tp.Listen("127.0.0.1:9991")
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fn := func(sock microTransport.Socket) {
		defer sock.Close()

		for {
			responseHeader, responseBody, _, ID, err := sock.Recv()
			if err != nil {
				lager.Logger.Errorf(err, "server receive err")
				return
			}
			if err := sock.Send(responseHeader, responseBody, nil, ID); err != nil {
				lager.Logger.Errorf(err, "server send err")
				return
			}
		}
	}

	done := make(chan bool)

	go func() {
		if err := l.Accept(fn); err != nil {
			select {
			case <-done:
			default:
				panic(err)
			}
		}
	}()
	c, err := tp.Dial(l.Addr())
	if err != nil {
		panic(err)
	}
	defer c.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var requestBody = []byte("ms name")
		metadata := make(map[string]string)
		metadata["requestID"] = "0"

		if err := c.Send(nil, requestBody, nil, 0); err != nil {
			panic(err)

		}

		_, respBody, _, _, err := c.Recv()
		if err != nil {
			panic(err)
		}
		if string(respBody) != string(requestBody) {
			panic("Expected, got")
		}

	}

	close(done)
	l.Close()
}
func TestTCPTransportCommunicationSendError(t *testing.T) {
	t.Log("Testing for transport listen function with error")
	tp := NewTransport(func(abc *microTransport.Options) {
		abc.Secure = true
	})
	_, err := tp.Listen(":55555", func(abc *microTransport.ListenOptions) {})
	assert.Error(t, err)
}
func TestTCPTransportDial(t *testing.T) {
	t.Log("Testing for transport Dial function")
	os.Setenv("CHASSIS_HOME", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication", "client"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	tp := NewTransport()
	l, err := tp.Listen("127.0.0.1:9990")
	assert.NoError(t, err)
	defer l.Close()

	c, err := tp.Dial(l.Addr(), func(abc *microTransport.DialOptions) {})
	assert.NoError(t, err)
	defer c.Close()
}
func TestTCPTransportDialSecureTrue(t *testing.T) {

	//	os.Setenv("CHASSIS_HOME", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication", "client"))
	//	lager.Initialize()
	certContent := `-----BEGIN CERTIFICATE-----
MIICLDCCAdYCAQAwDQYJKoZIhvcNAQEEBQAwgaAxCzAJBgNVBAYTAlBUMRMwEQYD
VQQIEwpRdWVlbnNsYW5kMQ8wDQYDVQQHEwZMaXNib2ExFzAVBgNVBAoTDk5ldXJv
bmlvLCBMZGEuMRgwFgYDVQQLEw9EZXNlbnZvbHZpbWVudG8xGzAZBgNVBAMTEmJy
dXR1cy5uZXVyb25pby5wdDEbMBkGCSqGSIb3DQEJARYMc2FtcG9AaWtpLmZpMB4X
DTk2MDkwNTAzNDI0M1oXDTk2MTAwNTAzNDI0M1owgaAxCzAJBgNVBAYTAlBUMRMw
EQYDVQQIEwpRdWVlbnNsYW5kMQ8wDQYDVQQHEwZMaXNib2ExFzAVBgNVBAoTDk5l
dXJvbmlvLCBMZGEuMRgwFgYDVQQLEw9EZXNlbnZvbHZpbWVudG8xGzAZBgNVBAMT
EmJydXR1cy5uZXVyb25pby5wdDEbMBkGCSqGSIb3DQEJARYMc2FtcG9AaWtpLmZp
MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAL7+aty3S1iBA/+yxjxv4q1MUTd1kjNw
L4lYKbpzzlmC5beaQXeQ2RmGMTXU+mDvuqItjVHOK3DvPK7lTcSGftUCAwEAATAN
BgkqhkiG9w0BAQQFAANBAFqPEKFjk6T6CKTHvaQeEAsX0/8YHPHqH/9AnhSjrwuX
9EBc0n6bVGhN7XaXd6sJ7dym9sbsWxb+pJdurnkxjx4=
-----END CERTIFICATE-----`

	keyContent := `-----BEGIN RSA PRIVATE KEY-----
MIIBPAIBAAJBAL7+aty3S1iBA/+yxjxv4q1MUTd1kjNwL4lYKbpzzlmC5beaQXeQ
2RmGMTXU+mDvuqItjVHOK3DvPK7lTcSGftUCAwEAAQJBALjkK+jc2+iihI98riEF
oudmkNziSRTYjnwjx8mCoAjPWviB3c742eO3FG4/soi1jD9A5alihEOXfUzloenr
8IECIQD3B5+0l+68BA/6d76iUNqAAV8djGTzvxnCxycnxPQydQIhAMXt4trUI3nc
a+U8YL2HPFA3gmhBsSICbq2OptOCnM7hAiEA6Xi3JIQECob8YwkRj29DU3/4WYD7
WLPgsQpwo1GuSpECICGsnWH5oaeD9t9jbFoSfhJvv0IZmxdcLpRcpslpeWBBAiEA
6/5B8J0GHdJq89FHwEG/H2eVVUYu5y/aD6sgcm+0Avg=
-----END RSA PRIVATE KEY-----`

	cert, err := tls.X509KeyPair([]byte(certContent), []byte(keyContent))
	assert.Nil(t, err)

	var certs = []tls.Certificate{cert}
	var serverTlsConfig *tls.Config = &tls.Config{
		Certificates:             certs,
		CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		PreferServerCipherSuites: true,
		ClientAuth:               tls.NoClientCert,
		MinVersion:               tls.VersionTLS12,
		MaxVersion:               tls.VersionTLS12,
	}

	var clientTlsConfig *tls.Config = &tls.Config{
		Certificates:       certs,
		CipherSuites:       []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS12,
	}

	t.Log("Testing for Transport Dial Secure true")
	tr := NewTransport(
		microTransport.Secure(true),
		microTransport.TLSConfig(serverTlsConfig))

	l, err := tr.Listen("127.0.0.1:9988")
	assert.NoError(t, err)

	fn := func(sock microTransport.Socket) {
		defer sock.Close()

		for {

			responseHeader, responseBody, _, ID, err := sock.Recv()
			if err != nil {
				lager.Logger.Errorf(err, "server receive err")
				return
			}
			if err := sock.Send(responseHeader, responseBody, nil, ID); err != nil {
				lager.Logger.Errorf(err, "server send err")
				return
			}
		}
	}

	finish := make(chan bool)

	go func() {
		if err := l.Accept(fn); err != nil {
			select {
			case <-finish:
			default:
				t.Errorf("Unexpected err: %v", err)
			}
		}
	}()

	clientTr := NewTransport(
		microTransport.Secure(true),
		microTransport.TLSConfig(clientTlsConfig))
	c, err := clientTr.Dial(l.Addr())
	assert.NoError(t, err)

	defer c.Close()
}
func TestTCPTransportString(t *testing.T) {
	t.Log("Testing transport string function")
	tp := NewTransport()
	protocol := tp.String()
	assert.Equal(t, protocol, "tcp")
}
func TestTCPTransportDialError(t *testing.T) {
	t.Log("Testing Dial function with errors")
	os.Setenv("CHASSIS_HOME", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication", "client"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	tp := NewTransport()
	l, err := tp.Listen("127.0.0.1:9991")
	assert.NoError(t, err)
	defer l.Close()

	_, err = tp.Dial("abc")
	assert.Error(t, err)
}
func TestTCPTransportListenError(t *testing.T) {
	t.Log("Testing transport Listen function with error")
	tp := NewTransport(func(abc *microTransport.Options) {
		abc.Secure = true
	})

	_, err := tp.Listen("abc")
	assert.Error(t, err)
}
