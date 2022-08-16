package registry

import (
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"crypto/tls"
	"fmt"

	"github.com/cenkalti/backoff"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	chassisTLS "github.com/go-chassis/go-chassis/v2/core/tls"
	"github.com/go-chassis/go-chassis/v2/pkg/util/iputil"
	"github.com/go-chassis/openlog"
)

const (
	protocolSymbol = "://"
)

// GetProtocolMap returns the protocol map
func GetProtocolMap(eps []string) (map[string]*Endpoint, string) {
	m := make(map[string]*Endpoint)
	var p string
	for _, addr := range eps {
		proto := ""
		ep := ""
		idx := strings.Index(addr, protocolSymbol)
		if idx == -1 {
			ep = addr
			proto = "unknown"
		} else {
			ep = addr[idx+len(protocolSymbol):]
			proto = addr[:idx]
		}
		u, err := NewEndPoint(ep)
		if err != nil {
			openlog.Error(fmt.Sprintf("Can not parse %s, error %s", ep, err))
			continue
		}
		m[proto] = u
		p = proto
	}
	return m, p
}

// GetProtocolList returns the protocol list
func GetProtocolList(m map[string]*Endpoint) []string {
	eps := []string{}
	for p, ep := range m {
		uri := p + protocolSymbol + ep.GenEndpoint()
		eps = append(eps, uri)
	}
	return eps
}

// MakeEndpoints returns the endpoints
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

// MakeEndpointMap returns the endpoints map
func MakeEndpointMap(m map[string]model.Protocol) (map[string]*Endpoint, error) {
	eps := make(map[string]*Endpoint)
	for name, protocol := range m {
		ep := protocol.Listen
		if len(protocol.Advertise) > 0 {
			ep = protocol.Advertise
		}

		host, port, err := net.SplitHostPort(ep)
		if err != nil {
			return nil, err
		}
		if host == "" || port == "" {
			return nil, fmt.Errorf("listen address is invalid [%s]", protocol.Listen)
		}

		_, err = FillUnspecifiedIP(host)
		if err != nil {
			return nil, err
		}
		if endpoint, err := NewEndPoint(ep); err == nil {
			eps[name] = endpoint
		}

	}
	return eps, nil
}

// FillUnspecifiedIP replace 0.0.0.0 or :: IPv4 and IPv6 unspecified IP address with local NIC IP.
func FillUnspecifiedIP(host string) (string, error) {
	var addr string
	ip := net.ParseIP(host)
	if ip == nil {
		return "", fmt.Errorf("invalid IP address %s", host)
	}

	addr = host
	if ip.IsUnspecified() {
		if iputil.IsIPv6Address(ip) {
			addr = iputil.GetLocalIPv6()
		} else {
			addr = iputil.GetLocalIP()
		}
		if len(addr) == 0 {
			return addr, fmt.Errorf("auto generate IP address failed, plz manually set listenAddress and advertiseAddress")
		}
	}
	return addr, nil
}

// Microservice2ServiceKeyStr prepares a microservice key
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
		openlog.Info(fmt.Sprintf("start backoff with initial interval %v", initialInterval))
		err := backoff.Retry(operation, backOff)
		if err == nil {
			return
		}
	}
}

// URIs2Hosts return hosts and scheme
func URIs2Hosts(uris []string) ([]string, string, error) {
	hosts := make([]string, 0)
	var scheme string
	var URIRegex = "(\\.*://.*)"
	reg, err := regexp.Compile(URIRegex)
	if err != nil {
		return nil, "", err
	}
	for _, addr := range uris {
		ok := reg.MatchString(addr)
		if ok {
			u, e := url.Parse(addr)
			if e != nil {
				openlog.Warn("registry address is invalid:" + addr)
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
		} else {
			//not uri. but still permitted, like zookeeper,file system
			hosts = append(hosts, addr)
			continue
		}
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
				openlog.Error(tmpErr.Error() + ", err: " + err.Error())
				return nil, tmpErr
			}
			openlog.Error(fmt.Sprintf("Load TLS config failed: %s", err))
			return nil, err
		}
		openlog.Warn(fmt.Sprintf("%s TLS mode, verify peer: %t, cipher plugin: %s.",
			sslTag, sslConfig.VerifyPeer, sslConfig.CipherPlugin))
		tlsConfig = tmpTLSConfig
	}
	return tlsConfig, nil
}

// GetDuration return the time.Duration type value by specified key
func GetDuration(key string, def time.Duration) time.Duration {
	str := strings.TrimSpace(key)
	if str == "" {
		return def
	}
	d, err := time.ParseDuration(str)
	if err != nil {
		return def
	}
	return d
}
