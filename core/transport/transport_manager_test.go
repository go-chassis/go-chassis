package transport_test

import (
	"crypto/tls"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
	"github.com/ServiceComb/go-chassis/core/transport"
	"github.com/ServiceComb/go-chassis/security"
	securityCommon "github.com/ServiceComb/go-chassis/security/common"
	_ "github.com/ServiceComb/go-chassis/security/plugins/aes"
	_ "github.com/ServiceComb/go-chassis/security/plugins/plain"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	microTransport "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	_ "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport/tcp"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransport(t *testing.T) {
	os.Setenv("CHASSIS_HOME", filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "communication", "client"))
	config.Init()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	testDir := filepath.Join(os.Getenv("GOPATH"), "test", "transport", "TestCreateTransport")
	tlsFileDir := filepath.Join(testDir, "tls")
	err := os.MkdirAll(tlsFileDir, 0600)
	assert.NoError(t, err)
	certpath := filepath.Join(tlsFileDir, "tls_cert.txt")
	_, err = os.Stat(certpath)

	// create file if not exists
	if os.IsNotExist(err) {
		file, err := os.Create(certpath)
		assert.NoError(t, err)
		defer file.Close()
	}
	file, err := os.OpenFile(certpath, os.O_RDWR, 0644)
	assert.NoError(t, err)
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(`-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`)
	assert.NoError(t, err)

	// save changes
	err = file.Sync()
	assert.NoError(t, err)

	keypath := filepath.Join(tlsFileDir, "tls_key.txt")
	var _, err1 = os.Stat(keypath)

	// create file if not exists
	if os.IsNotExist(err1) {
		file, err := os.Create(keypath)
		assert.NoError(t, err)
		defer file.Close()
	}
	file, err = os.OpenFile(keypath, os.O_RDWR, 0644)
	assert.NoError(t, err)
	defer file.Close()

	// write some text line-by-line to file

	_, err = file.WriteString(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAneF3oCSg1XllOgeQyfal
ph+EHCMHS0+lA8YP91TVi355gQDS6T30l/6EzVW9yY8hV4gGOZBmQSZ5LMo/lYcB
ES8vsOELQ/xfL09nBNtNt3JN0cV2c02RabBxFzbqqwo6zZWbdhuOIRePxQK/JMfA
QLE7xIB8caVR3Pc6WH+xB4GKENH2kxdx4PpReRXU14+tvW844SZ9vPA+gIm07I5p
kNuXivAjI4OCO2qxrOvnmXQqNY6pZP1GnujlSGExbub8GRhUwxtP1gBEhxw3Rer1
ycsPDFXsz2rCRSYjojFSTe4hff1YcsIoxY6p0O4Bdwil8CIrR3krz5pGtY/9ZKK1
7QIDAQAB
-----END PUBLIC KEY-----`)
	assert.NoError(t, err)

	// save changes
	err = file.Sync()
	assert.NoError(t, err)

	sslRootPath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "tls", "etc")
	os.Setenv("CIPHER_ROOT", filepath.Join(sslRootPath, "cipher"))
	os.Setenv("PAAS_CRYPTO_PATH", filepath.Join(sslRootPath, "cipher"))

	aesFunc, err := security.GetCipherNewFunc("default")
	assert.NoError(t, err)
	//fmt.Println("@@@@@@@@@@@@@@@",security.CipherPlugins["aes"])
	cipher := aesFunc()
	s, err := cipher.Encrypt("gochassis")

	keypwdpath := filepath.Join(tlsFileDir, "pwd_key.txt")
	var _, err2 = os.Stat(keypwdpath)

	// create file if not exists
	if os.IsNotExist(err2) {
		file, err := os.Create(keypwdpath)
		assert.NoError(t, err)
		defer file.Close()
	}
	file, err = os.OpenFile(keypwdpath, os.O_RDWR, 0644)
	assert.NoError(t, err)
	defer file.Close()

	// write some text line-by-line to file

	_, err = file.WriteString(s)
	assert.NoError(t, err)

	// save changes
	err = file.Sync()
	assert.NoError(t, err)

	/*	sslRootPath := filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "ServiceComb", "go-chassis", "examples", "tls", "etc")
		os.Setenv("CIPHER_ROOT", filepath.Join(sslRootPath, "cipher"))
		os.Setenv("PAAS_CRYPTO_PATH", filepath.Join(sslRootPath, "cipher"))*/
	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Ssl = make(map[string]string)
	sslConfig := chassisTLS.GetDefaultSSLConfig()

	sslConfig.VerifyPeer = false
	sslConfig.VerifyPeer = false
	sslConfig.CAFile = filepath.Join(sslRootPath, "ssl", "trust.cer")
	sslConfig.CertFile = certpath
	sslConfig.KeyFile = keypath
	sslConfig.CertPWDFile = keypwdpath
	sslConfig.CipherPlugin = "aes"

	tlsConfig, err := securityCommon.GetServerTLSConfig(sslConfig)
	assert.Error(t, err)

	transport.InstallPlugin("abc", nil)

	trF, err := transport.GetTransportFunc("tcp")
	assert.NoError(t, err)

	trF1, err := transport.GetTransportFunc("abc")
	assert.Nil(t, trF1)
	assert.Error(t, err)

	serverTr := trF(microTransport.TLSConfig(tlsConfig))
	l, err := serverTr.Listen("127.0.0.1:9991")
	assert.NoError(t, err)
	defer l.Close()

	fn := func(sock microTransport.Socket) {
		defer sock.Close()

		for {
			responseHeader, responseBody, _, ID, err := sock.Recv()
			if err != nil {
				return
			}

			if err := sock.Send(responseHeader, responseBody, nil, ID); err != nil {
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
				t.Errorf("Accept err: %v", err)
			}
		}
	}()
	tlsConfig, err = securityCommon.GetClientTLSConfig(sslConfig)
	assert.Error(t, err)
	clietTr := trF(microTransport.TLSConfig(tlsConfig))
	c, err := clietTr.Dial(l.Addr())
	assert.NoError(t, err)
	defer c.Close()

	//metadata := make(map[string]string)
	//metadata["requestID"] = "0"

	var requestBody = []byte("test")
	err = c.Send(nil, requestBody, nil, 0)
	assert.NoError(t, err)

	_, respBody, _, _, err := c.Recv()
	assert.NoError(t, err)

	close(finish)
	l.Close()
	assert.Equal(t, requestBody, respBody, "they should be equal")
	//
	//// assert inequality
	//assert.NotEqual(t, 123, 456, "they should not be equal")
	//
	//// assert for nil (good for errors)
	//assert.Nil(t, object)
	//
	//// assert for not nil (good when you expect something)
	//if assert.NotNil(t, object) {
	//
	//	// now we know that object isn't nil, we are safe to make
	//	// further assertions without causing any errors
	//	assert.Equal(t, "Something", object.Value)
	//
	//}

}
func TestCreateTransportFunc(t *testing.T) {
	transport.CreateTransport("tcp")
	transport.Init()
	tr := transport.GetTransport("tcp")
	assert.NotNil(t, tr)

	var o *microTransport.Options = new(microTransport.Options)
	var c codec.Codec
	var dopt *microTransport.DialOptions = new(microTransport.DialOptions)
	var tls1 = new(tls.Config)
	tls1.DynamicRecordSizingDisabled = true
	t1 := time.Second * 10

	var addrArr = []string{"abc"}
	op := microTransport.Addrs(addrArr...)
	op(o)
	assert.Equal(t, o.Addrs, addrArr)

	op = microTransport.Codec(c)
	op(o)
	assert.Equal(t, c, o.Codec)

	op = microTransport.Timeout(t1)
	op(o)
	assert.Equal(t, t1, o.Timeout)

	op = microTransport.Secure(true)
	op(o)
	assert.Equal(t, true, o.Secure)

	op = microTransport.TLSConfig(tls1)
	op(o)
	assert.Equal(t, tls1.DynamicRecordSizingDisabled, o.TLSConfig.DynamicRecordSizingDisabled)

	dop := microTransport.WithStream()
	dop(dopt)
	assert.Equal(t, true, dopt.Stream)

	dop = microTransport.WithTimeout(t1)
	dop(dopt)
	assert.Equal(t, t1, dopt.Timeout)

}
