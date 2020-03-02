package eventlistener

import (
	"github.com/go-chassis/go-archaius/event"
	"strings"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/qps"
)

const (
	//QPSLimitKey is a variable of type string
	QPSLimitKey = "cse.flowcontrol"
)

//QPSEventListener is a struct used for Event listener
type QPSEventListener struct {
	//Key []string
	Key string
}

//Event is a method for QPS event listening
func (el *QPSEventListener) Event(e *event.Event) {
	qpsLimiter := qps.GetRateLimiters()

	if strings.Contains(e.Key, "enabled") {
		return
	}
	//TODO watch new config
	switch e.EventType {
	case common.Update:
		qpsLimiter.UpdateRateLimit(e.Key, e.Value)
	case common.Create:
		qpsLimiter.UpdateRateLimit(e.Key, e.Value)
	case common.Delete:
		qpsLimiter.DeleteRateLimiter(e.Key)
	}
}
