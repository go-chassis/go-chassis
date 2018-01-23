package eventlistener

import (
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/qpslimiter"
	"strings"
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
func (e *QPSEventListener) Event(event *core.Event) {
	qpsLimiter := qpslimiter.GetQPSTrafficLimiter()

	if strings.Contains(event.Key, "enabled") {
		return
	}

	switch event.EventType {
	case common.Update:
		qpsLimiter.UpdateRateLimit(event.Key, event.Value)
	case common.Create:
		qpsLimiter.UpdateRateLimit(event.Key, event.Value)
	case common.Delete:
		qpsLimiter.DeleteRateLimiter(event.Key)
	}
}
