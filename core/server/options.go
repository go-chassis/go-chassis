package server

import (
	"crypto/tls"

	"github.com/go-chassis/go-chassis/core/provider"
	"google.golang.org/grpc"
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
	SchemaID    string
	GRPCSvcDesc *grpc.ServiceDesc
}

//RegisterOption is option when you register a schema to chassis
type RegisterOption func(*RegisterOptions)

//WithSchemaID you can specify a unique id for schema
func WithSchemaID(schemaID string) RegisterOption {
	return func(o *RegisterOptions) {
		o.SchemaID = schemaID
	}
}

//WithGRPCServiceDesc you can set grpc service desc
func WithGRPCServiceDesc(GRPCSvcDesc *grpc.ServiceDesc) RegisterOption {
	return func(o *RegisterOptions) {
		o.GRPCSvcDesc = GRPCSvcDesc
	}
}
