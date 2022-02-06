package model

type Cors struct {
	Enable         bool     `yaml:"enable"`
	ExposeHeaders  []string `yaml:"exposeHeaders"`
	AllowedHeaders []string `yaml:"allowedHeaders"`
	AllowedDomains []string `yaml:"allowedDomains"`
	AllowedMethods []string `yaml:"allowedMethods"`
	CookiesAllowed bool     `yaml:"cookiesAllowed"`
	MaxAge         int      `yaml:"maxAge"`
}
