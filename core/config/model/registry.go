package model

// RegistryStruct SC information
type RegistryStruct struct {
	Disable         bool                     `yaml:"disabled"`
	Type            string                   `yaml:"type"`
	Scope           string                   `yaml:"scope"`
	AutoDiscovery   bool                     `yaml:"autodiscovery"`
	AutoIPIndex     bool                     `yaml:"autoIPIndex"`
	Address         string                   `yaml:"address"`
	RefreshInterval string                   `yaml:"refreshInterval"`
	Watch           bool                     `yaml:"watch"`
	AutoRegister    string                   `yaml:"register"`
	APIVersion      RegistryAPIVersionStruct `yaml:"api"`

	HealthCheck bool      `yaml:"healthCheck"`
	CacheIndex  bool      `yaml:"cacheIndex"`
	ConfigPath  string    `yaml:"configPath"`
	Heartbeat   Heartbeat `yaml:"heartbeat"`
}

// RegistryAPIVersionStruct registry api version structure
type RegistryAPIVersionStruct struct {
	Version string `yaml:"version"`
}

// Heartbeat indicates heartbeat information
type Heartbeat struct {
	Mode     string `yaml:"mode"`
	Interval string `yaml:"interval"`
}
