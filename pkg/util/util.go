package util

import (
	"errors"
	"strings"
)

//ErrInvalidPortName happens if your port name is illegal
var ErrInvalidPortName = errors.New("invalid port name,port name must be {protocol}-<suffix>")

//ParsePortName a port name is composite by protocol-name,like http-admin,http-api,grpc-console,grpc-api
//ParsePortName return two string separately
func ParsePortName(n string) (string, string, error) {
	if n == "" {
		return "", "", ErrInvalidPortName
	}
	tmp := strings.Split(n, "-")
	switch len(tmp) {
	case 2:
		return tmp[0], tmp[1], nil
	case 1:
		return tmp[0], "", nil
	default:
		return "", "", ErrInvalidPortName
	}

}
