package rest

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/go-chassis/go-chassis/v2/core/client"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
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
	DefaultMaxConnsPerHost = 512 * 20
	//SchemaHTTP represents the http schema
	SchemaHTTP = "http"
	//SchemaHTTPS represents the https schema
	SchemaHTTPS = "https"
)

var (

	//ErrInvalidResp invalid input
	ErrInvalidResp = errors.New("rest consumer response arg is not *rest.Response type")
)

func init() {
	client.InstallPlugin(Name, NewRestClient)
}

// Client is a struct
type Client struct {
	c    *http.Client
	opts client.Options
}

func (c *Client) Status(rsp interface{}) (status int, err error) {
	if resp, ok := rsp.(*http.Response); ok {
		return resp.StatusCode, nil
	}
	return 0, fmt.Errorf("imcompatible type: %s", reflect.TypeOf(rsp))
}

// NewRestClient is a function
func NewRestClient(opts client.Options) (client.ProtocolClient, error) {
	tp := newTransport(opts)
	rc := &Client{
		opts: opts,

		c: &http.Client{
			Timeout:       opts.Timeout,
			Transport:     tp,
			CheckRedirect: opts.CheckRedirect,
		},
	}
	return rc, nil
}

func newTransport(opts client.Options) *http.Transport {
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
	return tp
}

// If a request fails, we generate an error.
func (c *Client) failure2Error(e error, r *http.Response, addr string) error {
	if e != nil {
		return e
	}
	if c.opts.Failure == nil {
		return nil
	}
	if r == nil {
		return nil
	}

	codeStr := strconv.Itoa(r.StatusCode)
	// The Failure map defines whether or not a request fail.
	if c.opts.Failure["http_"+codeStr] {
		return fmt.Errorf("http error status [%d], server addr: [%s], will not print response body, to protect service sensitive data", r.StatusCode, addr)
	}

	return nil
}

// Call is a method which uses client struct object
func (c *Client) Call(ctx context.Context, addr string, inv *invocation.Invocation, rsp interface{}) error {
	var err error
	reqSend, err := httputil.HTTPRequest(inv)
	if err != nil {
		return err
	}
	resp, ok := rsp.(*http.Response)
	if !ok {
		return ErrInvalidResp
	}

	c.contextToHeader(ctx, reqSend)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	if c.opts.TLSConfig != nil {
		reqSend.URL.Scheme = SchemaHTTPS
	} else {
		reqSend.URL.Scheme = SchemaHTTP
	}
	if addr != "" {
		reqSend.URL.Host = addr
	}

	var temp *http.Response
	errChan := make(chan error, 1)
	go func() {
		temp, err = c.c.Do(reqSend)
		errChan <- err
	}()

	select {
	case <-ctx.Done():
		err = client.ErrCanceled
	case err = <-errChan:
		if err == nil {
			*resp = *temp
		}
	}

	return c.failure2Error(err, resp, addr)
}

func (c *Client) String() string {
	return "rest_client"
}

// Close release the idle connection
func (c *Client) Close() error {
	c.c.CloseIdleConnections()
	return nil
}

// ReloadConfigs  reload configs for timeout and tls
func (c *Client) ReloadConfigs(opts client.Options) {
	c.opts = client.EqualOpts(c.opts, opts)
	c.c.Timeout = c.opts.Timeout
	tp := newTransport(opts)
	c.c.Transport = tp
}

// GetOptions method return opts
func (c *Client) GetOptions() client.Options {
	return c.opts
}

func (c *Client) contextToHeader(ctx context.Context, req *http.Request) {
	for k, v := range common.FromContext(ctx) {
		req.Header.Set(k, v)
	}

	if len(req.Header.Get("Content-Type")) == 0 {
		req.Header.Set("Content-Type", common.JSON)
	}
}
