package registry

import (
	"fmt"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/go-mesh/openlogging"
)

var sdFunc = make(map[string]func(opts Options) ServiceDiscovery)

//InstallServiceDiscovery install service discovery client
func InstallServiceDiscovery(name string, f func(opts Options) ServiceDiscovery) {
	sdFunc[name] = f
	openlogging.Info("Installed service discovery plugin: " + name)
}

//NewDiscovery create discovery service
func NewDiscovery(name string, opts Options) (ServiceDiscovery, error) {
	f := sdFunc[name]
	if f == nil {
		return nil, fmt.Errorf("no service discovery plugin: %s", name)
	}
	return f(opts), nil
}

//ServiceDiscovery fetch service and instances from remote or local
type ServiceDiscovery interface {
	GetMicroService(microServiceID string) (*MicroService, error)
	FindMicroServiceInstances(consumerID, microServiceName string, tags utiltags.Tags) ([]*MicroServiceInstance, error)
	AutoSync()
	Close() error
}

//DefaultServiceDiscoveryService supplies service discovery
var DefaultServiceDiscoveryService ServiceDiscovery

func enableServiceDiscovery(opts Options) error {
	if config.GetServiceDiscoveryDisable() {
		openlogging.Warn("discovery is disabled")
		return nil
	}

	t := config.GetServiceDiscoveryType()
	if t == "" {
		t = DefaultServiceDiscoveryPlugin
	}
	f := sdFunc[t]
	if f == nil {
		panic("No service discovery plugin")
	}
	var err error
	DefaultServiceDiscoveryService, err = NewDiscovery(t, opts)
	if err != nil {
		return err
	}

	DefaultServiceDiscoveryService.AutoSync()

	openlogging.GetLogger().Infof("Enable %s service discovery.", t)
	return nil
}
