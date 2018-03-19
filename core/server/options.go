package server

import (
	"crypto/tls"

	"github.com/ServiceComb/go-chassis/core/provider"
)

//Options is the options for service initiating
type Options struct {
	Address   string
	ChainName string
	Provider  provider.Provider
	TLSConfig *tls.Config
}

//RegisterOptions is options when you register a schema to chassis
type RegisterOptions struct {
	SchemaID string
}

//RegisterOption is option when you register a schema to chassis
type RegisterOption func(*RegisterOptions)

//WithSchemaID you can specify a unique id for schema
func WithSchemaID(schemaID string) RegisterOption {
	return func(o *RegisterOptions) {
		o.SchemaID = schemaID
	}
}
