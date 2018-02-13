package session_test

import (
	"github.com/ServiceComb/go-chassis/session"
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
}
