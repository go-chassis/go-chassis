package rest_test

import (
	"net/http"
	"testing"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/stretchr/testify/assert"
)

func TestNewRestRequest(t *testing.T) {
	t.Log("Testing all the restfull client functions")
	req, err := rest.NewRequest("GET", "cse://hello", nil)
	assert.NoError(t, err)

	req, err = rest.NewRequest("", "cse://hello", []byte("bodypart"))
	assert.NoError(t, err)

	httputil.SetURI(req, "cse://example/:id")
	uri := req.URL.String()
	assert.Equal(t, uri, "cse://example/:id")

	httputil.SetBody(req, []byte("hello"))

	req.Header.Set("a", "1")
	value := req.Header.Get("a")
	assert.Equal(t, value, "1")

	httputil.SetContentType(req, "application/json")
	value = httputil.GetContentType(req)
	assert.Equal(t, value, "application/json")

	req.Method = "POST"
	method := req.Method
	assert.Equal(t, "POST", method)

	resp := rest.NewResponse()
	body := httputil.ReadBody(resp)
	assert.Empty(t, body)

	header := resp.Header
	assert.Empty(t, header)

	c1 := new(http.Cookie)
	c1.Name = "test"

	sessionIDValue := "abcdefg"
	c1.Value = sessionIDValue

	httputil.SetRespCookie(resp, c1)
	val := httputil.GetRespCookie(resp, "test")
	assert.Equal(t, c1.Value, string(val))

	req, err = rest.NewRequest("GET", "cse://hello", nil)
	assert.NoError(t, err)

	testHeaderKey := "hello"
	testHeaderValue := "go-chassis"
	req.Header.Add(testHeaderKey, testHeaderValue)
	req.AddCookie(c1)
	assert.Equal(t, req.Header.Get(testHeaderKey), testHeaderValue)
	assert.Equal(t, httputil.GetCookie(req, c1.Name), c1.Value)

}
