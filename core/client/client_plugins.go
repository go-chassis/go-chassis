package client

import (
	"fmt"
	"log"
)

// NewFunc is function for the client
type NewFunc func(Options) (ProtocolClient, error)

var plugins = make(map[string]NewFunc)

// GetClientNewFunc is to get the client
func GetClientNewFunc(name string) (NewFunc, error) {
	f := plugins[name]
	if f == nil {
		return nil, fmt.Errorf("don't have client plugin %s", name)
	}
	return f, nil
}

// InstallPlugin is plugin for the new function
func InstallPlugin(protocol string, f NewFunc) {
	log.Printf("Install client plugin, protocol=%s", protocol)
	plugins[protocol] = f
}
