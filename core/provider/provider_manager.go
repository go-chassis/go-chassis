package provider

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/core/lager"
	"log"
)

// plugin name and schemas map
var providerPlugins = make(map[string]func(string) Provider)

// microservice name and schemas map
var providers = make(map[string]Provider)
var defaultProviderFunc = NewProvider

//TODO locks

// InstallProviderPlugin install provider plugin
func InstallProviderPlugin(pluginName string, newFunc func(string) Provider) {
	log.Printf("Install Provider Plugin, name=%s", pluginName)
	providerPlugins[pluginName] = newFunc
}

// todo: return error.

// RegisterProvider register provider
func RegisterProvider(pluginName string, microserviceName string) Provider {
	pFunc, exist := providerPlugins[pluginName]
	if !exist {
		lager.Logger.Errorf(nil, "provider type %s is not exist.", pluginName)
		return nil
	}
	p := pFunc(microserviceName)
	lager.Logger.Debugf("registered provider for service [%s]", microserviceName)
	RegisterCustomProvider(microserviceName, p)
	return p

}

// RegisterCustomProvider register customer provider
func RegisterCustomProvider(microserviceName string, p Provider) {
	if providers[microserviceName] != nil {
		lager.Logger.Warnf(nil, "Can not replace Provider,since it is not nil")
		return
	}
	providers[microserviceName] = p
}

// GetProvider get provider
func GetProvider(microserviceName string) (Provider, error) {
	p, exist := providers[microserviceName]
	if !exist {
		return nil, fmt.Errorf("Service [%s] doesn't have provider", microserviceName)
	}
	return p, nil
}

// RegisterSchemaWithName register schema with name
func RegisterSchemaWithName(microserviceName string, schemaID string, schema interface{}) error {
	p, exist := providers[microserviceName]
	if !exist {
		return fmt.Errorf("service: %s do not exist", microserviceName)
	}
	return p.RegisterName(schemaID, schema)
}

// RegisterSchema register schema
func RegisterSchema(microserviceName string, schema interface{}) (string, error) {
	p := providers[microserviceName]
	if p == nil {
		return "", fmt.Errorf("[%s] Provider is not exist ", microserviceName)
	}
	return p.Register(schema)
}

// GetOperation get operation
func GetOperation(microserviceName string, schemaID string, operationID string) (Operation, error) {
	p, ok := providers[microserviceName]
	if !ok {
		return nil, fmt.Errorf("MicroService [%s] doesn't exist", microserviceName)
	}
	return p.GetOperation(schemaID, operationID)
}
