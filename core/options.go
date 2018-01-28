package core

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	"time"
)

// Options is a struct to stores information about chain name, filters, and their invocation options
type Options struct {
	// chain for client
	ChainName         string
	Filters           []selector.Filter
	InvocationOptions InvokeOptions
}

// InvokeOptions struct having information about microservice API call parameters
type InvokeOptions struct {
	//microservice version
	Version string
	Stream  bool
	// Transport Dial Timeout
	DialTimeout time.Duration
	// Request/Response timeout
	RequestTimeout time.Duration
	// end to end，Directly call
	Endpoint string
	// end to end，Directly call
	Protocol string
	// end to end，Directly call
	ContentType string
	//loadbalance stratery
	//StrategyFunc loadbalance.Strategy
	StrategyFunc string
	Filters      []selector.Filter
	URLPath      string
	MethodType   string
	AppID        string
	// local data
	Metadata map[string]interface{}
}

//TODO a lot of options

// DefaultOptions for chain name
var DefaultOptions = Options{
	ChainName: "default",
}

// ChainName is request option
func ChainName(name string) Option {
	return func(o *Options) {
		o.ChainName = name
	}
}

// Filters is request option
func Filters(f []selector.Filter) Option {
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

// WithVersion is a request option
func WithVersion(v string) InvocationOption {
	return func(o *InvokeOptions) {
		o.Version = v
	}
}

// WithProtocol is a request option
func WithProtocol(p string) InvocationOption {
	return func(o *InvokeOptions) {
		o.Protocol = p
	}
}

// WithAppID is a request option
func WithAppID(p string) InvocationOption {
	return func(o *InvokeOptions) {
		o.AppID = p
	}
}

// WithContentType is a request option
func WithContentType(ct string) InvocationOption {
	return func(o *InvokeOptions) {
		o.ContentType = ct
	}
}

//Request Options
/*func WithStrategy(s loadbalance.Strategy) InvocationOption {
	return func(o *InvokeOptions) {
		o.StrategyFunc = s
	}
}*/

// WithStrategy is a request option
func WithStrategy(s string) InvocationOption {
	return func(o *InvokeOptions) {
		o.StrategyFunc = s
	}
}

// WithFilters is a request option
func WithFilters(f ...selector.Filter) InvocationOption {
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

// getOpts is to get the options
func getOpts(microservice string, options ...InvocationOption) InvokeOptions {
	opts := InvokeOptions{}
	for _, o := range options {
		o(&opts)
	}
	if opts.ContentType == "" {
		opts.ContentType = common.JSON
	}
	if opts.Version == "" {
		opts.Version = config.GlobalDefinition.Cse.References[microservice].Version
		if opts.Version == "" {
			opts.Version = common.DefaultVersion
		}
	}
	return opts
}

// wrapInvocationWithOpts is wrap invocation with options
func wrapInvocationWithOpts(i *invocation.Invocation, opts InvokeOptions) {
	i.Endpoint = opts.Endpoint
	i.Protocol = opts.Protocol
	i.ContentType = opts.ContentType
	i.Version = opts.Version
	i.Strategy = opts.StrategyFunc
	i.Filters = opts.Filters
	i.AppID = opts.AppID
	i.Metadata = opts.Metadata
}
