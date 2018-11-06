package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdd_GetSessionStickinessCache(t *testing.T) {
	SessionStickinessCache = initCache()
	testSlice := []string{"services-cookies1", "services-cookies2"}

	AddSessionStickinessToCache(testSlice[0], "")
	cookie := GetSessionID("")
	assert.Equal(t, "services-cookies1", cookie)

	AddSessionStickinessToCache(testSlice[1], "")
	cookie = GetSessionID("")
	assert.Equal(t, "services-cookies2", cookie)
}
