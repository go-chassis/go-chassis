package registry

import (
	"strings"
)

const (
	ssLEnabledQuery = "sslEnabled=true"
)

// Endpoint struct having full info about micro-service instance endpoint
type Endpoint struct {
	SSLEnabled bool   `json:"sslEnabled"`
	Address    string `json:"address"`
}

// NewEndPoint return a Endpoint object what parse from url
func NewEndPoint(schema string) (*Endpoint, error) {
	return parseAddress(schema)
}

//GenEndpoint return the endpoint string which it contain the sslEnabled=true query arg or not
func (e *Endpoint) GenEndpoint() string {
	if e.SSLEnabled {
		return e.Address + "?" + ssLEnabledQuery
	}
	return e.Address
}

//IsSSLEnable return it is use ssl or not
func (e *Endpoint) IsSSLEnable() bool {
	return e.SSLEnabled
}

//SetSSLEnable set ssl enable or not
func (e *Endpoint) SetSSLEnable(enabled bool) {
	e.SSLEnabled = enabled
}

func (e *Endpoint) String() string {
	return e.GenEndpoint()
}

func parseAddress(address string) (*Endpoint, error) {
	ep := Endpoint{}
	idx := strings.Index(address, "?")
	if idx != -1 {
		if strings.Contains(address, ssLEnabledQuery) {
			ep.SSLEnabled = true
		}
		address = address[:idx]
	}
	if pIdx := strings.Index(address, ":"); pIdx == -1 {
		ep.Address = address
		return &ep, nil
	}
	ep.Address = address
	return &ep, nil
}
