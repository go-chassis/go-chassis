package eventlistener

import (
	"github.com/ServiceComb/go-archaius/core"
	"github.com/ServiceComb/go-chassis/core/qpslimiter"
	"strings"
)

const (
	//QPSLimitKey is a variable of type string
	QPSLimitKey = "cse.flowcontrol"
)

//QpsEventListener is a struct used for Event listener
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
	case "UPDATE":
		qpsLimiter.UpdateRateLimit(event.Key, event.Value)
	case "CREATE":
		qpsLimiter.UpdateRateLimit(event.Key, event.Value)
	case "DELETE":
		qpsLimiter.DeleteRateLimiter(event.Key)
	}
}
