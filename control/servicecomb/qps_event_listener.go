package servicecomb

import (
	"fmt"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/v2/resilience/rate"
	"github.com/go-chassis/openlog"

	"strings"

	"github.com/go-chassis/go-chassis/v2/core/common"
)

//QPSEventListener is a struct used for Event listener
type QPSEventListener struct {
	//Key []string
	Key string
}

//Event is a method for QPS event listening
func (el *QPSEventListener) Event(e *event.Event) {
	qpsLimiter := rate.GetRateLimiters()

	if strings.Contains(e.Key, "enabled") {
		return
	}
	qps, ok := e.Value.(int)
	if !ok {
		openlog.Error(fmt.Sprintf("invalid qps config %s", e.Value))
	}
	openlog.Info("update rate limiter", openlog.WithTags(openlog.Tags{
		"module": "RateLimiting",
		"event":  e.EventType,
		"value":  qps,
	}))
	switch e.EventType {
	case common.Update:
		qpsLimiter.UpdateRateLimit(e.Key, qps, qps/5)
	case common.Create:
		qpsLimiter.UpdateRateLimit(e.Key, qps, qps/5)
	case common.Delete:
		qpsLimiter.DeleteRateLimiter(e.Key)
	}
}
