package storage

//DB is yaml file struct to set mongodb config
type Options struct {
	URI        string `yaml:"uri"`
	PoolSize   int    `yaml:"poolSize"`
	SSLEnabled bool   `yaml:"sslEnabled"`
	RootCA     string `yaml:"rootCAFile"`
	Timeout    string `yaml:"timeout"`
	VerifyPeer bool   `yaml:"verifyPeer"`
}

type Option func(opt *Options)

func PoolSize(poolSize int) Option {
	return func(opt *Options) {
		opt.PoolSize = poolSize
	}
}

func SSLEnabled(sslEnabled bool) Option {
	return func(opt *Options) {
		opt.SSLEnabled = sslEnabled
	}
}

func RootCA(rootCAFile string) Option {
	return func(opt *Options) {
		opt.RootCA = rootCAFile
	}
}

func Timeout(timeout string) Option {
	return func(opt *Options) {
		opt.Timeout = timeout
	}
}

func VerifyPeer(verifyPeer bool) Option {
	return func(opt *Options) {
		opt.VerifyPeer = verifyPeer
	}
}

func NewConfig(uri string, opts ...func(opt *Options)) Options {
	opt := Options{
		URI: uri,
	}
	for _, option := range opts {
		option(&opt)
	}
	return opt
}
