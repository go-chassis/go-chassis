package server

import (
	"crypto/tls"
	"time"

	"github.com/go-chassis/go-chassis/v2/core/provider"
)

//Options is the options for service initiating
type Options struct {
	Address            string
	ProtocolServerName string
	ChainName          string
	Provider           provider.Provider
	TLSConfig          *tls.Config
	BodyLimit          int64
	HeaderLimit        int
	Timeout            time.Duration

	ProfilingEnable bool
	ProfilingAPI    string

	MetricsEnable bool
	MetricsAPI    string
}

//RegisterOptions is options when you register a schema to chassis
type RegisterOptions struct {
	SchemaID   string
	Method     string
	Path       string
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

//WithPath specify a url pattern
func WithPath(Path string) RegisterOption {
	return func(o *RegisterOptions) {
		o.Path = Path
	}
}

//WithPah specify a method
func WithMethod(Method string) RegisterOption {
	return func(o *RegisterOptions) {
		o.Method = Method
	}
}

//WithRPCServiceDesc you can set rpc service desc, it cloud be *grpc.ServiceDesc
func WithRPCServiceDesc(RPCSvcDesc interface{}) RegisterOption {
	return func(o *RegisterOptions) {
		o.RPCSvcDesc = RPCSvcDesc
	}
}
