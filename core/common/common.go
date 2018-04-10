package common

import "context"

// constant for provider and consumer
const (
	Provider = "Provider"
	Consumer = "Consumer"
)

// constant for transport tcp
const (
	TransportTCP = "tcp"
)

// constant for microservice environment parameters
const (
	Env = "ServiceComb_ENV"

	EnvNodeIP     = "HOSTING_SERVER_IP"
	EnvSchemaRoot = "SCHEMA_ROOT"
	EnvProjectID  = "CSE_PROJECT_ID"
)

// constant environment keys service center, config center, monitor server addresses
const (
	CseRegistryAddress     = "CSE_REGISTRY_ADDR"
	CseConfigCenterAddress = "CSE_CONFIG_CENTER_ADDR"
	CseMonitorServer       = "CSE_MONITOR_SERVER_ADDR"
)

// env connect with "." like service_description.name and service_description.version which can not be used in k8s.
// So we can not use archaius to set env.
// To support this decalring constant for service name and version
// constant for service name and version.
const (
	ServiceName = "SERVICE_NAME"
	Version     = "VERSION"
)

// constant for microservice environment
const (
	EnvValueDev  = "development"
	EnvValueProd = "production"
)

// constant for secure socket layer parameters
const (
	SslCipherPluginKey = "cipherPlugin"
	SslVerifyPeerKey   = "verifyPeer"
	SslCipherSuitsKey  = "cipherSuits"
	SslProtocolKey     = "protocol"
	SslCaFileKey       = "caFile"
	SslCertFileKey     = "certFile"
	SslKeyFileKey      = "keyFile"
	SslCertPwdFileKey  = "certPwdFile"
	AKSKCustomCipher   = "cse.credentials.akskCustomCipher"
)

// constant for protocol types
const (
	ProtocolRest    = "rest"
	ProtocolHighway = "highway"
	LBSessionID     = "ServiceCombLB"
)

// DefaultKey default key
const DefaultKey = "default"

// DefaultValue default value
const DefaultValue = "default"

// BuildinTagApp build tag for the application
const BuildinTagApp = "app"

// BuildinTagVersion build tag version
const BuildinTagVersion = "version"

// CallerKey caller key
const CallerKey = "caller"

const (
	// HeaderSourceName is constant for header source name
	HeaderSourceName = "x-cse-src-microservice"
)

const (
	// RestMethod is the http method for restful protocol
	RestMethod = "method"
)

// constant for default application name and version
const (
	DefaultApp        = "default"
	DefaultVersion    = "0.0.1"
	LatestVersion     = "latest"
	AllVersion        = "0+"
	DefaultStatus     = "UP"
	DefaultLevel      = "BACK"
	DefaultHBInterval = 30
)

//constant used
const (
	HTTP              = "http"
	HTTPS             = "https"
	JSON              = "application/json"
	Create            = "CREATE"
	Update            = "UPDATE"
	Delete            = "DELETE"
	Size              = "size"
	Client            = "client"
	File              = "File"
	SessionID         = "sessionid"
	DefaultTenant     = "default"
	DefaultChainName  = "default"
	RollingPolicySize = "size"
	FileRegistry      = "File"
	DefaultUserName   = "default"
	DefaultDomainName = "default"
	DefaultProvider   = "default"
)

// const default config for config-center
const (
	DefaultRefreshMode = 1
)

//ContextValueKey is the key of value in context
type ContextValueKey struct{}

// NewContext transforms a metadata to context object
func NewContext(m map[string]string) context.Context {
	if m == nil {
		return context.WithValue(context.Background(), ContextValueKey{}, make(map[string]string, 0))
	}
	return context.WithValue(context.Background(), ContextValueKey{}, m)
}

// WithContext sets the KV and returns the context object
func WithContext(ctx context.Context, key, val string) context.Context {
	if ctx == nil {
		return context.WithValue(context.Background(), ContextValueKey{}, map[string]string{
			key: val,
		})
	}

	at, ok := ctx.Value(ContextValueKey{}).(map[string]string)
	if !ok {
		return context.WithValue(ctx, ContextValueKey{}, map[string]string{
			key: val,
		})
	}
	at[key] = val
	return ctx
}

// FromContext transforms a context object to metadata
func FromContext(ctx context.Context) map[string]string {
	if ctx == nil {
		return make(map[string]string, 0)
	}
	at, ok := ctx.Value(ContextValueKey{}).(map[string]string)
	if !ok {
		return make(map[string]string, 0)
	}
	return at
}
