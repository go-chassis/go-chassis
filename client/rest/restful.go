package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

// NewRequest is a function which creates new request
func NewRequest(method, urlStr string, body []byte) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		r = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(context.Background(), method, urlStr, r)
	if err != nil {
		return nil, err
	}
	return req, nil
}

// NewResponse is creating the object of response
func NewResponse() *http.Response {
	resp := &http.Response{
		Header: http.Header{},
	}
	return resp
}
