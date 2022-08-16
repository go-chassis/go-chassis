// Package rate supply functionality about QPS
// for example rate limiting
package rate

import (
	"github.com/go-chassis/openlog"
	"sync"

	"k8s.io/client-go/util/flowcontrol"
)

// constant qps default rate
const (
	DefaultRate = 2147483647
)

// Limiters manages all rate limiters. it is thread safe and singleton.
// it create new limiters and try to limit request.
// each limiter has a unique name.
type Limiters struct {
	sync.RWMutex
	m map[string]flowcontrol.RateLimiter
}

// variables of qps limiter and mutex variable
var (
	once       = new(sync.Once)
	qpsLimiter *Limiters
)

// GetRateLimiters get qps rate limiters
func GetRateLimiters() *Limiters {
	once.Do(func() {
		qpsLimiter = &Limiters{m: make(map[string]flowcontrol.RateLimiter)}
	})
	return qpsLimiter
}

// TryAccept try to accept a request. if limiter can not accept a request, it returns false
// name is the limiter name
// qps is not necessary if the limiter already exists
func (qpsL *Limiters) TryAccept(name string, qps, burst int) bool {
	qpsL.RLock()
	limiter, ok := qpsL.m[name]
	if !ok {
		qpsL.RUnlock()
		//If the name operation is not present in the map, then add the new name operation to the map
		return qpsL.addLimiter(name, qps, burst)
	}
	qpsL.RUnlock()
	return limiter.TryAccept()
}

// addLimiter create a new limiter and add it to limiter map
func (qpsL *Limiters) addLimiter(name string, qps, burst int) bool {
	var bucketSize int
	// add a limiter object for the newly found operation in the Default Hash map
	// so that the default rate will be applied to subsequent token requests to this new operation
	if qps >= 1 {
		bucketSize = qps
	} else {
		bucketSize = DefaultRate
	}
	qpsL.Lock()
	// Create a new bucket for the new operation
	r := flowcontrol.NewTokenBucketRateLimiter(float32(bucketSize), burst)
	qpsL.m[name] = r
	qpsL.Unlock()
	return r.TryAccept()
}

// UpdateRateLimit will update the old limiters
func (qpsL *Limiters) UpdateRateLimit(name string, qps, burst int) {
	openlog.Info("add limiter", openlog.WithTags(openlog.Tags{
		"module": "rateLimiter",
		"event":  "update",
		"mark":   name,
		"qps":    qps,
		"burst":  burst,
	}))
	qpsL.addLimiter(name, qps, burst)
}

// DeleteRateLimiter delete rate limiter
func (qpsL *Limiters) DeleteRateLimiter(name string) {
	qpsL.Lock()
	delete(qpsL.m, name)
	qpsL.Unlock()
}
