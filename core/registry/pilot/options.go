package pilot

import "crypto/tls"

// Options is the list of parameters which passed to the EnvoyDSClient while creating a new client
type Options struct {
	Addrs     []string
	TLSConfig *tls.Config
}
