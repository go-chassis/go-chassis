package model

// MonitorCfg monitoring.yaml 配置项
type MonitorCfg struct {
	ServiceComb ServiceCombStruct `yaml:"servicecomb"`
}

// ServiceCombStruct structure is for config of servicecomb
type ServiceCombStruct struct {
	APM APMStruct `yaml:"apm"`
}

// APMStruct is for Application Performance Management
type APMStruct struct {
	Tracing TracingStruct `yaml:"tracing"`
}
