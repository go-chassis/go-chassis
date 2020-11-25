package storage

//DB is yaml file struct to set mongodb config
type Options struct {
	URI        *string `yaml:"uri"`
	PoolSize   *int    `yaml:"poolSize"`
	SSLEnabled *bool   `yaml:"sslEnabled"`
	RootCA     *string `yaml:"rootCAFile"`
	Timeout    *string `yaml:"timeout"`
	VerifyPeer *bool   `yaml:"verifyPeer"`
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) SetURI(uri string) *Options {
	o.URI = &uri
	return o
}

func (o *Options) SetPoolSize(poolSize int) *Options {
	o.PoolSize = &poolSize
	return o
}

func (o *Options) SetSSLEnabled(sslEnabled bool) *Options {
	o.SSLEnabled = &sslEnabled
	return o
}

func (o *Options) SetRootCA(rootCAFile string) *Options {
	o.RootCA = &rootCAFile
	return o
}

func (o *Options) SetTimeout(timeout string) *Options {
	o.Timeout = &timeout
	return o
}

func (o *Options) SetVerifyPeer(verifyPeer bool) *Options {
	o.VerifyPeer = &verifyPeer
	return o
}

func MergeOptions(opts ...*Options) *Options {
	options := NewOptions()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if opt.URI != nil {
			options.URI = opt.URI
		}
		if opt.PoolSize != nil {
			options.PoolSize = opt.PoolSize
		}
		if opt.SSLEnabled != nil {
			options.SSLEnabled = opt.SSLEnabled
		}
		if opt.RootCA != nil {
			options.RootCA = opt.RootCA
		}
		if opt.Timeout != nil {
			options.Timeout = opt.Timeout
		}
		if opt.VerifyPeer != nil {
			options.VerifyPeer = opt.VerifyPeer
		}
	}
	return options
}

func NewConfig(uri string, opts ...*Options) *Options {
	options := MergeOptions(opts...)
	if uri != "" && len(uri) != 0 {
		options.SetURI(uri)
	}
	return options
}
