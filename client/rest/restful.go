package rest

import (
	"bytes"
	"github.com/go-chassis/go-chassis/core/client"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

//NewRequest is a function which creates new request
func NewRequest(method, urlStr string, body []byte) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}

	req, err := http.NewRequest(method, urlStr, r)
	if err != nil {
		return nil, err
	}
	return req, nil
}

//Client is a struct
type Client struct {
	c    *http.Client
	opts client.Options
	mu   sync.Mutex // protects following
}

//Do is a method
func (c *Client) Do(req *http.Request, resp *Response) (err error) {
	resp.Resp, err = c.c.Do(req)
	return
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
	if resp.Resp != nil {
		return resp.Resp.StatusCode
	}
	return 0
}

// SetStatusCode sets the status code
func (resp *Response) SetStatusCode(s int) {
	resp.Resp.StatusCode = s
}

// ReadBody read body from the from the response
func (resp *Response) ReadBody() []byte {
	if resp.Resp != nil && resp.Resp.Body != nil {
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
	if resp.Resp != nil && resp.Resp.Body != nil {
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
