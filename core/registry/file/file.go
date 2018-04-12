package file

import (
	"fmt"

	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/registry/servicecenter"

	"github.com/ServiceComb/go-sc-client/model"
)

// constant string for file
const (
	Name = "file"
)

// Registrator struct represents file parameters
type Registrator struct {
	Name           string
	registryClient *fileClient
	opts           Options
}

// Close close the file
func (f *Registrator) Close() error {
	return nil
}

// RegisterServiceInstance register service instance
func (f *Registrator) RegisterServiceInstance(sid string, instance *registry.MicroServiceInstance) (string, error) {
	return "", nil
}

// RegisterService register service
func (f *Registrator) RegisterService(microservice *registry.MicroService) (string, error) {
	return "", nil
}

// RegisterServiceAndInstance register service and instance
func (f *Registrator) RegisterServiceAndInstance(microService *registry.MicroService, instance *registry.MicroServiceInstance) (string, string, error) {
	return "", "", nil
}

// Heartbeat check heartbeat of micro-service instance
func (f *Registrator) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	return true, nil
}

// AddDependencies add dependencies
func (f *Registrator) AddDependencies(request *registry.MicroServiceDependency) error {
	return nil
}

// UnRegisterMicroServiceInstance unregister micro-service instances
func (f *Registrator) UnRegisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error {
	return nil
}

// UpdateMicroServiceInstanceStatus update micro-service instance status
func (f *Registrator) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error {
	return nil
}

// UpdateMicroServiceProperties update micro-service properities
func (f *Registrator) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	return nil
}

// UpdateMicroServiceInstanceProperties update micro-service instance properities
func (f *Registrator) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error {
	return nil
}

//AddSchemas add schema
func (f *Registrator) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	return nil
}

// Discovery struct represents file service
type Discovery struct {
	Name           string
	registryClient *fileClient
	opts           Options
}

// Close close the file
func (f *Discovery) Close() error {
	return nil
}

// GetMicroServiceID get micro-service id
func (f *Discovery) GetMicroServiceID(appID, microServiceName, version, env string) (string, error) {
	return "helloService", nil
}

// GetAllMicroServices get all microservices
func (f *Discovery) GetAllMicroServices() ([]*registry.MicroService, error) {
	return []*registry.MicroService{
		{},
	}, nil
}

// GetAllApplications get all applications
func (f *Discovery) GetAllApplications() ([]string, error) {
	return []string{}, nil
}

// GetMicroService get micro-service
func (f *Discovery) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	return &registry.MicroService{
		ServiceID:   "helloService",
		ServiceName: "helloService",
		Status:      "UP",
		Metadata:    map[string]string{},
	}, nil
}

// GetMicroServiceInstances get micro-service instances
func (f *Discovery) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	return []*registry.MicroServiceInstance{}, nil
}

// WatchMicroService watch micro-service
func (f *Discovery) WatchMicroService(selfMicroServiceID string, callback func(*model.MicroServiceInstanceChangedEvent)) {
	return
}

// AutoSync auto sync
func (f *Discovery) AutoSync() {

}

// FindMicroServiceInstances find micro-service instances
func (f *Discovery) FindMicroServiceInstances(consumerID, appID, microServiceName, version, env string) ([]*registry.MicroServiceInstance, error) {
	providerInstances, err := f.registryClient.FindMicroServiceInstances(microServiceName)
	if err != nil {
		return nil, fmt.Errorf("FindMicroServiceInstances failed, err: %s", err)
	}
	instances := filterInstances(providerInstances)

	return instances, nil
}

// filterInstances filter instances
func filterInstances(providerInstances []*model.MicroServiceInstance) []*registry.MicroServiceInstance {
	instances := make([]*registry.MicroServiceInstance, 0)
	for _, ins := range providerInstances {
		msi := servicecenter.ToMicroServiceInstance(ins)
		instances = append(instances, msi)
	}
	return instances
}

// newFileRegistry new file registry
func newFileRegistry(options registry.Options) registry.Registrator {
	fileOption := Options{}
	fileOption.Addrs = options.Addrs
	f := &fileClient{}
	f.Initialize(fileOption)

	return &Registrator{
		Name:           Name,
		registryClient: f,
		opts:           fileOption,
	}
}
func newDiscovery(options registry.Options) registry.ServiceDiscovery {
	fileOption := Options{}
	fileOption.Addrs = options.Addrs
	f := &fileClient{}
	f.Initialize(fileOption)

	return &Discovery{
		Name:           Name,
		registryClient: f,
		opts:           fileOption,
	}
}

// init install plugin of new file registry
func init() {
	registry.InstallRegistrator(Name, newFileRegistry)
	registry.InstallServiceDiscovery(Name, newDiscovery)
}
