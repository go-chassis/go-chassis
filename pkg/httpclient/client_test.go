package httpclient_test

import (
	"github.com/go-chassis/go-chassis/pkg/httpclient"
	"github.com/go-chassis/go-chassis/security"
	_ "github.com/go-chassis/go-chassis/security/plugins/aes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpDo(t *testing.T) {

	var htc = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc = new(httpclient.URLClient)
	uc.Client = htc

	htServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	resp, err := uc.HTTPDo(http.MethodGet, htServer.URL, nil, nil)
	assert.NotNil(t, resp)
	assert.NoError(t, err)
}

func TestHttpDoHeadersNil(t *testing.T) {

	var htc *http.Client = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc *httpclient.URLClient = new(httpclient.URLClient)
	uc.Client = htc

	resp, err := uc.HTTPDo("GET", "https://fakeURL", nil, nil)
	assert.Nil(t, resp)
	assert.Error(t, err)

}

func TestHttpDoURLInvalid(t *testing.T) {

	var htc *http.Client = new(http.Client)
	htc.Timeout = time.Second * 2

	var uc *httpclient.URLClient = new(httpclient.URLClient)
	uc.Client = htc

	resp, err := uc.HTTPDo("abc", "url", nil, nil)
	assert.Nil(t, resp)
	assert.Error(t, err)

}
func TestGetURLClient(t *testing.T) {

	tduration := time.Second * 2

	var uc *httpclient.URLClientOption = new(httpclient.URLClientOption)
	uc.Compressed = true
	uc.SSLEnabled = true
	uc.Verbose = true
	uc.HandshakeTimeout = tduration
	uc.ResponseHeaderTimeout = tduration

	c, err := httpclient.GetURLClient(uc)
	expectedc := &httpclient.URLClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   tduration,
				ResponseHeaderTimeout: tduration,
				DisableCompression:    false,
			},
		},
	}

	assert.Equal(t, expectedc.Client, c.Client)
	assert.NoError(t, err)

}

func TestGetURLClientURLClientOptionNil(t *testing.T) {

	option := httpclient.DefaultURLClientOption
	expectedclient := &httpclient.URLClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   option.HandshakeTimeout,
				ResponseHeaderTimeout: option.ResponseHeaderTimeout,
				DisableCompression:    !option.Compressed,
			},
		},
		TLS: option.TLSConfig,
	}

	var uc1 *httpclient.URLClientOption

	c1, err := httpclient.GetURLClient(uc1)

	assert.Equal(t, expectedclient.Client, c1.Client)
	assert.NoError(t, err)

}

func TestGetURLClientSSLEnabledFalse(t *testing.T) {

	tduration := time.Second * 2

	expectedc := &httpclient.URLClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   tduration,
				ResponseHeaderTimeout: tduration,
				DisableCompression:    false,
			},
		},
	}

	var uc2 *httpclient.URLClientOption = new(httpclient.URLClientOption)
	uc2.Compressed = true
	uc2.SSLEnabled = false
	uc2.Verbose = true
	uc2.HandshakeTimeout = tduration
	uc2.ResponseHeaderTimeout = tduration

	c2, err := httpclient.GetURLClient(uc2)

	assert.Equal(t, expectedc.Client, c2.Client)
	assert.NoError(t, err)

}

//func TestGetX509CACertPoolFileNotExist(t *testing.T) {
//	_, err := httpclient.GetX509CACertPool("abc.txt")
//	assert.EqualError(t, err, "read ca cert file abc.txt failed")
//}
//func TestGetX509CACertPoolFileExist(t *testing.T) {
//
//	path := "/home/ca_cert.txt"
//	var _, err = os.Stat(path)
//
//	// create file if not exists
//	if os.IsNotExist(err) {
//		file, err := os.Create(path)
//		assert.NoError(t, err)
//		defer file.Close()
//	}
//
//	_, err = httpclient.GetX509CACertPool(path)
//	assert.NoError(t, err)
//
//	err = os.Remove(path)
//	assert.NoError(t, err)
//}
func TestLoadTLSCertificateFileNotExist(t *testing.T) {

	var cip security.Cipher
	tlsCert, err := httpclient.LoadTLSCertificate("abc.txt", "abc.txt", "fakepassphase", cip)
	assert.Nil(t, tlsCert)
	assert.Error(t, err)
}

//func TestLoadTLSCertificateFileExist(t *testing.T) {
//
//	path := "/home/tls_cert.txt"
//	var _, err = os.Stat(path)
//
//	// create file if not exists
//	if os.IsNotExist(err) {
//		file, err := os.Create(path)
//		assert.NoError(t, err)
//		defer file.Close()
//	}
//	file, err := os.OpenFile(path, os.O_RDWR, 0644)
//	assert.NoError(t, err)
//	defer file.Close()
//
//	// write some text line-by-line to file
//	_, err = file.WriteString(`-----BEGIN RSA PRIVATE KEY-----
//MIIBPAIBAAJBAL7+aty3S1iBA/+yxjxv4q1MUTd1kjNwL4lYKbpzzlmC5beaQXeQ
//2RmGMTXU+mDvuqItjVHOK3DvPK7lTcSGftUCAwEAAQJBALjkK+jc2+iihI98riEF
//oudmkNziSRTYjnwjx8mCoAjPWviB3c742eO3FG4/soi1jD9A5alihEOXfUzloenr
//8IECIQD3B5+0l+68BA/6d76iUNqAAV8djGTzvxnCxycnxPQydQIhAMXt4trUI3nc
//a+U8YL2HPFA3gmhBsSICbq2OptOCnM7hAiEA6Xi3JIQECob8YwkRj29DU3/4WYD7
//WLPgsQpwo1GuSpECICGsnWH5oaeD9t9jbFoSfhJvv0IZmxdcLpRcpslpeWBBAiEA
//6/5B8J0GHdJq89FHwEG/H2eVVUYu5y/aD6sgcm+0Avg=
//-----END RSA PRIVATE KEY-----`)
//	assert.NoError(t, err)
//
//	// save changes
//	err = file.Sync()
//	assert.NoError(t, err)
//	aesFunc, _ := security.GetCipherNewFunc("aes")
//	cipher := aesFunc()
//	s, err := cipher.Encrypt("gochassis")
//	//var cip security.Cipher
//	tlsCert, err := httpclient.LoadTLSCertificate(path, path, s, cipher)
//	assert.Empty(t, tlsCert)
//	assert.Error(t, err)
//
//	err = os.Remove(path)
//	assert.NoError(t, err)
//}
//func TestLoadTLSCertificateCorFileExist(t *testing.T) {
//
//	certpath := "/home/tls_cert.txt"
//	var _, err = os.Stat(certpath)
//
//	// create file if not exists
//	if os.IsNotExist(err) {
//		file, err := os.Create(certpath)
//		assert.NoError(t, err)
//		defer file.Close()
//	}
//	file, err := os.OpenFile(certpath, os.O_RDWR, 0644)
//	assert.NoError(t, err)
//	defer file.Close()
//
//	// write some text line-by-line to file
//	_, err = file.WriteString(`-----BEGIN CERTIFICATE-----
//MIICLDCCAdYCAQAwDQYJKoZIhvcNAQEEBQAwgaAxCzAJBgNVBAYTAlBUMRMwEQYD
//VQQIEwpRdWVlbnNsYW5kMQ8wDQYDVQQHEwZMaXNib2ExFzAVBgNVBAoTDk5ldXJv
//bmlvLCBMZGEuMRgwFgYDVQQLEw9EZXNlbnZvbHZpbWVudG8xGzAZBgNVBAMTEmJy
//dXR1cy5uZXVyb25pby5wdDEbMBkGCSqGSIb3DQEJARYMc2FtcG9AaWtpLmZpMB4X
//DTk2MDkwNTAzNDI0M1oXDTk2MTAwNTAzNDI0M1owgaAxCzAJBgNVBAYTAlBUMRMw
//EQYDVQQIEwpRdWVlbnNsYW5kMQ8wDQYDVQQHEwZMaXNib2ExFzAVBgNVBAoTDk5l
//dXJvbmlvLCBMZGEuMRgwFgYDVQQLEw9EZXNlbnZvbHZpbWVudG8xGzAZBgNVBAMT
//EmJydXR1cy5uZXVyb25pby5wdDEbMBkGCSqGSIb3DQEJARYMc2FtcG9AaWtpLmZp
//MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAL7+aty3S1iBA/+yxjxv4q1MUTd1kjNw
//L4lYKbpzzlmC5beaQXeQ2RmGMTXU+mDvuqItjVHOK3DvPK7lTcSGftUCAwEAATAN
//BgkqhkiG9w0BAQQFAANBAFqPEKFjk6T6CKTHvaQeEAsX0/8YHPHqH/9AnhSjrwuX
//9EBc0n6bVGhN7XaXd6sJ7dym9sbsWxb+pJdurnkxjx4=
//-----END CERTIFICATE-----`)
//	assert.NoError(t, err)
//
//	// save changes
//	err = file.Sync()
//	assert.NoError(t, err)
//
//	keypath := "/home/tls_key.txt"
//	var _, err1 = os.Stat(keypath)
//
//	// create file if not exists
//	if os.IsNotExist(err1) {
//		file, err := os.Create(keypath)
//		assert.NoError(t, err)
//		defer file.Close()
//	}
//	file, err = os.OpenFile(keypath, os.O_RDWR, 0644)
//	assert.NoError(t, err)
//	defer file.Close()
//
//	// write some text line-by-line to file
//	_, err = file.WriteString(`-----BEGIN RSA PRIVATE KEY-----
//MIIBPAIBAAJBAL7+aty3S1iBA/+yxjxv4q1MUTd1kjNwL4lYKbpzzlmC5beaQXeQ
//2RmGMTXU+mDvuqItjVHOK3DvPK7lTcSGftUCAwEAAQJBALjkK+jc2+iihI98riEF
//oudmkNziSRTYjnwjx8mCoAjPWviB3c742eO3FG4/soi1jD9A5alihEOXfUzloenr
//8IECIQD3B5+0l+68BA/6d76iUNqAAV8djGTzvxnCxycnxPQydQIhAMXt4trUI3nc
//a+U8YL2HPFA3gmhBsSICbq2OptOCnM7hAiEA6Xi3JIQECob8YwkRj29DU3/4WYD7
//WLPgsQpwo1GuSpECICGsnWH5oaeD9t9jbFoSfhJvv0IZmxdcLpRcpslpeWBBAiEA
//6/5B8J0GHdJq89FHwEG/H2eVVUYu5y/aD6sgcm+0Avg=
//-----END RSA PRIVATE KEY-----`)
//	assert.NoError(t, err)
//
//	// save changes
//	err = file.Sync()
//	assert.NoError(t, err)
//
//	aesFunc, _ := security.GetCipherNewFunc("aes")
//	cipher := aesFunc()
//	s, err := cipher.Encrypt("gochassis")
//	//var cip security.Cipher
//	tlsCert, err := httpclient.LoadTLSCertificate(certpath, keypath, s, cipher)
//	assert.NotEmpty(t, tlsCert)
//	assert.NoError(t, err)
//
//	err = os.Remove(certpath)
//	assert.NoError(t, err)
//
//	err = os.Remove(keypath)
//	assert.NoError(t, err)
//}
//func TestLoadTLSKeyCorFileExist(t *testing.T) {
//	certpath := "/home/tls_cert.txt"
//	var _, err = os.Stat(certpath)
//
//	// create file if not exists
//	if os.IsNotExist(err) {
//		file, err := os.Create(certpath)
//		assert.NoError(t, err)
//		defer file.Close()
//	}
//	file, err := os.OpenFile(certpath, os.O_RDWR, 0644)
//	assert.NoError(t, err)
//	defer file.Close()
//
//	// write some text line-by-line to file
//	_, err = file.WriteString(`-----BEGIN CERTIFICATE-----
//MIICLDCCAdYCAQAwDQYJKoZIhvcNAQEEBQAwgaAxCzAJBgNVBAYTAlBUMRMwEQYD
//VQQIEwpRdWVlbnNsYW5kMQ8wDQYDVQQHEwZMaXNib2ExFzAVBgNVBAoTDk5ldXJv
//bmlvLCBMZGEuMRgwFgYDVQQLEw9EZXNlbnZvbHZpbWVudG8xGzAZBgNVBAMTEmJy
//dXR1cy5uZXVyb25pby5wdDEbMBkGCSqGSIb3DQEJARYMc2FtcG9AaWtpLmZpMB4X
//DTk2MDkwNTAzNDI0M1oXDTk2MTAwNTAzNDI0M1owgaAxCzAJBgNVBAYTAlBUMRMw
//EQYDVQQIEwpRdWVlbnNsYW5kMQ8wDQYDVQQHEwZMaXNib2ExFzAVBgNVBAoTDk5l
//dXJvbmlvLCBMZGEuMRgwFgYDVQQLEw9EZXNlbnZvbHZpbWVudG8xGzAZBgNVBAMT
//EmJydXR1cy5uZXVyb25pby5wdDEbMBkGCSqGSIb3DQEJARYMc2FtcG9AaWtpLmZp
//MFwwDQYJKoZIhvcNAQEBBQADSwAwSAJBAL7+aty3S1iBA/+yxjxv4q1MUTd1kjNw
//L4lYKbpzzlmC5beaQXeQ2RmGMTXU+mDvuqItjVHOK3DvPK7lTcSGftUCAwEAATAN
//BgkqhkiG9w0BAQQFAANBAFqPEKFjk6T6CKTHvaQeEAsX0/8YHPHqH/9AnhSjrwuX
//9EBc0n6bVGhN7XaXd6sJ7dym9sbsWxb+pJdurnkxjx4=
//-----END CERTIFICATE-----`)
//	assert.NoError(t, err)
//
//	// save changes
//	err = file.Sync()
//	assert.NoError(t, err)
//
//	keypath := "/home/tls_key.txt"
//	var _, err1 = os.Stat(keypath)
//
//	// create file if not exists
//	if os.IsNotExist(err1) {
//		file, err := os.Create(keypath)
//		assert.NoError(t, err)
//		defer file.Close()
//	}
//	file, err = os.OpenFile(keypath, os.O_RDWR, 0644)
//	assert.NoError(t, err)
//	defer file.Close()
//
//	// write some text line-by-line to file
//
//	_, err = file.WriteString(`-----BEGIN PUBLIC KEY-----
//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAneF3oCSg1XllOgeQyfal
//ph+EHCMHS0+lA8YP91TVi355gQDS6T30l/6EzVW9yY8hV4gGOZBmQSZ5LMo/lYcB
//ES8vsOELQ/xfL09nBNtNt3JN0cV2c02RabBxFzbqqwo6zZWbdhuOIRePxQK/JMfA
//QLE7xIB8caVR3Pc6WH+xB4GKENH2kxdx4PpReRXU14+tvW844SZ9vPA+gIm07I5p
//kNuXivAjI4OCO2qxrOvnmXQqNY6pZP1GnujlSGExbub8GRhUwxtP1gBEhxw3Rer1
//ycsPDFXsz2rCRSYjojFSTe4hff1YcsIoxY6p0O4Bdwil8CIrR3krz5pGtY/9ZKK1
//7QIDAQAB
//-----END PUBLIC KEY-----`)
//	assert.NoError(t, err)
//
//	// save changes
//	err = file.Sync()
//	assert.NoError(t, err)
//
//	aesFunc, _ := security.GetCipherNewFunc("aes")
//	cipher := aesFunc()
//	s, err := cipher.Encrypt("gochassis")
//	//var cip security.Cipher
//	tlsCert, err := httpclient.LoadTLSCertificate(certpath, keypath, s, cipher)
//	assert.NotEmpty(t, tlsCert)
//	assert.NoError(t, err)
//
//	err = os.Remove(certpath)
//	assert.NoError(t, err)
//
//	err = os.Remove(keypath)
//	assert.NoError(t, err)
//}
//func TestLoadTLSCertificateFileExistDecodeFail(t *testing.T) {
//
//	path := "/home/tls_cert.txt"
//	var _, err = os.Stat(path)
//
//	// create file if not exists
//	if os.IsNotExist(err) {
//		file, err := os.Create(path)
//		assert.NoError(t, err)
//		defer file.Close()
//	}
//
//	aesFunc, _ := security.GetCipherNewFunc("aes")
//	cipher := aesFunc()
//	s, err := cipher.Encrypt("gochassis")
//	//var cip security.Cipher
//	tlsCert, err := httpclient.LoadTLSCertificate(path, path, s, cipher)
//	assert.Empty(t, tlsCert)
//	assert.Error(t, err)
//
//	err = os.Remove(path)
//	assert.NoError(t, err)
//}

//func TestLoadTLSCertificateKeyFileNotExist(t *testing.T) {
//
//	path := "/home/tls_cert.txt"
//	var _, err = os.Stat(path)
//
//	// create file if not exists
//	if os.IsNotExist(err) {
//		file, err := os.Create(path)
//		assert.NoError(t, err)
//		defer file.Close()
//	}
//
//	file, err := os.OpenFile(path, os.O_RDWR, 0644)
//	assert.NoError(t, err)
//	defer file.Close()
//
//	// write some text line-by-line to file
//	_, err = file.WriteString(`-----BEGIN RSA PRIVATE KEY-----
//MIIBPAIBAAJBAN+FmbxmHVOp/RxtpMGz0DvQEBz1sDktHp19hIoMSu0YZift5MAu
//4xAEJYvWVCshDiyOTWsUBXwZkrkt87FyctkCAwEAAQJAG/vxBGpQb6IPo1iC0RF/
//F430BnwoBPCGLbeCOXpSgx5X+19vuTSdEqMgeNB6+aNb+XY/7mvVfCjyD6WZ0oxs
//JQIhAPO+uL9cP40lFs62pdL3QSWsh3VNDByvOtr9LpeaxBm/AiEA6sKVfXsDQ5hd
//SHt9U61r2r8Lcxmzi9Kw6JNqjMmzqWcCIQCKoRy+aZ8Tjdas9yDVHh+FZ90bEBkl
//b1xQFNOdEj8aTQIhAOJWrO6INYNsWTPS6+hLYZtLamyUsQj0H+B8kNQge/mtAiEA
//nBfvUl243qbqN8gF7Az1u33uc9FsPVvQPiBzLxZ4ixw=
//-----END RSA PRIVATE KEY-----`)
//	assert.NoError(t, err)
//
//	// save changes
//	err = file.Sync()
//	assert.NoError(t, err)
//	aesFunc, _ := security.GetCipherNewFunc("aes")
//	cipher := aesFunc()
//	s, err := cipher.Encrypt("gochassis")
//	//var cip security.Cipher
//	tlsCert, err := httpclient.LoadTLSCertificate(path, "abc.txt", s, cipher)
//	assert.Empty(t, tlsCert)
//	assert.Error(t, err)
//
//	err = os.Remove(path)
//	assert.NoError(t, err)
//}
