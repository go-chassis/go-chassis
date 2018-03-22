package client

import (
	"crypto/tls"
	"time"
)

//Options is configs for client creation
type Options struct {
	PoolSize  int
	PoolTTL   time.Duration
	TLSConfig *tls.Config
	Failure   map[string]bool
}
