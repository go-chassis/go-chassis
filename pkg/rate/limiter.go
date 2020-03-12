//Package rate supply functionality about QPS
//for example rate limiting
package rate

import (
	"strconv"
	"sync"

	"github.com/go-mesh/openlogging"
	"k8s.io/client-go/util/flowcontrol"
)

// constant qps default rate
const (
	DefaultRate = 2147483647
)

//Limiters manages all rate limiter
//it create new limiters and try to limit processes
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

//TryAccept process request, if it can not process a request, it returns false
func (qpsL *Limiters) TryAccept(key string, qpsRate int) bool {
	qpsL.RLock()
	limiter, ok := qpsL.m[key]
	if !ok {
		qpsL.RUnlock()
		//If the key operation is not present in the map, then add the new key operation to the map
		return qpsL.addLimiter(key, qpsRate)
	}
	qpsL.RUnlock()
	return limiter.TryAccept()

}

// addLimiter process create a new limiter and add it to limiter map
func (qpsL *Limiters) addLimiter(key string, qps int) bool {
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
	r := flowcontrol.NewTokenBucketRateLimiter(float32(bucketSize), 1)
	qpsL.m[key] = r
	qpsL.Unlock()
	return r.TryAccept()
}

// UpdateRateLimit will update the old limiters
func (qpsL *Limiters) UpdateRateLimit(key string, value interface{}) {
	switch v := value.(type) {
	case int:
		qpsL.addLimiter(key, value.(int))
	case string:
		convertedIntValue, err := strconv.Atoi(value.(string))
		if err != nil {
			openlogging.GetLogger().Warnf("invalid value type received for rate limiter: %v", v, err)
		} else {
			qpsL.addLimiter(key, convertedIntValue)
		}
	default:
		openlogging.GetLogger().Warnf("invalid value type received for rate limiter: %v", v)
	}
}

// DeleteRateLimiter delete rate limiter
func (qpsL *Limiters) DeleteRateLimiter(key string) {
	qpsL.Lock()
	delete(qpsL.m, key)
	qpsL.Unlock()
}
