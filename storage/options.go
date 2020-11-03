package storage

//DB is yaml file struct to set mongodb config
type DB struct {
	URI        string `yaml:"uri"`
	PoolSize   int    `yaml:"poolSize"`
	SSLEnabled bool   `yaml:"sslEnabled"`
	RootCA     string `yaml:"rootCAFile"`
	Timeout    string `yaml:"timeout"`
	VerifyPeer bool   `yaml:"verifyPeer"`
}
