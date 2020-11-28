package client

import (
	"crypto/tls"
	"time"

	"context"
)

// Options is the list of dynamic parameter's which can be passed to the RegistryClient while creating a new client
type Options struct {
	Addrs     []string
	EnableSSL bool
	Timeout   time.Duration
	TLSConfig *tls.Config
	// Other options can be stored in a context
	Context    context.Context
	Compressed bool
	Verbose    bool
	Version    string
}

//CallOptions is options when you call a API
type CallOptions struct {
	WithoutRevision bool
	Revision        string
	WithGlobal      bool
}

//WithoutRevision ignore current revision number
func WithoutRevision() CallOption {
	return func(o *CallOptions) {
		o.WithoutRevision = true
	}
}

//WithGlobal query resources include other aggregated SC
func WithGlobal() CallOption {
	return func(o *CallOptions) {
		o.WithGlobal = true
	}
}

//CallOption is receiver for options and chang the attribute of it
type CallOption func(*CallOptions)
