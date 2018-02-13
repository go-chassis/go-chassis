package session

import (
	"time"
)

// Save for setting the service id, address, timeout
func Save(sid string, addr string, timeOut time.Duration) {
	SessionCache.Set(sid, addr, timeOut)
}

// Get return session id based on session key
func Get(k string) (sid interface{}, ok bool) {
	sid, ok = SessionCache.Get(k)
	return
}

//ClearExpired delete all expired session
func ClearExpired() {
	SessionCache.DeleteExpired()
}

// Delete for deleting the sid
func Delete(key string) {
	SessionCache.Delete(key)
}
