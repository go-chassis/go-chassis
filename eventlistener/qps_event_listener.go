package eventlistener

import (
	"fmt"
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/resilience/rate"
	"github.com/go-mesh/openlogging"

	"strings"

	"github.com/go-chassis/go-chassis/core/common"
)

const (
	//QPSLimitKey is a variable of type string
	QPSLimitKey = "servicecomb.flowcontrol"
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
		openlogging.Error(fmt.Sprintf("invalid qps config %s", e.Value))
	}
	openlogging.Info("update rate limiter", openlogging.WithTags(openlogging.Tags{
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
