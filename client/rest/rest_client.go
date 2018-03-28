package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/common"
)

const (
	// Name is a constant of type string
	Name = "rest"
	// FailureTypePrefix is a constant of type string
	FailureTypePrefix = "http_"
	//DefaultTimoutBySecond defines the default timeout for http connections
	DefaultTimoutBySecond = 60
	//DefaultMaxConnsPerHost defines the maximum number of concurrent connections
	DefaultMaxConnsPerHost = 512
	//SchemaHTTP represents the http schema
	SchemaHTTP = "http"
	//SchemaHTTPS represents the https schema
	SchemaHTTPS = "https"
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

	tp := &http.Transport{}
	if opts.TLSConfig != nil {
		tp.TLSClientConfig = opts.TLSConfig
	}
	// There differences between MaxIdleConnsPerHost and MaxConnsPerHost
	// See https://github.com/golang/go/issues/13957
	tp.MaxIdleConnsPerHost = poolSize
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
func (c *Client) failure2Error(e error, r *Response) error {
	if e != nil {
		return e
	}

	if r == nil {
		return nil
	}

	codeStr := strconv.Itoa(r.GetStatusCode())
	// The Failure map defines whether or not a request fail.
	if c.opts.Failure["http_"+codeStr] {
		return fmt.Errorf("Get error status code: %d from http response: %s", r.GetStatusCode(), string(r.ReadBody()))
	}

	return nil
}

//Call is a method which uses client struct object
func (c *Client) Call(ctx context.Context, addr string, req *client.Request, rsp interface{}) error {
	reqSend, ok := req.Arg.(*Request)
	if !ok {
		return errors.New("Rest consumer call arg is not *rest.Request type")
	}

	resp, ok := rsp.(*Response)
	if !ok {
		return errors.New("Rest consumer response arg is not *rest.Response type")
	}

	c.contextToHeader(ctx, reqSend)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if c.opts.TLSConfig != nil {
		reqSend.Req.URL.Scheme = SchemaHTTPS
	} else {
		reqSend.Req.URL.Scheme = SchemaHTTP
	}

	reqSend.Req.URL.Host = addr

	//increase the max connection per host to prevent error "no free connection available" error while sending more requests.
	c.c.Transport.(*http.Transport).MaxIdleConnsPerHost = 512 * 20

	errChan := make(chan error, 1)
	go func() { errChan <- c.Do(reqSend, resp) }()

	var err error
	select {
	case <-ctx.Done():
		err = errors.New("Request Cancelled")
	case err = <-errChan:
	}
	return c.failure2Error(err, resp)
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
