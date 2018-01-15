package registry

import (
	"net"
	"strings"

	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	stringutil "github.com/ServiceComb/go-chassis/core/util/string"
	"github.com/ServiceComb/go-chassis/util/iputil"
)

const protocolSymbol = "://"

//GetProtocolMap returns the protocol map
func GetProtocolMap(eps []string) (map[string]string, string) {
	m := make(map[string]string)
	var p string
	for _, ep := range eps {
		proto, ipPort := stringutil.SplitToTwo(ep, protocolSymbol)
		if proto == "" {
			m["unknown"] = ipPort
		} else {
			m[proto] = ipPort
			p = proto
		}
	}
	return m, p
}

//GetProtocolList returns the protocol list
func GetProtocolList(m map[string]string) []string {
	eps := []string{}
	for p, ep := range m {
		uri := p + protocolSymbol + ep
		eps = append(eps, uri)
	}
	return eps
}

//MakeEndpoints returns the endpoints
func MakeEndpoints(m map[string]model.Protocol) []string {
	var eps = make([]string, 0)
	for name, protocol := range m {
		ep := protocol.Advertise
		if ep == "" {
			if protocol.Listen != "" {
				ep = protocol.Listen
			} else {
				ep = iputil.DefaultEndpoint4Protocol(name)
			}
		}
		ep = strings.Join([]string{name, ep}, protocolSymbol)
		eps = append(eps, ep)
	}
	return eps
}

//MakeEndpointMap returns the endpoints map
func MakeEndpointMap(m map[string]model.Protocol) map[string]string {
	eps := make(map[string]string, 0)
	for name, protocol := range m {

		if len(protocol.Advertise) == 0 {
			host, port, err := net.SplitHostPort(protocol.Listen)
			if err != nil {
				lager.Logger.Warn("get port from listen addr failed.", err)
				port = iputil.DefaultPort4Protocol(name)
				host = iputil.Localhost()
			}

			if host != "" {
				if host == "0.0.0.0" {
					host = iputil.GetLocalIP()
				}
				eps[name] = strings.Join([]string{host, port}, ":")

			} else {
				eps[name] = iputil.DefaultEndpoint4Protocol(name)
			}
		} else {
			eps[name] = protocol.Advertise
		}
	}
	return eps
}

//Microservice2ServiceKeyStr prepares a microservice key
func Microservice2ServiceKeyStr(m *MicroService) string {
	return strings.Join([]string{m.ServiceName, m.Version, m.AppID}, ":")
}
