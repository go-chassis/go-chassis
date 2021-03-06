package servicecomb

import (
	"github.com/go-chassis/go-archaius/event"
	"github.com/go-chassis/go-chassis/v2/control"
	"github.com/go-chassis/go-chassis/v2/resilience/rate"
	"github.com/go-chassis/openlog"
	"strconv"

	"fmt"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"strings"
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
	if e.Value == nil {
		openlog.Error(fmt.Sprintf("nil qps value %s", e.Key))
		return
	}
	qps, ok := e.Value.(int)
	if !ok {
		var err error
		qpsString, ok := e.Value.(string)
		if !ok {
			openlog.Error(fmt.Sprintf("invalid qps config %s", e.Value))
			return
		}
		qps, err = strconv.Atoi(qpsString)
		if err != nil {
			openlog.Error(fmt.Sprintf("invalid qps config %s", e.Value))
			return
		}
	}
	openlog.Info("update rate limiter", openlog.WithTags(openlog.Tags{
		"module": "RateLimiting",
		"event":  e.EventType,
		"value":  qps,
	}))
	burst := qps / 5
	if burst == 0 {
		burst = control.DefaultBurst
	}
	switch e.EventType {
	case common.Update:
		qpsLimiter.UpdateRateLimit(e.Key, qps, burst)
	case common.Create:
		qpsLimiter.UpdateRateLimit(e.Key, qps, burst)
	case common.Delete:
		qpsLimiter.DeleteRateLimiter(e.Key)
	}
}
