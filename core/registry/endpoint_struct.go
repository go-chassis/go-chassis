package registry

import (
	"net"
	"strings"
)

const (
	ssLEnabledQuery = "sslEnabled=true"
)

// EndPoint struct having full info about micro-service instance endpoint
type EndPoint struct {
	SslEnabled bool
	HostOrIP   string
	Port       string
}

// NewEndPoint return a Endpoint object what parse from url
func NewEndPoint(schema string) (*EndPoint, error) {
	return parseAddress(schema)
}

func (e *EndPoint) Host() string {
	if e.Port == "" {
		return e.HostOrIP
	}
	return net.JoinHostPort(e.HostOrIP, e.Port)
}

func (e *EndPoint) GenEndpoint() string {
	sslFlag := ""
	if e.SslEnabled {
		sslFlag = "?" + ssLEnabledQuery
	}

	if e.Port == "" {
		return e.HostOrIP + sslFlag
	}
	return net.JoinHostPort(e.HostOrIP, e.Port) + sslFlag
}

func (e *EndPoint) IsSSLEnable() bool {
	return e.SslEnabled
}

func (e *EndPoint) SetSSLEnable(enabled bool) {
	e.SslEnabled = enabled
}

func (e *EndPoint) String() string {
	return e.GenEndpoint()
}

func parseAddress(address string) (*EndPoint, error) {
	ep := EndPoint{}
	idx := strings.Index(address, "?")
	if idx != -1 {
		if strings.Contains(address, ssLEnabledQuery) {
			ep.SslEnabled = true
		}
		address = address[:idx]
	}
	if pIdx := strings.Index(address, ":"); pIdx == -1 {
		ep.HostOrIP = address
		return &ep, nil
	}
	ip, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	ep.HostOrIP = ip
	ep.Port = port
	return &ep, nil
}
