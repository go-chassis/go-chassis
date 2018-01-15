package handler

import (
	"errors"
	"fmt"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/util/string"
)

var errViolateBuildIn = errors.New("Can not replace build-in handler func")
var buildIn = []string{BizkeeperConsumer, BizkeeperProvider, Loadbalance, Router, TracingConsumer,
	TracingProvider, RatelimiterConsumer, RatelimiterProvider, Transport, FaultInject}

// HandlerFuncMap handler function map
var HandlerFuncMap = make(map[string]func() Handler)

// constant keys for handlers
const (
	Transport           = "transport"
	Loadbalance         = "loadbalance"
	BizkeeperConsumer   = "bizkeeper-consumer"
	BizkeeperProvider   = "bizkeeper-provider"
	TracingConsumer     = "tracing-consumer"
	TracingProvider     = "tracing-provider"
	RatelimiterConsumer = "ratelimiter-consumer"
	RatelimiterProvider = "ratelimiter-provider"
	Router              = "router"
	FaultInject         = "fault-inject"
)

// init is for to initialize the all handlers at boot time
func init() {
	//register build-in handler,don't need to call RegisterHandlerFunc
	HandlerFuncMap[Transport] = newTransportHandler
	HandlerFuncMap[Loadbalance] = newLBHandler
	HandlerFuncMap[BizkeeperConsumer] = newBizKeeperConsumerHandler
	HandlerFuncMap[BizkeeperProvider] = newBizKeeperProviderHandler
	HandlerFuncMap[RatelimiterConsumer] = newConsumerRateLimiterHandler
	HandlerFuncMap[RatelimiterProvider] = newProviderRateLimiterHandler
	HandlerFuncMap[TracingProvider] = newTracingProviderHandler
	HandlerFuncMap[TracingConsumer] = newTracingConsumerHandler
	HandlerFuncMap[Router] = newRouterHandler
	HandlerFuncMap[FaultInject] = FaultHandle
}

// Handler interface for handlers
type Handler interface {
	// handle invocation transportation,and tr response
	Handle(*Chain, *invocation.Invocation, invocation.ResponseCallBack)
	Name() string
}

func writeErr(err error, cb invocation.ResponseCallBack) {
	r := &invocation.InvocationResponse{
		Err: err,
	}
	cb(r)
}

// RegisterHandler Let developer custom handler
func RegisterHandler(name string, f func() Handler) error {
	if stringutil.StringInSlice(name, buildIn) {
		return errViolateBuildIn
	}
	HandlerFuncMap[name] = f
	return nil
}

// CreateHandler create a new handler by name your registered
func CreateHandler(name string) (Handler, error) {
	f := HandlerFuncMap[name]
	if f == nil {
		return nil, fmt.Errorf("Don't have handler [%s]", name)
	}
	return f(), nil
}
