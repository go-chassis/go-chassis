package rest_test

import (
	"net/http"
	"testing"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/stretchr/testify/assert"
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

	req.SetHeader("a", "1")
	value := req.GetHeader("a")
	assert.Equal(t, value, "1")

	req.SetContentType("application/json")
	value = req.GetContentType()
	assert.Equal(t, value, "application/json")

	req.SetMethod("POST")
	method := req.GetMethod()
	assert.Equal(t, "POST", method)

	resp := rest.NewResponse()
	body := resp.ReadBody()
	assert.Empty(t, body)

	header := resp.GetHeader()
	assert.Empty(t, header)

	_ = resp.GetStatusCode()

	c1 := new(http.Cookie)
	c1.Name = "test"

	sessionIDValue := "abcdefg"
	c1.Value = sessionIDValue

	resp.SetCookie(c1)
	val := resp.GetCookie("test")
	assert.Equal(t, c1.Value, string(val))

	req, err = rest.NewRequest("GET", "cse://hello")
	assert.NoError(t, err)

	testHeaderKey := "hello"
	testHeaderValue := "ServiceComb"
	req.Req.Header.Add(testHeaderKey, testHeaderValue)
	req.Req.AddCookie(c1)
	newRequest := req.Copy()
	assert.Equal(t, newRequest.Req.Header.Get(testHeaderKey), testHeaderValue)
	assert.Equal(t, newRequest.GetCookie(c1.Name), c1.Value)

	req.Close()
	resp.Close()
}
