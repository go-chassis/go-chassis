package provider

import (
	"fmt"

	"github.com/go-mesh/openlogging"
)

// plugin name and schemas map
var providerPlugins = make(map[string]func(string) Provider)

// micro service name and schemas map
var providers = make(map[string]Provider)

//TODO locks

// InstallProviderPlugin install provider plugin
func InstallProviderPlugin(pluginName string, newFunc func(string) Provider) {
	openlogging.Info("Install Provider Plugin, name: " + pluginName)
	providerPlugins[pluginName] = newFunc
}

// todo: return error.

// RegisterProvider register provider
func RegisterProvider(pluginName string, microServiceName string) Provider {
	pFunc, exist := providerPlugins[pluginName]
	if !exist {
		openlogging.GetLogger().Errorf("provider type %s is not exist.", pluginName)
		return nil
	}
	p := pFunc(microServiceName)
	openlogging.GetLogger().Debugf("registered provider for service [%s]", microServiceName)
	RegisterCustomProvider(microServiceName, p)
	return p

}

// RegisterCustomProvider register customer provider
func RegisterCustomProvider(microServiceName string, p Provider) {
	if providers[microServiceName] != nil {
		openlogging.GetLogger().Warnf("Can not replace Provider,since it is not nil")
		return
	}
	providers[microServiceName] = p
}

// GetProvider get provider
func GetProvider(microServiceName string) (Provider, error) {
	p, exist := providers[microServiceName]
	if !exist {
		return nil, fmt.Errorf("service [%s] doesn't have provider", microServiceName)
	}
	return p, nil
}

// RegisterSchemaWithName register schema with name
func RegisterSchemaWithName(microServiceName string, schemaID string, schema interface{}) error {
	p, exist := providers[microServiceName]
	if !exist {
		return fmt.Errorf("service: %s do not exist", microServiceName)
	}
	return p.RegisterName(schemaID, schema)
}

// RegisterSchema register schema
func RegisterSchema(microServiceName string, schema interface{}) (string, error) {
	p := providers[microServiceName]
	if p == nil {
		return "", fmt.Errorf("[%s] Provider is not exist ", microServiceName)
	}
	return p.Register(schema)
}

// GetOperation get operation
func GetOperation(microServiceName string, schemaID string, operationID string) (Operation, error) {
	p, ok := providers[microServiceName]
	if !ok {
		return nil, fmt.Errorf("MicroService [%s] doesn't exist", microServiceName)
	}
	return p.GetOperation(schemaID, operationID)
}
