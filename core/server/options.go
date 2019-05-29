package server

import (
	"crypto/tls"

	"github.com/go-chassis/go-chassis/core/provider"
)

//Options is the options for service initiating
type Options struct {
	Address            string
	ProtocolServerName string
	ChainName          string
	Provider           provider.Provider
	TLSConfig          *tls.Config
	BodyLimit          int64
}

//RegisterOptions is options when you register a schema to chassis
type RegisterOptions struct {
	SchemaID   string
	RPCSvcDesc interface{}
}

//RegisterOption is option when you register a schema to chassis
type RegisterOption func(*RegisterOptions)

//WithSchemaID you can specify a unique id for schema
func WithSchemaID(schemaID string) RegisterOption {
	return func(o *RegisterOptions) {
		o.SchemaID = schemaID
	}
}

//WithRPCServiceDesc you can set rpc service desc, it cloud be *grpc.ServiceDesc
func WithRPCServiceDesc(RPCSvcDesc interface{}) RegisterOption {
	return func(o *RegisterOptions) {
		o.RPCSvcDesc = RPCSvcDesc
	}
}
