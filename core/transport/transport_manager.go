package transport

import (
	"fmt"
	"log"

	"github.com/ServiceComb/go-chassis/core/lager"
	microTransport "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
)

// TransortFunc transport function
type TransortFunc func(...microTransport.Option) microTransport.Transport

var transportFuncMap = make(map[string]TransortFunc)

// TransportMap transport map
var TransportMap = make(map[string]microTransport.Transport)

// InstallPlugin install plugin of transport
func InstallPlugin(protocol string, newFunc TransortFunc) {
	log.Printf("Install transport Plugin, protocol=%s", protocol)
	transportFuncMap[protocol] = newFunc
}

// GetTransportFunc get transport function
func GetTransportFunc(protocol string) (TransortFunc, error) {
	newFunc := transportFuncMap[protocol]
	if newFunc == nil {
		lager.Logger.Errorf(nil, "Do not support, protocol:%s", protocol)
		return nil, fmt.Errorf("Don't support [%s] transport", protocol)
	}
	return transportFuncMap[protocol], nil
}

// CreateTransport create transport function
func CreateTransport(protocol string, opts ...microTransport.Option) {
	trFunc := transportFuncMap[protocol]
	if trFunc == nil {
		lager.Logger.Warnf(nil, "Doesn't support this protocol:%s", protocol)
		return
	}

	TransportMap[protocol] = trFunc(opts...)
}

// GetTransport get transport
func GetTransport(protocol string) microTransport.Transport {
	return TransportMap[protocol]
}

// Init intilize transport
func Init() {
	//TODO hard code. must read from
	/*
	  grpc:
	    address: 10.146.207.197:37325
	  rest:
	    address: 0.0.0.0:8080
	  highway:
	    address: 0.0.0.0:7070
	*/
	protocols := []string{"grpc", "tcp"}
	for _, protocol := range protocols {
		CreateTransport(protocol)
	}

}
