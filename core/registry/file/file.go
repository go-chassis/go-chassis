package file

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/registry/servicecenter"
	"github.com/ServiceComb/go-sc-client/model"
)

// constant string for file
const (
	Name = "File"
)

// File struct represents file parameters
type File struct {
	Name           string
	registryClient *fileClient
	opts           Options
}

// Close close the file
func (f *File) Close() error {
	return nil
}

// RegisterServiceInstance register service instance
func (f *File) RegisterServiceInstance(sid string, instance *registry.MicroServiceInstance) (string, error) {
	return "", nil
}

// RegisterService register service
func (f *File) RegisterService(microservice *registry.MicroService) (string, error) {
	return "", nil
}

// RegisterServiceAndInstance register service and instance
func (f *File) RegisterServiceAndInstance(microService *registry.MicroService, instance *registry.MicroServiceInstance) (string, string, error) {
	return "", "", nil
}

// Heartbeat check heartbeat of micro-service instance
func (f *File) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	return true, nil
}

// AddDependencies add dependencies
func (f *File) AddDependencies(request *registry.MicroServiceDependency) error {
	return nil
}

// GetMicroServiceID get micro-service id
func (f *File) GetMicroServiceID(appID, microServiceName, version string) (string, error) {
	return "helloService", nil
}

// GetAllMicroServices get all microservices
func (f *File) GetAllMicroServices() ([]*registry.MicroService, error) {
	return []*registry.MicroService{
		{},
	}, nil
}

// GetAllApplications get all applications
func (f *File) GetAllApplications() ([]string, error) {
	return []string{}, nil
}

// GetMicroService get micro-service
func (f *File) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	return &registry.MicroService{
		ServiceID:   "helloService",
		ServiceName: "helloService",
		Status:      "UP",
		Metadata:    map[string]string{},
	}, nil
}

// GetMicroServiceInstances get micro-service instances
func (f *File) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	return []*registry.MicroServiceInstance{}, nil
}

// UnregisterMicroServiceInstance unregister micro-service instances
func (f *File) UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error {
	return nil
}

// WatchMicroService watch micro-service
func (f *File) WatchMicroService(selfMicroServiceID string, callback func(*model.MicroServiceInstanceChangedEvent)) {
	return
}

// UpdateMicroServiceInstanceStatus update micro-service instance status
func (f *File) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error {
	return nil
}

// UpdateMicroServiceProperties update micro-service properities
func (f *File) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	return nil
}

// UpdateMicroServiceInstanceProperties update micro-service instance properities
func (f *File) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error {
	return nil
}

// String returns empty string
func (f *File) String() string {
	return ""
}

// AutoSync auto sync
func (f *File) AutoSync() {

}

// AddSchemas add schemas
func (f *File) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	return nil
}

// GetMicroServicesByInterface get micro-services by interface
func (f *File) GetMicroServicesByInterface(interfaceName string) (services []*registry.MicroService) {
	return services
}

// GetSchemaContentByInterface get schema content by interface
func (f *File) GetSchemaContentByInterface(interfaceName string) (content registry.SchemaContent) {
	return content
}

// GetSchemaContentByServiceName get schema content by service name
func (f *File) GetSchemaContentByServiceName(svcName, version, appID, env string) (content []*registry.SchemaContent) {
	return content
}

// FindMicroServiceInstances find micro-service instances
func (f *File) FindMicroServiceInstances(consumerID, appID, microServiceName, version string) ([]*registry.MicroServiceInstance, error) {
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
func newFileRegistry(opts ...registry.Option) registry.Registry {
	var options registry.Options
	for _, o := range opts {
		o(&options)
	}
	fileOption := Options{}
	fileOption.Addrs = options.Addrs
	f := &fileClient{}
	f.Initialize(fileOption)

	return &File{
		Name:           Name,
		registryClient: f,
		opts:           fileOption,
	}
}

// init install plugin of new file registry
func init() {
	registry.InstallPlugin(Name, newFileRegistry)
}
