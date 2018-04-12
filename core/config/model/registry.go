package model

//ServiceStruct SC注册中心地址信息结构体
type ServiceStruct struct {
	Registry          RegistryStruct `yaml:"registry"`
	Registrator       RegistryStruct `yaml:"registrator"`
	ServiceDiscovery  RegistryStruct `yaml:"serviceDiscovery"`
	ContractDiscovery RegistryStruct `yaml:"contractDiscovery"`
}

//RegistryStruct SC注册中心地址信息
type RegistryStruct struct {
	Disable         bool                     `yaml:"disabled"`
	Type            string                   `yaml:"type"`
	Scope           string                   `yaml:"scope"`
	AutoDiscovery   bool                     `yaml:"autodiscovery"`
	AutoIPIndex     bool                     `yaml:"autoIPIndex"`
	Address         string                   `yaml:"address"`
	RefreshInterval string                   `yaml:"refreshInterval"`
	Watch           bool                     `yaml:"watch"`
	Tenant          string                   `yaml:"tenant"`
	AutoRegister    string                   `yaml:"register"`
	APIVersion      RegistryAPIVersionStruct `yaml:"api"`
}

// RegistryAPIVersionStruct registry api version structure
type RegistryAPIVersionStruct struct {
	Version string `yaml:"version"`
}
