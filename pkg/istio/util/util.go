package util

import (
	"fmt"
	"os"
	"strings"
)

const (
	// PODNAMESPACE means pod's namespace in kubernetes
	PODNAMESPACE = "POD_NAMESPACE"
	// PODNAME means pod's name in kubernetes
	PODNAME = "POD_NAME"
	// PODIP means pod's instance ip in kubernetes
	PODIP = "INSTANCE_IP"

	defaultSuffix = "svc.cluster.local"
	defaultNS     = "default"
)

const (
	// RDSHttpProxy query all route configuration
	RDSHttpProxy = "http_proxy"
	// EnvoyAPIV2 defines prefix of type
	EnvoyAPIV2 = "type.googleapis.com/envoy.api.v2."
	// RouteType defines ADS type
	RouteType = EnvoyAPIV2 + "RouteConfiguration"
)

// ServiceKey returns service key from a service name
func ServiceKey(service string) string {
	ns := os.Getenv(PODNAMESPACE)
	if ns == "" {
		ns = defaultNS
	}
	return strings.Join([]string{service, ns, defaultSuffix}, ".")
}

// ServiceKeyToLabel returns label from service key
func ServiceKeyToLabel(service string) string {
	ss := strings.Split(service, "|")
	if len(ss) != 4 {
		return ""
	}
	return ss[2]
}

// ServiceAndPort returns service and port
func ServiceAndPort(host string) (string, string) {
	sp := strings.Split(host, ":")
	if len(sp) <= 1 {
		return host, "0"
	}
	ss := strings.Split(sp[0], ".")
	if len(ss) <= 1 {
		return sp[0], sp[1]
	}
	return ss[0], sp[1]
}

// BuildNodeID returns nodeID
func BuildNodeID() string {
	ns := os.Getenv(PODNAMESPACE)
	if ns == "" {
		ns = defaultNS
	}
	return fmt.Sprintf("sidecar~%s~%s.%s~%s.svc.cluster.local",
		os.Getenv(PODIP), os.Getenv(PODNAME), ns, ns)
}
