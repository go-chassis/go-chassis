package model

type TestCfg struct {
	ThreadCount            int    `yaml:"threadcount"`
	MessageSize            int    `yaml:"messagesize"`
	PrintCount             int    `yaml:"printcount"`
	RegistryEnable         bool   `yaml:"registryenable"`
	Protocol               string `yaml:"protocol"`
	EndPoint               string `yaml:"endpoint"`
	GoMaxProcs             int    `yaml:"gomaxprocs"`
	MicroServiceName       string `yaml:"microservicename"`
	DependMicroServiceName string `yaml:"dependmicroservicename"`
	EnvStage               string `yaml:"envstage"`
}
