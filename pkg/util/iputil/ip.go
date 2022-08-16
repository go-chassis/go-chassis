package iputil

import (
	"crypto/tls"
	"fmt"
	"github.com/go-chassis/openlog"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chassis/go-chassis/v2/core/common"
)

// Localhost is a function which returns localhost IP address
func Localhost() string { return "127.0.0.1" }

// GetLocalIP 获得本机IP
func GetLocalIP() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return ""
	}
	for _, address := range addresses {
		// Parse IP
		var ip net.IP
		if ip, _, err = net.ParseCIDR(address.String()); err != nil {
			log.Println(err)
			return ""
		}
		// Check if valid global unicast IPv4 address
		if ip != nil && (ip.To4() != nil) && ip.IsGlobalUnicast() {
			return ip.String()
		}
	}
	return ""
}

// DefaultEndpoint4Protocol : To ensure consistency, we generate default addr for listenAddress and advertiseAddress by one method. To avoid unnecessary port allocation work, we allocate fixed port for user defined protocol.
func DefaultEndpoint4Protocol(proto string) string {
	return strings.Join([]string{Localhost(), DefaultPort4Protocol(proto)}, ":")
}

// DefaultPort4Protocol returns the default port for different protocols
func DefaultPort4Protocol(proto string) string {
	switch proto {
	case common.ProtocolRest:
		return "5000"
	case common.ProtocolHighway:
		return "6000"
	default:
		return "7000"
	}
}

// URIs2Hosts returns hosts and schema
func URIs2Hosts(uris []string) ([]string, string, error) {
	hosts := make([]string, 0, len(uris))
	var scheme string
	for _, addr := range uris {
		u, e := url.Parse(addr)
		if e != nil {
			//not uri. but still permitted, like zookeeper,file system
			openlog.Warn(fmt.Sprintf("parse address failed, %s", e.Error()))
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

// GetLocalIPv6 Get IPv6 address of NIC.
func GetLocalIPv6() string {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addresses {
		// Parse IP
		var ip net.IP
		if ip, _, err = net.ParseCIDR(address.String()); err != nil {
			return ""
		}
		// Check if valid IPv6 address
		if ip != nil && (ip.To16() != nil) && IsIPv6Address(ip) && ip.IsGlobalUnicast() {
			return ip.String()
		}
	}
	return ""
}

// IsIPv6Address check whether the IP is IPv6 address.
func IsIPv6Address(ip net.IP) bool {
	if ip != nil && strings.Contains(ip.String(), ":") {
		return true
	}
	return false
}

// StartListener start listener with address and tls(if has), returns the listener and the real listened ip/port
func StartListener(listenAddress string, tlsConfig *tls.Config) (listener net.Listener, listenedIP string, port string, err error) {
	if tlsConfig == nil {
		listener, err = net.Listen("tcp", listenAddress)
	} else {
		listener, err = tls.Listen("tcp", listenAddress, tlsConfig)
	}
	if err != nil {
		return
	}
	realAddr := listener.Addr().String()
	listenedIP, port, err = net.SplitHostPort(realAddr)
	if err != nil {
		return
	}
	ip := net.ParseIP(listenedIP)
	if ip.IsUnspecified() {
		if IsIPv6Address(ip) {
			listenedIP = GetLocalIPv6()
			if listenedIP == "" {
				listenedIP = GetLocalIP()
			}
		} else {
			listenedIP = GetLocalIP()
		}
	}
	return
}

// ClientIP returns client ip
func ClientIP(r *http.Request) string {
	ips := ForwardedIPs(r)
	if len(ips) > 0 {
		ip := ips[0]
		if !strings.Contains(ip, ":") {
			return ip
		}
		rip, _, err := net.SplitHostPort(ip)
		if err != nil {
			openlog.Warn(fmt.Sprintf("get client ip catch a err, %s", err.Error()))
			return ip
		}
		return rip
	}

	realIP := RealIP(r)
	if len(realIP) > 0 {
		return realIP
	}
	return RemoteIP(r)
}

// RemoteIP returns remote ip
func RemoteIP(r *http.Request) string {
	remoteIP := r.RemoteAddr
	rip, _, err := net.SplitHostPort(remoteIP)
	if err != nil {
		openlog.Warn(fmt.Sprintf("get remote ip catch a err, %s", err.Error()))
		return remoteIP
	}
	return rip
}

// ForwardedIPs returns forwarded for ips
func ForwardedIPs(r *http.Request) []string {
	ips := r.Header.Get("X-Forwarded-For")
	if len(ips) == 0 {
		return []string{}
	}
	return strings.Split(ips, ",")
}

// RealIP returns real ip
func RealIP(r *http.Request) string {
	return r.Header.Get("X-Real-Ip")
}
