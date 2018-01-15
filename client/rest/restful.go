package rest

import (
	"sync"

	clientOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
)

//Client is a struct
type Client struct {
	c    *fasthttp.Client
	opts clientOption.Options
	mu   sync.Mutex // protects following
}

//Do is a method
func (c *Client) Do(req *Request, resp *Response) error {
	return c.c.Do(req.r, resp.r)
}

//Request is struct
type Request struct {
	r *fasthttp.Request
}

//NewRequest is a function which creates new request
func NewRequest(method, urlStr string, body ...[]byte) (*Request, error) {
	if method == "" {
		method = "GET"
	}
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.Header.SetRequestURI(urlStr)
	if body != nil && len(body) == 1 {
		req.SetBody(body[0])
	}

	return &Request{r: req}, nil
}

//SetURI sets host for the request.
func (req *Request) SetURI(url string) {
	req.r.SetRequestURI(url)
}

//Copy is method
func (req *Request) Copy() *Request {
	newReq := fasthttp.AcquireRequest()
	req.r.CopyTo(newReq)
	return &Request{
		r: newReq,
	}
}

//GetRequest is a method
func (req *Request) GetRequest() *fasthttp.Request {
	return req.r
}

//SetBody is a method used for setting body for a request
func (req *Request) SetBody(body []byte) {
	req.r.SetBody(body)
}

//GetURI is a method
func (req *Request) GetURI() string {
	return string(req.r.RequestURI())
}

//SetHeader is a method used for setting header in a request
func (req *Request) SetHeader(key, value string) {
	req.r.Header.Set(key, value)
}

//SetHeaderCookie is a method used to setting header cookie
func (req *Request) SetHeaderCookie(key, value string) {
	req.r.Header.Add(key, value)
}

//GetHeader is a method which gets head from a request
func (req *Request) GetHeader(key string) string {
	return string(req.r.Header.Peek(key))
}

//SetMethod is a method
func (req *Request) SetMethod(method string) {
	req.r.Header.SetMethodBytes([]byte(method))
}

//GetMethod is a method
func (req *Request) GetMethod() string {
	return string(req.r.Header.Method())
}

//Close is used for closing a request
func (req *Request) Close() {
	fasthttp.ReleaseRequest(req.r)
}

//Response is a struct used for handling response
type Response struct {
	r *fasthttp.Response
}

// NewResponse is creating the object of response
func NewResponse() *Response {
	res := fasthttp.AcquireResponse()
	return &Response{
		r: res,
	}
}

// GetResponse is a method used to get response
func (resp *Response) GetResponse() *fasthttp.Response {
	return resp.r
}

// GetStatusCode returns response status code.
func (resp *Response) GetStatusCode() int {
	return resp.r.Header.StatusCode()
}

// SetStatusCode sets the status code
func (resp *Response) SetStatusCode(s int) {
	resp.r.Header.SetStatusCode(s)
}

// ReadBody read body from the from the response
func (resp *Response) ReadBody() []byte {
	return resp.r.Body()
}

// GetHeader get header from the response
func (resp *Response) GetHeader() []byte {
	return resp.r.Header.Header()
}

// Close closes the file descriptor
func (resp *Response) Close() {
	fasthttp.ReleaseResponse(resp.r)
}

// GetCookie returns response Cookie.
func (resp *Response) GetCookie(key string) []byte {
	var c []byte
	resp.r.Header.VisitAllCookie(func(k, v []byte) {
		if string(k) == key {
			c = v
		}
	})
	return c
}

// SetCookie sets the cookie.
func (resp *Response) SetCookie(cookie *fasthttp.Cookie) {
	resp.r.Header.SetCookie(cookie)
}
