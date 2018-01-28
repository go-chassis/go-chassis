package selector

import (
	"github.com/ServiceComb/go-chassis/core/registry"
	"golang.org/x/net/context"
)

// Options is having registry, strategy, context variables
type Options struct {
	Registry registry.Registry
	Strategy Strategy

	// Other options can be stored in a context
	Context context.Context
}

// SelectOptions is having micro-service filters, strategy, appid, consumerid, context
type SelectOptions struct {
	Filters    []Filter
	Strategy   Strategy
	AppID      string
	Metadata   interface{}
	ConsumerID string
	// Other options can be stored in a context
	Context context.Context
}

// Option used to initialise Options struct
type Option func(*Options)

// SelectOption used to initialise SelectOptions struct
type SelectOption func(*SelectOptions)

// Registry sets the registry that selector used
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// SetStrategy sets the default strategy that selector used
func SetStrategy(fn Strategy) Option {
	return func(o *Options) {
		o.Strategy = fn
	}
}

// WithFilter adds a filter func to the list of filters
func WithFilter(fns []Filter) SelectOption {
	return func(o *SelectOptions) {
		o.Filters = append(o.Filters, fns...)
	}
}

// WithAppID adds a application id to the function
func WithAppID(a string) SelectOption {
	return func(o *SelectOptions) {
		o.AppID = a
	}
}

// WithConsumerID consumer id is added to the selectoptions
func WithConsumerID(id string) SelectOption {
	return func(o *SelectOptions) {
		o.ConsumerID = id
	}
}

// WithStrategy sets the selector strategy
func WithStrategy(fn Strategy) SelectOption {
	return func(o *SelectOptions) {
		o.Strategy = fn
	}
}

// WithMetadata sets the selector metadata
func WithMetadata(a interface{}) SelectOption {
	return func(o *SelectOptions) {
		o.Metadata = a
	}
}
