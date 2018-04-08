package rest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/ServiceComb/go-chassis/core/client"
)

//Client is a struct
type Client struct {
	c    *http.Client
	opts client.Options
	mu   sync.Mutex // protects following
}

//Do is a method
func (c *Client) Do(req *Request, resp *Response) error {
	c.c.Timeout = DefaultTimoutBySecond * time.Second
	tempResponse, err := c.c.Do(req.Req)
	if err != nil {
		return err
	}
	resp.Resp = tempResponse
	return nil
}

//Request is struct
type Request struct {
	Req *http.Request
}

//NewRequest is a function which creates new request
func NewRequest(method, urlStr string, body ...[]byte) (*Request, error) {
	if method == "" {
		method = "GET"
	}
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body[0])
	}

	req, err := http.NewRequest(method, urlStr, r)
	if err != nil {
		return nil, err
	}
	return &Request{Req: req}, nil
}

//SetURI sets host for the request.
func (req *Request) SetURI(url string) {
	if tempURL, err := req.Req.URL.Parse(url); err == nil {
		req.Req.URL = tempURL
	}
}

//Copy is method
func (req *Request) Copy() *Request {
	newReq, err := http.NewRequest(req.Req.Method, req.Req.URL.String(), req.Req.Body)
	if err != nil {
		return nil
	}
	//Copy headers
	for key := range req.Req.Header {
		newReq.Header.Add(key, req.Req.Header.Get(key))
	}
	//Copy cookies
	for _, c := range req.Req.Cookies() {
		newReq.AddCookie(c)
	}

	return &Request{Req: newReq}
}

//GetRequest is a method
func (req *Request) GetRequest() *http.Request {
	return req.Req
}

//SetBody is a method used for setting body for a request
func (req *Request) SetBody(body []byte) {
	req.Req.Body = ioutil.NopCloser(bytes.NewReader(body))
}

//SetCookie set key value in request cookie
func (req *Request) SetCookie(k, v string) {
	c := &http.Cookie{
		Name:  k,
		Value: v,
	}
	req.Req.AddCookie(c)
}

//GetURI is a method
func (req *Request) GetURI() string {
	return req.Req.URL.String()
}

// SetContentType is a method used for setting content-type in a request
func (req *Request) SetContentType(ct string) {
	req.Req.Header.Set("Content-Type", ct)
}

// GetContentType is a method used for getting content-type in a request
func (req *Request) GetContentType() string {
	return req.Req.Header.Get("Content-Type")
}

//SetHeader is a method used for setting header in a request
func (req *Request) SetHeader(key, value string) {
	req.Req.Header.Set(key, value)
}

//SetHeaderCookie is a method used to setting header cookie
func (req *Request) SetHeaderCookie(key, value string) {
	req.Req.Header.Set(key, value)
}

//GetHeader is a method which gets head from a request
func (req *Request) GetHeader(key string) string {
	return string(req.Req.Header.Get(key))
}

//GetCookie is a method which gets cookie from a request
func (req *Request) GetCookie(key string) string {
	cookie, err := req.Req.Cookie(key)
	if err == http.ErrNoCookie {
		return ""
	}
	return cookie.Value
}

//SetMethod is a method
func (req *Request) SetMethod(method string) {
	req.Req.Method = method
}

//GetMethod is a method
func (req *Request) GetMethod() string {
	return req.Req.Method
}

//Close is used for closing a request
//TODO Confirm it's necessary or not
func (req *Request) Close() {
	//req.Req.Body.Close()
}

//Response is a struct used for handling response
type Response struct {
	Resp *http.Response
}

// NewResponse is creating the object of response
func NewResponse() *Response {

	resp := http.Response{
		Header: http.Header{},
	}
	return &Response{
		Resp: &resp,
	}
}

// GetResponse is a method used to get response
func (resp *Response) GetResponse() *http.Response {
	return resp.Resp
}

// GetStatusCode returns response status code.
func (resp *Response) GetStatusCode() int {
	return resp.Resp.StatusCode
}

// SetStatusCode sets the status code
func (resp *Response) SetStatusCode(s int) {
	resp.Resp.StatusCode = s
}

// ReadBody read body from the from the response
func (resp *Response) ReadBody() []byte {
	if resp.Resp.Body != nil {
		body, err := ioutil.ReadAll(resp.Resp.Body)
		if err != nil {
			return nil
		}
		return body
	}
	return nil
}

// GetHeader get header from the response
//TODO Confirm it's necessary or not
func (resp *Response) GetHeader() []byte {
	bf := new(bytes.Buffer)
	resp.Resp.Header.Write(bf)
	return bf.Bytes()
}

// Close closes the file descriptor
//TODO Confirm it's necessary or not
func (resp *Response) Close() {
	if resp.Resp.Body != nil {
		resp.Resp.Body.Close()
	}
}

// GetCookie returns response Cookie.
func (resp *Response) GetCookie(key string) []byte {
	for _, c := range resp.Resp.Cookies() {
		if c.Name == key {
			return []byte(c.Value)
		}
	}
	return nil
}

// SetCookie sets the cookie.
func (resp *Response) SetCookie(cookie *http.Cookie) {
	resp.Resp.Header.Add("Set-Cookie", cookie.String())
}
