package loadbalance

import (
	"time"
)

// Save for setting the service id, address, timeout
func Save(sid string, addr string, timeOut time.Duration) {
	SessionCache.Set(sid, addr, timeOut)
}

// Delete for deleting the sid
func Delete(sid string) {
	SessionCache.Delete(sid)
}
