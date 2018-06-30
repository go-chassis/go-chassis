package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"net"
	"time"
)

const (
	// Name is a constant of type string
	Name = "rest"
	// FailureTypePrefix is a constant of type string
	FailureTypePrefix = "http_"
	//DefaultTimeoutBySecond defines the default timeout for http connections
	DefaultTimeoutBySecond = 60 * time.Second
	//DefaultKeepAliveSecond defines the connection time
	DefaultKeepAliveSecond = 60 * time.Second
	//DefaultMaxConnsPerHost defines the maximum number of concurrent connections
	DefaultMaxConnsPerHost = 512
	//SchemaHTTP represents the http schema
	SchemaHTTP = "http"
	//SchemaHTTPS represents the https schema
	SchemaHTTPS = "https"
)

var (
	//ErrCanceled means Request is canceled by context management
	ErrCanceled = errors.New("request cancelled")
	//ErrInvalidReq invalid input
	ErrInvalidReq = errors.New("rest consumer call arg is not *rest.Request type")
	//ErrInvalidResp invalid input
	ErrInvalidResp = errors.New("rest consumer response arg is not *rest.Response type")
)

//HTTPFailureTypeMap is a variable of type map
var HTTPFailureTypeMap = map[string]bool{
	FailureTypePrefix + strconv.Itoa(http.StatusInternalServerError): true, //http_500
	FailureTypePrefix + strconv.Itoa(http.StatusBadGateway):          true, //http_502
	FailureTypePrefix + strconv.Itoa(http.StatusServiceUnavailable):  true, //http_503
	FailureTypePrefix + strconv.Itoa(http.StatusGatewayTimeout):      true, //http_504
	FailureTypePrefix + strconv.Itoa(http.StatusTooManyRequests):     true, //http_429
}

func init() {
	client.InstallPlugin(Name, NewRestClient)
}

//NewRestClient is a function
func NewRestClient(opts client.Options) client.ProtocolClient {
	if opts.Failure == nil || len(opts.Failure) == 0 {
		opts.Failure = HTTPFailureTypeMap
	} else {
		tmpFailureMap := make(map[string]bool)
		for k := range opts.Failure {
			if HTTPFailureTypeMap[k] {
				tmpFailureMap[k] = true
			}

		}

		opts.Failure = tmpFailureMap
	}

	poolSize := DefaultMaxConnsPerHost
	if opts.PoolSize != 0 {
		poolSize = opts.PoolSize
	}

	tp := &http.Transport{
		MaxIdleConns:        poolSize,
		MaxIdleConnsPerHost: poolSize,
		DialContext: (&net.Dialer{
			KeepAlive: DefaultKeepAliveSecond,
			Timeout:   DefaultTimeoutBySecond,
		}).DialContext}
	if opts.TLSConfig != nil {
		tp.TLSClientConfig = opts.TLSConfig
	}
	rc := &Client{
		opts: opts,
		c: &http.Client{
			Transport: tp,
		},
	}

	return rc
}

//Init is a method

// If a request fails, we generate an error.
func (c *Client) failure2Error(e error, r *Response, addr string) error {
	if e != nil {
		return e
	}

	if r == nil {
		return nil
	}

	codeStr := strconv.Itoa(r.GetStatusCode())
	// The Failure map defines whether or not a request fail.
	if c.opts.Failure["http_"+codeStr] {
		return fmt.Errorf("http error status %d, server addr: %s", r.GetStatusCode(), addr)
	}

	return nil
}
func invocation2HttpRequest(inv *invocation.Invocation) (*Request, error) {
	reqSend, ok := inv.Args.(*Request)
	if !ok {
		return nil, ErrInvalidReq
	}
	return reqSend, nil
}

//Call is a method which uses client struct object
func (c *Client) Call(ctx context.Context, addr string, inv *invocation.Invocation, rsp interface{}) error {
	var err error
	reqSend, err := invocation2HttpRequest(inv)
	if err != nil {
		return err
	}

	resp, ok := rsp.(*Response)
	if !ok {
		return ErrInvalidResp
	}

	c.contextToHeader(ctx, reqSend)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if c.opts.TLSConfig != nil {
		reqSend.Req.URL.Scheme = SchemaHTTPS
	} else {
		reqSend.Req.URL.Scheme = SchemaHTTP
	}
	if addr != "" {
		reqSend.Req.URL.Host = addr
	}

	//increase the max connection per host to prevent error "no free connection available" error while sending more requests.
	c.c.Transport.(*http.Transport).MaxIdleConnsPerHost = 512 * 20

	errChan := make(chan error, 1)
	go func() { errChan <- c.Do(reqSend, resp) }()

	select {
	case <-ctx.Done():
		err = ErrCanceled
	case err = <-errChan:
	}
	return c.failure2Error(err, resp, addr)
}

func (c *Client) String() string {
	return "rest_client"
}

func (c *Client) contextToHeader(ctx context.Context, req *Request) {
	for k, v := range common.FromContext(ctx) {
		req.Req.Header.Set(k, v)
	}

	if len(req.GetContentType()) == 0 {
		req.SetContentType(common.JSON)
	}
}
