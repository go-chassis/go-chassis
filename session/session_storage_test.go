package session_test

import (
	"context"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/session"
	"github.com/stretchr/testify/assert"
	"net/http"
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
	cookieValue := session.GetSessionCookie(nil, nil)
	assert.Equal(t, "", cookieValue)

	resp := &http.Response{
		Header: http.Header{},
	}
	cookieValue = session.GetSessionCookie(nil, resp)
	assert.Equal(t, "", cookieValue)
	session.DeletingKeySuccessiveFailure(nil)
	session.DeletingKeySuccessiveFailure(resp)
	cookieValue = session.GetSessionFromResp("abc", resp)
	assert.Equal(t, "", cookieValue)
	session.CheckForSessionID("", 1, resp, new(http.Request))

	ctx := context.Background()
	ctx = session.SetContextMetadata(ctx, "key", "value")
	val := session.GetContextMetadata(ctx, "key")
	assert.Equal(t, val, "value")
	ctx = session.CheckForSessionIDFromContext(ctx, "", 1)
	val = session.GetContextMetadata(ctx, "ServiceCombLB")
	assert.NotEqual(t, val, "")
}
