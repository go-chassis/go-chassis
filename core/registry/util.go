package registry

import (
	"net"
	"net/url"
	"strings"
	"time"

	"crypto/tls"
	"fmt"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	chassisTLS "github.com/ServiceComb/go-chassis/core/tls"
	"github.com/ServiceComb/go-chassis/util/iputil"
	"github.com/cenkalti/backoff"
)

const protocolSymbol = "://"

//GetProtocolMap returns the protocol map
func GetProtocolMap(eps []string) (map[string]string, string) {
	m := make(map[string]string)
	var p string
	for _, ep := range eps {
		u, err := url.Parse(ep)
		if err != nil {
			lager.Logger.Error("Can not parse "+ep, err)
			continue
		}
		proto := u.Scheme
		ipPort := u.Host
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
				lager.Logger.Warnf("get port from listen addr failed.", err)
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
			var (
				ip  net.IP
				err error
			)

			// check the provided Advertise ip is IPV4 or IPV6
			ipWithoutPort := strings.Split(protocol.Advertise, ":")
			if len(ipWithoutPort) > 2 {
				ip, _, err = net.ParseCIDR(protocol.Advertise + "/0")
			} else {
				ip, _, err = net.ParseCIDR(ipWithoutPort[0] + "/0")
			}

			if err != nil {
				lager.Logger.Errorf(err, "failed to parse ip address")
			} else {
				if ip != nil && ip.To4() != nil {
					eps[name] = ip.String() + ":" + ipWithoutPort[1]
				}
			}
		}
	}
	return eps
}

//Microservice2ServiceKeyStr prepares a microservice key
func Microservice2ServiceKeyStr(m *MicroService) string {
	return strings.Join([]string{m.ServiceName, m.Version, m.AppID}, ":")
}

const (
	initialInterval = 5 * time.Second
	maxInterval     = 3 * time.Minute
)

func startBackOff(operation func() error) {
	backOff := &backoff.ExponentialBackOff{
		InitialInterval:     initialInterval,
		MaxInterval:         maxInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		Clock:               backoff.SystemClock,
	}
	for {
		lager.Logger.Infof("start backoff with initial interval %v", initialInterval)
		err := backoff.Retry(operation, backOff)
		if err == nil {
			return
		}
	}
}

//URIs2Hosts return hosts and scheme
func URIs2Hosts(uris []string) ([]string, string, error) {
	hosts := make([]string, 0, len(uris))
	var scheme string
	for _, addr := range uris {
		u, e := url.Parse(addr)
		if e != nil {
			//not uri. but still permitted, like zookeeper,file system
			hosts = append(hosts, u.Host)
			continue
		}
		if len(u.Host) == 0 {
			continue
		}
		if len(scheme) != 0 && u.Scheme != scheme {
			return nil, "", fmt.Errorf("inconsistent scheme found in registry address")
		}
		scheme = u.Scheme
		hosts = append(hosts, u.Host)

	}
	return hosts, scheme, nil
}
func getTLSConfig(scheme, t string) (*tls.Config, error) {
	var tlsConfig *tls.Config
	secure := scheme == common.HTTPS
	if secure {
		sslTag := t + "." + common.Consumer
		tmpTLSConfig, sslConfig, err := chassisTLS.GetTLSConfigByService(t, "", common.Consumer)
		if err != nil {
			if chassisTLS.IsSSLConfigNotExist(err) {
				tmpErr := fmt.Errorf("%s tls mode, but no ssl config", sslTag)
				lager.Logger.Error(tmpErr.Error(), err)
				return nil, tmpErr
			}
			lager.Logger.Errorf(err, "Load %s TLS config failed.", sslTag)
			return nil, err
		}
		lager.Logger.Warnf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
			sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin)
		tlsConfig = tmpTLSConfig
	}
	return tlsConfig, nil
}
