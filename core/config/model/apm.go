package model

//MonitorCfg monitoring.yaml 配置项
type MonitorCfg struct {
	ServiceComb ServiceCombStruct `yaml:"servicecomb"`
}

// ServiceComb structure
type ServiceCombStruct struct {
	APM APMStruct `yaml:"apm"`
}

//APM is for Application Performance Management
type APMStruct struct {
	Tracing TracingStruct `yaml:"tracing"`
}
