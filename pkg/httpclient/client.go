package httpclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-chassis/go-chassis/security"
)

//SignRequest sign a http request so that it can talk to API server
//this is global implementation, if you do not set SignRequest in URLClientOption
//client will use this function
var SignRequest func(*http.Request) error

//URLClientOption is a struct which provides options for client
type URLClientOption struct {
	SSLEnabled            bool
	TLSConfig             *tls.Config
	Compressed            bool
	HandshakeTimeout      time.Duration
	ResponseHeaderTimeout time.Duration
	Verbose               bool
	SignRequest           func(*http.Request) error
}

//URLClient is a struct used for storing details of a client
type URLClient struct {
	*http.Client
	TLS     *tls.Config
	Request *http.Request
	options URLClientOption
}

//HTTPDo is a method used for http connection
func (client *URLClient) HTTPDo(method string, rawURL string, headers http.Header, body []byte) (resp *http.Response, err error) {
	client.clientHasPrefix(rawURL, "https")

	if headers == nil {
		headers = make(http.Header)
	}

	if _, ok := headers["Accept"]; !ok {
		headers["Accept"] = []string{"*/*"}
	}
	if _, ok := headers["Accept-Encoding"]; !ok && client.options.Compressed {
		headers["Accept-Encoding"] = []string{"deflate, gzip"}
	}

	req, err := http.NewRequest(method, rawURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	client.Request = req

	req.Header = headers
	//sign a request, first use function in client options
	//if there is not, use global function
	if client.options.SignRequest != nil {
		if err = client.options.SignRequest(req); err != nil {
			return nil, errors.New("Add auth info failed, err: " + err.Error())
		}
	} else if SignRequest != nil {
		if err = SignRequest(req); err != nil {
			return nil, errors.New("Add auth info failed, err: " + err.Error())
		}
	}
	resp, err = client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if client.options.Verbose {
		fmt.Printf("> %s / %s\n", client.Request.Method, client.Request.Proto)
		for key, header := range client.Request.Header {
			for _, value := range header {
				fmt.Printf("> %s: %s\n", key, value)
			}
		}
		fmt.Println(">")
		fmt.Printf("< %s %s\n", resp.Proto, resp.Status)
		for key, header := range resp.Header {
			for _, value := range header {
				fmt.Printf("< %s: %s\n", key, value)
			}
		}
		fmt.Println("<")
	}
	return resp, nil
}

func (client *URLClient) clientHasPrefix(url, pro string) {
	if strings.HasPrefix(url, pro) {
		if transport, ok := client.Client.Transport.(*http.Transport); ok {
			transport.TLSClientConfig = client.TLS
		}
	}
}

//DefaultURLClientOption is a struct object which has default client option
var DefaultURLClientOption = &URLClientOption{
	Compressed:            true,
	HandshakeTimeout:      30 * time.Second,
	ResponseHeaderTimeout: 60 * time.Second,
}

//GetURLClient is a function which which sets client option
func GetURLClient(option *URLClientOption) (client *URLClient, err error) {
	if option == nil {
		option = DefaultURLClientOption
	} else {
		switch {
		case option.HandshakeTimeout == 0:
			option.HandshakeTimeout = DefaultURLClientOption.HandshakeTimeout
			fallthrough
		case option.ResponseHeaderTimeout == 0:
			option.ResponseHeaderTimeout = DefaultURLClientOption.ResponseHeaderTimeout
		}
	}

	if !option.SSLEnabled {
		client = &URLClient{
			Client: &http.Client{
				Transport: &http.Transport{
					TLSHandshakeTimeout:   option.HandshakeTimeout,
					ResponseHeaderTimeout: option.ResponseHeaderTimeout,
					DisableCompression:    !option.Compressed,
				},
			},
			options: *option,
		}

		return
	}

	client = &URLClient{
		Client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout:   option.HandshakeTimeout,
				ResponseHeaderTimeout: option.ResponseHeaderTimeout,
				DisableCompression:    !option.Compressed,
			},
		},
		TLS:     option.TLSConfig,
		options: *option,
	}
	return
}

//GetX509CACertPool is a function used to get certificate
func GetX509CACertPool(caCertFile string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("read ca cert file %s failed", caCertFile)
	}

	pool.AppendCertsFromPEM(caCert)
	return pool, nil
}

//LoadTLSCertificate is a function used to load a certificate
func LoadTLSCertificate(certFile, keyFile, passphase string, cipher security.Cipher) ([]tls.Certificate, error) {
	certContent, err := ioutil.ReadFile(certFile)
	if err != nil {
		errorMsg := "read cert file" + certFile + "failed."
		return nil, errors.New(errorMsg)
	}

	keyContent, err := ioutil.ReadFile(keyFile)
	if err != nil {
		errorMsg := "read key file" + keyFile + "failed."
		return nil, errors.New(errorMsg)
	}

	keyBlock, _ := pem.Decode(keyContent)
	if keyBlock == nil {
		errorMsg := "decode key file " + keyFile + " failed"
		return nil, errors.New(errorMsg)
	}

	plainpass, err := cipher.Decrypt(passphase)
	if err != nil {
		return nil, err
	}

	if x509.IsEncryptedPEMBlock(keyBlock) {
		keyData, err := x509.DecryptPEMBlock(keyBlock, []byte(plainpass))
		if err != nil {
			errorMsg := "decrypt key file " + keyFile + " failed."
			return nil, errors.New(errorMsg)
		}

		// 解密成功，重新编码为无加密的PEM格式文件
		plainKeyBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyData,
		}

		keyContent = pem.EncodeToMemory(plainKeyBlock)
	}

	cert, err := tls.X509KeyPair(certContent, keyContent)
	if err != nil {
		errorMsg := "load X509 key pair from cert file " + certFile + " with key file " + keyFile + " failed."
		return nil, errors.New(errorMsg)
	}

	var certs []tls.Certificate
	certs = append(certs, cert)

	return certs, nil
}
