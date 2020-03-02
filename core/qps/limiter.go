//Package qps supply functionality about QPS
//for example rate limiting
package qps

import (
	"strconv"
	"sync"

	"github.com/go-chassis/go-archaius"
	"github.com/go-mesh/openlogging"
	"k8s.io/client-go/util/flowcontrol"
)

// constant qps default rate
const (
	DefaultRate = 2147483647
)

// RateLimiters qps limiter map struct
type RateLimiters struct {
	sync.RWMutex
	m map[string]flowcontrol.RateLimiter
}

// variables of qps limiter and mutex variable
var (
	once       = new(sync.Once)
	qpsLimiter *RateLimiters
)

// GetRateLimiters get qps rate limiters
func GetRateLimiters() *RateLimiters {
	once.Do(func() {
		qpsLimiter = &RateLimiters{m: make(map[string]flowcontrol.RateLimiter)}
	})
	return qpsLimiter
}

// TryAccept process qps token request
func (qpsL *RateLimiters) TryAccept(key string, qpsRate int) bool {
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

// addLimiter process default rate pps token request
func (qpsL *RateLimiters) addLimiter(key string, qpsRate int) bool {
	var bucketSize int
	// add a limiter object for the newly found operation in the Default Hash map
	// so that the default rate will be applied to subsequent token requests to this new operation
	if qpsRate >= 1 {
		bucketSize = qpsRate
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

// GetQPSRate get qps rate
func GetQPSRate(rateConfig string) (int, bool) {
	qpsRate := archaius.GetInt(rateConfig, DefaultRate)
	if qpsRate == DefaultRate {
		return qpsRate, false
	}

	return qpsRate, true
}

// GetQPSRateWithPriority get qps rate with priority
func (qpsL *RateLimiters) GetQPSRateWithPriority(cmd ...string) (int, string) {
	var (
		qpsVal      int
		configExist bool
	)
	for _, c := range cmd {
		qpsVal, configExist = GetQPSRate(c)
		if configExist {
			return qpsVal, c
		}
	}

	return DefaultRate, cmd[len(cmd)-1]

}

// UpdateRateLimit update or add rate limiter
func (qpsL *RateLimiters) UpdateRateLimit(key string, value interface{}) {
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
func (qpsL *RateLimiters) DeleteRateLimiter(key string) {
	qpsL.Lock()
	delete(qpsL.m, key)
	qpsL.Unlock()
}
