package server

import (
	"crypto/tls"
	"time"

	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"golang.org/x/net/context"
)

type Options struct {
	Codecs    map[string]codec.Codec
	Transport transport.Transport
	Metadata  map[string]string
	//protocol
	Name      string
	Address   string
	Advertise string
	ID        string
	Version   string
	//singleton
	ChainName string
	//each microservice is a provider
	Provider    provider.Provider
	RegisterTTL time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context   context.Context
	TLSConfig *tls.Config
}

type Option func(*Options)

type RegisterOption func(*RegisterOptions)

var DefaultOptions = Options{
	ChainName: "default",
	Metadata:  map[string]string{},
}

//WithCodecs set Codecs
func WithCodecs(c map[string]codec.Codec) Option {
	return func(o *Options) {
		o.Codecs = c
	}
}

// Server name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Unique server id
func ID(id string) Option {
	return func(o *Options) {
		o.ID = id
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// The address to advertise for discovery - host:port
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

//// Registry used for discovery
//func Registry(r api.Register) Option {
//	return func(o *Options) {
//		o.Registry = r
//	}
//}

// Transport mechanism for communication e.g http, rabbitmq, etc
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

func ChainName(n string) Option {
	return func(o *Options) {
		o.ChainName = n
	}
}
func Provider(p provider.Provider) Option {
	return func(o *Options) {
		o.Provider = p
	}
}

// Metadata associated with the server
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// Register the service with a TTL
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

type RegisterOptions struct {
	MicroServiceName string
	SchemaID         string
	Provider         provider.Provider
	//only grpc protocol use the param
	GrpcRegister interface{}
}

func WithMicroServiceName(microservice string) RegisterOption {
	return func(o *RegisterOptions) {
		o.MicroServiceName = microservice
	}
}
func WithSchemaID(schemaID string) RegisterOption {
	return func(o *RegisterOptions) {
		o.SchemaID = schemaID
	}
}

func WithServiceProvider(provider provider.Provider) RegisterOption {
	return func(o *RegisterOptions) {
		o.Provider = provider
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(o *Options) {
		o.TLSConfig = t
	}
}

//WithGrpcRegister set GrpcRegister
func WithGrpcRegister(grpcRegister interface{}) RegisterOption {
	return func(o *RegisterOptions) {
		o.GrpcRegister = grpcRegister
	}
}
