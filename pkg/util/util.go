package util

import (
	"errors"
	"strings"
)

//ErrInvalidPortName happens if your port name is illegal
var ErrInvalidPortName = errors.New("invalid port name")

//ParsePortName a port name is composite by protocol-name,like http-admin,http-api,grpc-console,grpc-api
//ParsePortName return two string separately
func ParsePortName(n string) (string, string, error) {
	tmp := strings.Split(n, "-")
	if len(tmp) != 2 {
		return "", "", ErrInvalidPortName
	}
	return tmp[0], tmp[1], nil
}
