package model

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

// Protocol protocol structure
type Protocol struct {
	Listen       string `yaml:"listenAddress"`
	Advertise    string `yaml:"advertiseAddress"`
	WorkerNumber int    `yaml:"workerNumber"`
	Transport    string `yaml:"transport"`
	Failure      string `yaml:"failure"`
}

// MicroserviceCfg microservice.yaml 配置项
type MicroserviceCfg struct {
	AppID              string           `yaml:"APPLICATION_ID"`
	Provider           string           `yaml:"Provider"`
	ServiceDescription MicServiceStruct `yaml:"service_description"`
}

// MicServiceStruct ServiceStruct 设置微服务的私有属性
type MicServiceStruct struct {
	Name               string            `yaml:"name"`
	Version            string            `yaml:"version"`
	Environment        string            `yaml:"environment"`
	Level              string            `yaml:"level"`
	Properties         map[string]string `yaml:"properties"`
	InstanceProperties map[string]string `yaml:"instance_properties"`
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
