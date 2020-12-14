package handler

import (
	"errors"
	"fmt"

	"github.com/go-chassis/foundation/stringutil"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
)

var errViolateBuildIn = errors.New("can not replace build-in handler func")

//ErrDuplicatedHandler means you registered more than 1 handler with same name
var ErrDuplicatedHandler = errors.New("duplicated handler registration")
var buildIn = []string{LoadBalancing, Router, TracingConsumer,
	TracingProvider, FaultInject}

// funcMap saves handler functions
var funcMap = make(map[string]func() Handler)

// constant keys for handlers
const (
	//consumer chain
	Transport       = "transport"
	LoadBalancing   = "loadbalance"
	Router          = "router"
	FaultInject     = "fault-inject"
	TracingConsumer = "tracing-consumer"
	TracingProvider = "tracing-provider"
)

// init is for to initialize the all handlers at boot time
func init() {
	//register build-in handler,don't need to call RegisterHandlerFunc
	funcMap[LoadBalancing] = newLBHandler
	funcMap[Router] = newRouterHandler
	funcMap[TracingProvider] = newTracingProviderHandler
	funcMap[TracingConsumer] = newTracingConsumerHandler
	funcMap[FaultInject] = newFaultHandler
	funcMap[TrafficMarker] = newMarkHandler
}

// Handler interface for handlers
type Handler interface {
	// handle invocation transportation,and tr response
	Handle(*Chain, *invocation.Invocation, invocation.ResponseCallBack)
	Name() string
}

//WriteBackErr write err and callback
func WriteBackErr(err error, status int, cb invocation.ResponseCallBack) {
	r := &invocation.Response{
		Err:    err,
		Status: status,
	}
	cb(r)
}

// RegisterHandler Let developer custom handler
func RegisterHandler(name string, f func() Handler) error {
	if stringutil.StringInSlice(name, buildIn) {
		return errViolateBuildIn
	}
	_, ok := funcMap[name]
	if ok {
		return ErrDuplicatedHandler
	}
	funcMap[name] = f
	return nil
}

// CreateHandler create a new handler by name your registered
func CreateHandler(name string) (Handler, error) {
	f := funcMap[name]
	if f == nil {
		return nil, fmt.Errorf("don't have handler [%s]", name)
	}
	return f(), nil
}
