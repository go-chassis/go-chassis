package httputil

import (
	"bytes"
	"errors"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"io/ioutil"
	"net/http"
)

//ErrInvalidReq invalid input
var ErrInvalidReq = errors.New("rest consumer call arg is not *http.Request type")

//SetURI sets host for the request.
func SetURI(req *http.Request, url string) {
	if tempURL, err := req.URL.Parse(url); err == nil {
		req.URL = tempURL
	}
}

//SetBody is a method used for setting body for a request
func SetBody(req *http.Request, body []byte) {
	req.Body = ioutil.NopCloser(bytes.NewReader(body))
}

//SetCookie set key value in request cookie
func SetCookie(req *http.Request, k, v string) {
	c := &http.Cookie{
		Name:  k,
		Value: v,
	}
	req.AddCookie(c)
}

//GetCookie is a method which gets cookie from a request
func GetCookie(req *http.Request, key string) string {
	cookie, err := req.Cookie(key)
	if err == http.ErrNoCookie {
		return ""
	}
	return cookie.Value
}

// SetContentType is a method used for setting content-type in a request
func SetContentType(req *http.Request, ct string) {
	req.Header.Set("Content-Type", ct)
}

// GetContentType is a method used for getting content-type in a request
func GetContentType(req *http.Request) string {
	return req.Header.Get("Content-Type")
}

//HTTPRequest convert invocation to http request
func HTTPRequest(inv *invocation.Invocation) (*http.Request, error) {
	reqSend, ok := inv.Args.(*http.Request)
	if !ok {
		return nil, ErrInvalidReq
	}
	m := common.FromContext(inv.Ctx)
	for k, v := range m {
		reqSend.Header.Set(k, v)
	}
	return reqSend, nil
}
