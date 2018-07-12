package kuberegistry

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	corelisters "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// DiscoveryController defines discovery controller for kube registry
type DiscoveryController struct {
	client kubernetes.Interface

	sLister corelisters.ServiceLister
	eLister corelisters.EndpointsLister

	sListerSynced cache.InformerSynced
	eListerSynced cache.InformerSynced
}

// NewDiscoveryController returns new discovery controller
func NewDiscoveryController(
	sInformer coreinformers.ServiceInformer,
	eInformer coreinformers.EndpointsInformer,
	client kubernetes.Interface,
) *DiscoveryController {

	dc := &DiscoveryController{
		client: client,
	}

	sInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: dc.addService,
	})
	eInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: dc.addEndpoints,
	})
	dc.sListerSynced = sInformer.Informer().HasSynced
	dc.eListerSynced = eInformer.Informer().HasSynced
	dc.sLister = sInformer.Lister()
	dc.eLister = eInformer.Lister()
	return dc
}

// Run begins discovery controller
func (dc *DiscoveryController) Run(stop <-chan struct{}) {
	lager.Logger.Info("Starting Discovery Controller")
	if !cache.WaitForCacheSync(stop, dc.sListerSynced, dc.eListerSynced) {
		lager.Logger.Error("Time out waiting for caches to sync", nil)
		return
	}
	lager.Logger.Info("Finish Waiting For Cache Sync")
}

func (dc *DiscoveryController) addService(obj interface{}) {
	svc := obj.(*v1.Service)
	lager.Logger.Infof("Add Service: %s", svc.Name)
}

func (dc *DiscoveryController) addEndpoints(obj interface{}) {
	ep := obj.(*v1.Endpoints)
	lager.Logger.Infof("Add Endpoint: %s", ep.Name)
}

// FindEndpoints returns microservice instances of kube registry
func (dc *DiscoveryController) FindEndpoints(service string, tags registry.Tags) ([]*registry.MicroServiceInstance, error) {
	name, namespace := splitServiceKey(service)
	// TODO: use labels.ToLabelSelector to trans endpoint
	// use cache lister to get specific endpoints or use kubeclient instead
	ep, err := dc.eLister.Endpoints(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return toMicroServiceInstances(ep), nil
}

// GetAllServices returns microservice of kube registry
func (dc *DiscoveryController) GetAllServices() ([]*registry.MicroService, error) {
	microServices, err := dc.sLister.List(labels.Everything())
	if err != nil {
		lager.Logger.Errorf(err, "get all microservices from kube failed")
		return nil, err
	}
	ms := make([]*registry.MicroService, len(microServices))
	for i, s := range microServices {
		ms[i] = toMicroService(s)
	}
	lager.Logger.Debugf("get all microservices success, microservices: %v", microServices)
	return ms, nil
}
