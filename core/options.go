package core

import (
	"time"

	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/pkg/util/tags"
)

// Options is a struct to stores information about chain name, filters, and their invocation options
type Options struct {
	// chain for client
	ChainName         string
	Filters           []string
	InvocationOptions InvokeOptions
}

// InvokeOptions struct having information about microservice API call parameters
type InvokeOptions struct {
	Stream bool
	// Transport Dial Timeout
	DialTimeout time.Duration
	// Request/Response timeout
	RequestTimeout time.Duration
	// end to end, Directly call
	Endpoint string
	// end to end, Directly call
	Protocol string
	Port     string
	//loadbalancer stratery
	//StrategyFunc loadbalancer.Strategy
	StrategyFunc string
	Filters      []string
	URLPath      string
	MethodType   string
	// local data
	Metadata map[string]interface{}
	// tags for router
	RouteTags utiltags.Tags
}

//TODO a lot of options

// DefaultOptions for chain name
var DefaultOptions = Options{
	ChainName: "default",
}

//ChainName is able to custom handler chain for a invoker.
//you can specify a handler chain name under "servicecomb.handler.chain.Consumer" in chassis.yaml file.
//so that you can define different invoker with different handler chain.
//a handler chain is bind to a invoker instance.
func ChainName(name string) Option {
	return func(o *Options) {
		o.ChainName = name
	}
}

// Filters is request option
func Filters(f []string) Option {
	return func(o *Options) {
		o.Filters = f
	}
}

// DefaultCallOptions is request option
func DefaultCallOptions(io InvokeOptions) Option {
	return func(o *Options) {
		o.InvocationOptions = io
	}
}

// Option used by the invoker
type Option func(*Options)

// InvocationOption is a requestOption used by invocation
type InvocationOption func(*InvokeOptions)

// StreamingRequest is request option
func StreamingRequest() InvocationOption {
	return func(o *InvokeOptions) {
		o.Stream = true
	}
}

// WithEndpoint is request option
func WithEndpoint(ep string) InvocationOption {
	return func(o *InvokeOptions) {
		o.Endpoint = ep
	}
}

// WithProtocol is a request option
func WithProtocol(p string) InvocationOption {
	return func(o *InvokeOptions) {
		o.Protocol = p
	}
}

// WithStrategy is a request option
func WithStrategy(s string) InvocationOption {
	return func(o *InvokeOptions) {
		o.StrategyFunc = s
	}
}

// WithFilters is a request option
func WithFilters(f ...string) InvocationOption {
	return func(o *InvokeOptions) {
		o.Filters = append(o.Filters, f...)
	}
}

// WithMetadata is a request option
func WithMetadata(h map[string]interface{}) InvocationOption {
	return func(o *InvokeOptions) {
		o.Metadata = h
	}
}

// WithRouteTags is a request option
func WithRouteTags(t map[string]string) InvocationOption {
	return func(o *InvokeOptions) {
		o.RouteTags.Label = utiltags.LabelOfTags(t)
		o.RouteTags.KV = t
	}
}

// getOpts is to get the options
func getOpts(microservice string, options ...InvocationOption) InvokeOptions {
	opts := InvokeOptions{}
	for _, o := range options {
		o(&opts)
	}
	return opts
}

// wrapInvocationWithOpts is wrap invocation with options
func wrapInvocationWithOpts(i *invocation.Invocation, opts InvokeOptions) {
	i.Endpoint = opts.Endpoint
	i.Protocol = opts.Protocol
	i.Strategy = opts.StrategyFunc
	i.Filters = opts.Filters
	i.PortName = opts.Port
	if opts.Metadata != nil {
		i.Metadata = opts.Metadata
	}

	i.RouteTags = opts.RouteTags
}
