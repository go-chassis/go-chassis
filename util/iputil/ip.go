package iputil

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/ServiceComb/go-chassis/core/common"
)

//Localhost is a function which returns localhost IP address
func Localhost() string { return "127.0.0.1" }

//GetLocalIP 获得本机IP
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
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

//DefaultIPv4BindAddress is a function which binds IP address for corresponding interface
func DefaultIPv4BindAddress() (net.IP, error) {
	intfs, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("Unable to get interfaces: %v", err)
	}
	if len(intfs) == 0 {
		return nil, fmt.Errorf("no interfaces found on host")
	}

	for _, intf := range intfs {
		ip, err := getIPFromInterface(&intf)
		if err != nil {
			return nil, fmt.Errorf("Unable to get ip from interface %q: %v", intf.Name, err)
		}
		if ip != nil {
			return ip, nil
		}
	}

	return nil, fmt.Errorf("no acceptable interface with global unicast address found on host")
}

func getIPFromInterface(intf *net.Interface) (net.IP, error) {
	if !isInterfaceUp(intf) || isLoopbackOrPointToPoint(intf) {
		return nil, nil
	}
	addrs, err := intf.Addrs()
	if err != nil {
		return nil, err
	}
	if len(addrs) == 0 {
		return nil, nil
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, fmt.Errorf("Unable to parse CIDR: %s", err)
		}
		if ip.To4() != nil && ip.IsGlobalUnicast() {
			return ip, nil
		}
	}
	return nil, nil
}

func isInterfaceUp(intf *net.Interface) bool {
	if intf == nil {
		return false
	}
	if intf.Flags&net.FlagUp != 0 {
		return true
	}
	return false
}

func isLoopbackOrPointToPoint(intf *net.Interface) bool {
	return intf.Flags&(net.FlagLoopback|net.FlagPointToPoint) != 0
}
