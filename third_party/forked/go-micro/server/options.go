package server

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"crypto/tls"
	"time"

	"github.com/ServiceComb/go-chassis/core/provider"
	"golang.org/x/net/context"
)

type Options struct {
	Metadata map[string]string
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
