package kuberegistry

import (
	"math/rand"
	"time"

	"github.com/ServiceComb/go-chassis/core/registry"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	// KubeRegistry constant string
	KubeRegistry = "kube"
	// DefaultMinResyncPeriod determins the minmum resync period
	DefaultMinResyncPeriod = 1 * time.Second
)

// init initialize the plugin of service center registry
func init() { registry.InstallServiceDiscovery(KubeRegistry, newServiceDiscovery) }

// ServiceDiscovery to represent the object of service center to call the APIs of service center
type ServiceDiscovery struct {
	Controller *DiscoveryController

	Name string
}

// ResyncPeriod returns func to calculate time duration of resync
func ResyncPeriod(options registry.Options) func() time.Duration {
	return func() time.Duration {
		factor := rand.Float64() + 1
		return time.Duration(float64(DefaultMinResyncPeriod.Nanoseconds()) * factor)
	}
}

func createClientOrDie(kubeconfig string) kubernetes.Interface {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic("build config from flags failed" + err.Error())
	}
	client, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic("new client from config failed" + err.Error())
	}
	return client
}

func newServiceDiscovery(options registry.Options) registry.ServiceDiscovery {
	client := createClientOrDie(options.ConfigPath)
	sharedInformers := informers.NewSharedInformerFactory(client, ResyncPeriod(options)())

	controller := NewDiscoveryController(
		sharedInformers.Core().V1().Services(),
		sharedInformers.Core().V1().Endpoints(),
		client,
	)
	stop := make(chan struct{})
	sharedInformers.Start(stop)
	controller.Run(stop)

	return &ServiceDiscovery{
		Name:       KubeRegistry,
		Controller: controller,
	}
}

// GetAllMicroServices Get all MicroService information.
func (r *ServiceDiscovery) GetAllMicroServices() ([]*registry.MicroService, error) {
	return r.Controller.GetAllServices()
}

// FindMicroServiceInstances find micro-service instances
func (r *ServiceDiscovery) FindMicroServiceInstances(consumerID, microServiceName string, tags registry.Tags) ([]*registry.MicroServiceInstance, error) {
	return r.Controller.FindEndpoints(microServiceName, tags)
}

// GetMicroServiceID get microServiceID
func (r *ServiceDiscovery) GetMicroServiceID(appID, microServiceName, version, env string) (string, error) {
	return "", nil
}

// GetMicroServiceInstances return instances
func (r *ServiceDiscovery) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	return nil, nil
}

// GetMicroService return service
func (r *ServiceDiscovery) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	return nil, nil
}

// AutoSync updating the cache manager
func (r *ServiceDiscovery) AutoSync() {}

// Close close all websocket connection
func (r *ServiceDiscovery) Close() error { return nil }
