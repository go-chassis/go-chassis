package client

import (
	"crypto/tls"
	"time"

	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	microTransport "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"golang.org/x/net/context"
)

type Options struct {
	// Used to select codec
	ContentType string

	// Plugged interfaces
	Codecs map[string]codec.Codec
	//TODO
	ClientCodecs map[string]codec.NewClientCodec
	Transport    microTransport.Transport

	// Connection Pool
	PoolSize int
	PoolTTL  time.Duration

	// Default Call Options
	CallOptions CallOptions

	// Other options for implementations of the interface
	// can be stored in a context
	Context   context.Context
	TLSConfig *tls.Config
	Failure   map[string]bool
}
type PublishOptions struct {
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}
type CallOptions struct {
	Message     []byte
	ContentType string
	// Transport Dial Timeout
	DialTimeout time.Duration
	// Number of Call attempts
	Retries int
	// Request/Response timeout
	RequestTimeout time.Duration

	// Other options for implementations of the interface
	// can be stored in a context
	Context    context.Context
	Header     map[string]string
	UrlPath    string
	MethodType string
}
type RequestOption func(*RequestOptions)
type RequestOptions struct {
	Stream bool
	// Transport Dial Timeout
	DialTimeout time.Duration
	// Request/Response timeout
	RequestTimeout time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

func newOptions(options ...Option) Options {
	opts := Options{
		CallOptions: CallOptions{
			Retries:        DefaultRetries,
			RequestTimeout: DefaultRequestTimeout,
			DialTimeout:    microTransport.DefaultDialTimeout,
		},
		PoolSize: DefaultPoolSize,
		PoolTTL:  DefaultPoolTTL,
	}

	for _, o := range options {
		o(&opts)
	}
	return opts
}

// Default content type of the client
func ContentType(ct string) Option {
	return func(o *Options) {
		o.ContentType = ct
	}
}

//WithCodecs set Codecs
func WithCodecs(c map[string]codec.Codec) Option {
	return func(o *Options) {
		o.Codecs = c
	}
}

// PoolSize sets the connection pool size
func PoolSize(d int) Option {
	return func(o *Options) {
		o.PoolSize = d
	}
}

// PoolSize sets the connection pool size
func PoolTTL(d time.Duration) Option {
	return func(o *Options) {
		o.PoolTTL = d
	}
}

// Transport to use for communication e.g http, rabbitmq, etc
func Transport(t microTransport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

// Number of retries when making the request.
// Should this be a Call Option?
func Retries(i int) Option {
	return func(o *Options) {
		o.CallOptions.Retries = i
	}
}

// The request timeout.
// Should this be a Call Option?
func RequestTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.RequestTimeout = d
	}
}

// Transport dial timeout
func DialTimeout(d time.Duration) Option {
	return func(o *Options) {
		o.CallOptions.DialTimeout = d
	}
}

// Call Options

// WithRetries is a CallOption which overrides that which
// set in Options.CallOptions
func WithRetries(i int) CallOption {
	return func(o *CallOptions) {
		o.Retries = i
	}
}

// WithRequestTimeout is a CallOption which overrides that which
// set in Options.CallOptions
func WithRequestTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.RequestTimeout = d
	}
}

// WithDialTimeout is a CallOption which overrides that which
// set in Options.CallOptions
func WithDialTimeout(d time.Duration) CallOption {
	return func(o *CallOptions) {
		o.DialTimeout = d
	}
}
func WithContentType(ct string) CallOption {
	return func(o *CallOptions) {
		o.ContentType = ct
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}
func WithUrlPath(s string) CallOption {
	return func(o *CallOptions) {
		o.UrlPath = s
	}
}

func WithMethodType(s string) CallOption {
	return func(o *CallOptions) {
		o.MethodType = s
	}
}

func WithConnectionPoolSize(n int) Option {
	return func(o *Options) {
		o.PoolSize = n
	}
}

func WithFailure(m map[string]bool) Option {
	return func(o *Options) {
		o.Failure = m
	}
}

// Option used by the Client
type Option func(*Options)

// CallOption used by Call or Stream
type CallOption func(*CallOptions)

var (
	// DefaultRetries is the default number of times a request is tried
	DefaultRetries = 1
	// DefaultRequestTimeout is the default request timeout
	DefaultRequestTimeout = time.Second * 5
	// DefaultPoolSize sets the connection pool size
	DefaultPoolSize = 50
	// DefaultPoolTTL sets the connection pool ttl
	DefaultPoolTTL = time.Minute
)
