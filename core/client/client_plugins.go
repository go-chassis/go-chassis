package client

import (
	"fmt"
	"log"
)

// ClientNewFunc is function for the client
type ClientNewFunc func(Options) ProtocolClient

var rpcClientPlugins = make(map[string]ClientNewFunc)

// GetClientNewFunc is to get the client
func GetClientNewFunc(name string) (ClientNewFunc, error) {
	f := rpcClientPlugins[name]
	if f == nil {
		return nil, fmt.Errorf("Don't have client plugin %s", name)
	}
	return f, nil
}

// InstallPlugin is plugin for the new function
func InstallPlugin(protocol string, f ClientNewFunc) {
	log.Printf("Install client plugin, protocol=%s", protocol)
	rpcClientPlugins[protocol] = f
}
