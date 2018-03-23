package session_test

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/session"
	"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSessionStorage(t *testing.T) {
	session.Save("abc", "127.0.0.1:8080", time.Second)
	addr, ok := session.Get("abc")
	assert.Equal(t, true, ok)
	assert.Equal(t, "127.0.0.1:8080", addr)
	session.Delete("abc")
	addr, ok = session.Get("abc")
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, addr)
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	cookieValue := session.GetSessionCookie(nil)
	assert.Equal(t, "", cookieValue)
	var resp *fasthttp.Response
	resp = new(fasthttp.Response)
	cookieValue = session.GetSessionCookie(resp)
	assert.Equal(t, "", cookieValue)
	session.DeletingKeySuccessiveFailure(nil)
	session.DeletingKeySuccessiveFailure(resp)
	cookieValue = session.GetSessionFromResp("abc", resp)
	assert.Equal(t, "", cookieValue)
	session.CheckForSessionID("", 1, resp, new(fasthttp.Request))
}
