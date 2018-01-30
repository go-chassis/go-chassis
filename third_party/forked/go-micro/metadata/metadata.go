// Package metadata is a way of defining message headers
package metadata

import (
	context17 "context"
	"golang.org/x/net/context"
)

type metaKey struct{}

// Metadata is used to represent request headers internally.
// It is used at the RPC level and translate back and forth
// from request headers.
type Metadata map[string]string

// FromContext get the context contents and returns meta data
func FromContext(ctx context.Context) (Metadata, bool) {
	if ctx == nil {
		return nil, false
	}
	m, ok := ctx.Value(metaKey{}).(Metadata)
	return m, ok
}

// FromContext17 get the context contents and returns meta data
func FromContext17(ctx context17.Context) (Metadata, bool) {
	if ctx == nil {
		return nil, false
	}
	m, ok := ctx.Value(metaKey{}).(Metadata)
	return m, ok
}

// NewContext returns the context object
func NewContext(ctx context.Context, md Metadata) context.Context {
	return context.WithValue(ctx, metaKey{}, md)
}

// NewContext17 returns the context object
func NewContext17(ctx context17.Context, md Metadata) context.Context {
	return context17.WithValue(ctx, metaKey{}, md)
}
