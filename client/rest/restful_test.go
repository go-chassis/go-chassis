package rest_test

import (
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
	"github.com/stretchr/testify/assert"
	"testing"
)

var req *rest.Request

func TestNewRestRequest(t *testing.T) {
	t.Log("Testing all the restfull client functions")
	req, err := rest.NewRequest("GET", "cse://hello")
	assert.NoError(t, err)

	req, err = rest.NewRequest("", "cse://hello", []byte("bodypart"))
	assert.NoError(t, err)

	req.SetURI("cse://example/:id")
	uri := req.GetURI()
	assert.Equal(t, uri, "cse://example/:id")

	req.SetBody([]byte("hello"))
	req.SetHeader("Content-Type", "application/json")
	value := req.GetHeader("Content-Type")
	assert.Equal(t, value, "application/json")

	req.SetMethod("POST")
	method := req.GetMethod()
	assert.Equal(t, "POST", method)

	resp := rest.NewResponse()
	body := resp.ReadBody()
	assert.Empty(t, body)

	header := resp.GetHeader()
	assert.NotEmpty(t, header)

	_ = resp.GetStatusCode()

	var c1 *fasthttp.Cookie
	c1 = new(fasthttp.Cookie)
	c1.SetKey("test")

	sessionIDValue := "abcdefg"
	c1.SetValue(sessionIDValue)

	resp.SetCookie(c1)
	val := resp.GetCookie("test")
	assert.Equal(t, c1.Cookie(), val)

	req, err = rest.NewRequest("GET", "cse://hello")
	assert.NoError(t, err)
	_ = req.Copy()
	req.Close()
	resp.Close()
}
