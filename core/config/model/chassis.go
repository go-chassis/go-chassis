package model

import "time"

//GlobalCfg chassis.yaml 配置项
type GlobalCfg struct {
	AppID      string            `yaml:"APPLICATION_ID"`
	Cse        CseStruct         `yaml:"cse"`
	Ssl        map[string]string `yaml:"ssl"`
	Tracing    TracingStruct     `yaml:"tracing"`
	DataCenter *DataCenterInfo   `yaml:"region"`
}

// DataCenterInfo gives data center information
type DataCenterInfo struct {
	Name          string `yaml:"name"`
	Region        string `yaml:"region"`
	AvailableZone string `yaml:"availableZone"`
}

//PassLagerCfg is the struct for lager information(passlager.yaml)
type PassLagerCfg struct {
	Writers        string `yaml:"writers"`
	LoggerLevel    string `yaml:"logger_level"`
	LoggerFile     string `yaml:"logger_file"`
	LogFormatText  bool   `yaml:"log_format_text"`
	RollingPolicy  string `yaml:"rollingPolicy"`
	LogRotateDate  int    `yaml:"log_rotate_date"`
	LogRotateSize  int    `yaml:"log_rotate_size"`
	LogBackupCount int    `yaml:"log_backup_count"`
}

//CseStruct 设置注册中心SC的地址，要开哪些传输协议， 调用链信息等
type CseStruct struct {
	Loadbalance LoadBalanceStruct           `yaml:"loadbalance"`
	Config      ConfigStruct                `yaml:"config"`
	Service     ServiceStruct               `yaml:"service"`
	Protocols   map[string]Protocol         `yaml:"protocols"`
	Handler     HandlerStruct               `yaml:"handler"`
	References  map[string]ReferencesStruct `yaml:"references"`
	FlowControl ServiceTypes                `yaml:"flowcontrol"`
	Monitor     MonitorStruct               `yaml:"monitor"`
	Metrics     MetricsStruct               `yaml:"metrics"`
	Credentials CredentialStruct            `yaml:"credentials"`
}

// FaultProtocolStruct fault protocol struct
type FaultProtocolStruct struct {
	Fault map[string]Fault `yaml:"protocols"`
}

// Fault fault struct
type Fault struct {
	Abort Abort `yaml:"abort"`
	Delay Delay `yaml:"delay"`
}

// Abort abort struct
type Abort struct {
	Percent    int `yaml:"percent"`
	HTTPStatus int `yaml:"httpStatus"`
}

// Delay delay struct
type Delay struct {
	Percent    int           `yaml:"percent"`
	FixedDelay time.Duration `yaml:"fixedDelay"`
}

// MetricsStruct metrics struct
type MetricsStruct struct {
	APIPath                string `yaml:"apiPath"`
	Enable                 bool   `yaml:"enable"`
	EnableGoRuntimeMetrics bool   `yaml:"enableGoRuntimeMetrics"`
}

// MonitorStruct is the struct for monitoring parameters
type MonitorStruct struct {
	Client MonitorClientStruct `yaml:"client"`
}

// MonitorClientStruct monitor client struct
type MonitorClientStruct struct {
	ServerURI  string                  `yaml:"serverUri"`
	Enable     bool                    `yaml:"enable"`
	UserName   string                  `yaml:"userName"`
	DomainName string                  `yaml:"domainName"`
	APIVersion MonitorAPIVersionStruct `yaml:"api"`
}

// MonitorAPIVersionStruct monitor API version struct
type MonitorAPIVersionStruct struct {
	Version string `yaml:"version"`
}

// ServiceTypes gives the information of service types
type ServiceTypes struct {
	Consumer TypesStruct `yaml:"Consumer"`
	Provider TypesStruct `yaml:"Provider"`
}

// TypesStruct is the struct for QPS
type TypesStruct struct {
	QPS QPSStruct `yaml:"qps"`
}

// QPSStruct QPS struct
type QPSStruct struct {
	Enabled bool              `yaml:"enabled"`
	Global  map[string]int    `yaml:"global"`
	Limit   map[string]string `yaml:"limit"`
}

// ConfigStruct configuration structure
type ConfigStruct struct {
	Client ClientStruct `yaml:"client"`
}

// ClientStruct client structure
type ClientStruct struct {
	ServerURI       string                 `yaml:"serverUri"`
	TenantName      string                 `yaml:"tenantName"`
	RefreshMode     int                    `yaml:"refreshMode"`
	RefreshInterval int                    `yaml:"refreshInterval"`
	RefreshPort     string                 `yaml:"refreshPort"`
	Autodiscovery   bool                   `yaml:"autodiscovery"`
	APIVersion      ConfigAPIVersionStruct `yaml:"api"`
}

// ConfigAPIVersionStruct is the structure for configuration API version
type ConfigAPIVersionStruct struct {
	Version string `yaml:"version"`
}

// ReferencesStruct references structure
type ReferencesStruct struct {
	Version   string `yaml:"version"`
	Transport string `yaml:"transport"`
}

// LoadBalanceStruct loadbalancing structure
type LoadBalanceStruct struct {
	Strategy     map[string]string `yaml:"strategy"`
	RetryEnabled bool              `yaml:"retryEnabled"`
	RetryOnNext  int               `yaml:"retryOnNext"`
	RetryOnSame  int               `yaml:"retryOnSame"`
	Backoff      BackoffStrategy   `yaml:"backoff"`
}

// BackoffStrategy back off strategy
type BackoffStrategy struct {
	Kind  string `yaml:"kind"`
	MinMs uint   `yaml:"minMs"`
	MaxMs uint   `yaml:"maxMs"`
}

// Protocol protocol structure
type Protocol struct {
	Listen       string `yaml:"listenAddress"`
	Advertise    string `yaml:"advertiseAddress"`
	WorkerNumber int    `yaml:"workerNumber"`
	Transport    string `yaml:"transport"`
	Failure      string `yaml:"failure"`
}

//ServiceStruct SC注册中心地址信息结构体
type ServiceStruct struct {
	Registry RegistryStruct `yaml:"registry"`
}

//RegistryStruct SC注册中心地址信息
type RegistryStruct struct {
	Disable         bool                     `yaml:"disabled"`
	Type            string                   `yaml:"type"`
	Scope           string                   `yaml:"scope"`
	AutoDiscovery   bool                     `yaml:"autodiscovery"`
	AutoIPIndex     bool                     `yaml:"autoIPIndex"`
	Address         string                   `yaml:"address"`
	RefreshInterval string                   `yaml:"refeshInterval"`
	Watch           bool                     `yaml:"watch"`
	Tenant          string                   `yaml:"tenant"`
	AutoRegister    string                   `yaml:"register"`
	APIVersion      RegistryAPIVersionStruct `yaml:"api"`
}

// RegistryAPIVersionStruct registry api version structure
type RegistryAPIVersionStruct struct {
	Version string `yaml:"version"`
}

// MicroserviceCfg microservice.yaml 配置项
type MicroserviceCfg struct {
	AppID               string           `yaml:"APPLICATION_ID"`
	Provider            string           `yaml:"Provider"`
	ServiceDescription  MicServiceStruct `yaml:"service_description"`
	InstanceDescription MicServiceStruct `yaml:"instance_description"`
	Cse                 MicCseStruct     `yaml:"cse"`
}

// MicServiceStruct ServiceStruct 设置微服务的私有属性
type MicServiceStruct struct {
	Name               string            `yaml:"name"`
	Version            string            `yaml:"version"`
	Level              string            `yaml:"level"`
	Properties         map[string]string `yaml:"properties"`
	InstanceProperties map[string]string `yaml:"instance_properties"`
}

// InstanceDesc is the struct for instance description
type InstanceDesc struct {
	Env string `yaml:"environment"`
}

// MicCseStruct 设置注册中心SC的地址，要开哪些传输协议， 调用链信息等
type MicCseStruct struct {
	RPC     map[string]string `yaml:"rpc"`
	TCP     map[string]string `yaml:"tcp"`
	Rest    map[string]string `yaml:"rest"`
	HighWay map[string]string `yaml:"highway"`
	Handler HandlerStruct     `yaml:"handler"`
}

// HandlerStruct 调用链信息
type HandlerStruct struct {
	Chain ChainStruct `yaml:"chain"`
}

// ChainStruct 调用链信息
type ChainStruct struct {
	Consumer map[string]string `yaml:"Consumer"`
	Provider map[string]string `yaml:"Provider"`
}

// CredentialStruct aksk信息
type CredentialStruct struct {
	AccessKey        string `yaml:"accessKey"`
	SecretKey        string `yaml:"secretKey"`
	AkskCustomCipher string `yaml:"akskCustomCipher"`
	Project          string `yaml:"project"`
}

// TracingStruct tracing structure
type TracingStruct struct {
	SamplingRate  float64 `yaml:"samplingRate"`
	CollectorType string  `yaml:"collectorType"` // http|log
	// if collectorType is http, the target is zipkin server
	// if collectorType is log, the target is log file
	CollectorTarget string `yaml:"collectorTarget"`
}
