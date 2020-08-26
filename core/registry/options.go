package registry

import (
	"crypto/tls"
	"time"
)

// Options having micro-service parameters
type Options struct {
	Addrs      []string
	EnableSSL  bool
	Timeout    time.Duration
	TLSConfig  *tls.Config
	Compressed bool
	Verbose    bool
	Version    string
	ConfigPath string
}
