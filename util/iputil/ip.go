package iputil

import (
	"net"
	"os"
	"strings"

	"github.com/ServiceComb/go-chassis/core/common"
)

//Localhost is a function which returns localhost IP address
func Localhost() string { return "127.0.0.1" }

//GetLocalIP 获得本机IP
func GetLocalIP() string {
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
		// Check if Isloopback and IPV4
		if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
			return ip.String()
		}
	}
	return ""
}

//GetHostName is function which returns hostname
func GetHostName() string {
	if hostName, err := os.Hostname(); err == nil {
		return hostName
	}
	return "localhost"
}

// DefaultEndpoint4Protocol : To ensure consistency, we generate default addr for listenAddress and advertiseAddress by one method. To avoid unnecessary port allocation work, we allocate fixed port for user defined protocol.
func DefaultEndpoint4Protocol(proto string) string {
	return strings.Join([]string{Localhost(), DefaultPort4Protocol(proto)}, ":")
}

//DefaultPort4Protocol returns the default port for different protocols
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
