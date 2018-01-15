package registry

import (
	"crypto/tls"
	"time"

	"golang.org/x/net/context"
)

// Options having micro-service parameters
type Options struct {
	Addrs        []string
	EnableSSL    bool
	ConfigTenant string
	Timeout      time.Duration
	TLSConfig    *tls.Config
	// Other options can be stored in a context
	Context    context.Context
	Compressed bool
	Verbose    bool
	Version    string
}

// Option is the function of the type *Options
type Option func(*Options)

// Compressed set the Compressed parameter
func Compressed(bo bool) Option {
	return func(o *Options) {
		o.Compressed = bo
	}
}

// Verbose sets the Verbose parameter
func Verbose(bo bool) Option {
	return func(o *Options) {
		o.Verbose = bo
	}
}

// Version sets the version
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// EnableSSL enable the ssl
func EnableSSL(bo bool) Option {
	return func(o *Options) {
		o.EnableSSL = bo
	}
}

// Addrs sets the address
func Addrs(addrs ...string) Option {
	return func(o *Options) {
		o.Addrs = addrs
	}
}

// Timeout sets the timeout
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// TLSConfig sets tls configurations
func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

// Tenant sets tenant parameter
func Tenant(str string) Option {
	return func(o *Options) {
		o.ConfigTenant = str
	}
}
