package common

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
	EnvInstance   = "instance_description.environment"
	EnvSchemaRoot = "SCHEMA_ROOT"
	EnvProjectID  = "CSE_PROJECT_ID"
)

// constant for environment stage
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

// constant for default application name and version
const (
	DefaultApp     = "default"
	DefaultVersion = "0.0.1"
)

//constant used
const (
	HTTPS             = "https"
	JSON              = "application/json"
	Create            = "CREATE"
	Update            = "UPDATE"
	Delete            = "DELETE"
	Size              = "size"
	Client            = "client"
	File              = "File"
	SessionID         = "sessionid"
	ContentTypeJSON   = "application/json"
	DefaultTenant     = "default"
	DefaultChainName  = "default"
	RollingPolicySize = "size"
	FileRegistry      = "File"
	DefaultUserName   = "default"
	DefaultDomainName = "default"
	DefaultProvider   = "default"
)
