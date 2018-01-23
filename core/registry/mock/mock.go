package mock

import (
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-sc-client/model"
	"github.com/stretchr/testify/mock"
)

// RegistryMock struct for register mock
type RegistryMock struct {
	mock.Mock
}

// Close to close the registry mock
func (m *RegistryMock) Close() error {
	return nil
}

// RegisterServiceInstance register service center instance
func (m *RegistryMock) RegisterServiceInstance(sid string, instance *registry.MicroServiceInstance) (string, error) {
	return "", nil
}

// RegisterService register service
func (m *RegistryMock) RegisterService(microservice *registry.MicroService) (string, error) {
	return "", nil
}

// RegisterServiceAndInstance register service and instance
func (m *RegistryMock) RegisterServiceAndInstance(microService *registry.MicroService, instance *registry.MicroServiceInstance) (string, string, error) {
	return "", "", nil
}

// Heartbeat heart beat
func (m *RegistryMock) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	args := m.Called(microServiceID, microServiceInstanceID)
	return args.Bool(0), args.Error(1)
}

// AddDependencies add dependencies
func (m *RegistryMock) AddDependencies(request *registry.MicroServiceDependency) error {
	return nil
}

// GetMicroServiceID get micro-service id
func (m *RegistryMock) GetMicroServiceID(appID, microServiceName, version string) (string, error) {
	args := m.Called(appID, microServiceName, version)
	return args.String(0), args.Error(1)
}

// GetAllMicroServices get all microservices
func (m *RegistryMock) GetAllMicroServices() ([]*registry.MicroService, error) {
	return []*registry.MicroService{
		{},
	}, nil
}

// GetAllApplications get all applications
func (m *RegistryMock) GetAllApplications() ([]string, error) {
	return []string{}, nil
}

// GetMicroService get micro service
func (m *RegistryMock) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	return &registry.MicroService{}, nil
}

// GetMicroServiceInstances get micro-service instances
func (m *RegistryMock) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	args := m.Called(consumerID, providerID)
	return args.Get(0).([]*registry.MicroServiceInstance), args.Error(1)
}

// GetMicroServicesByInterface get services by interface
func (m *RegistryMock) GetMicroServicesByInterface(interfaceName string) []*registry.MicroService {
	return []*registry.MicroService{
		{},
	}
}

// GetSchemaContentByInterface get schema content by interface
func (m *RegistryMock) GetSchemaContentByInterface(interfaceName string) registry.SchemaContent {
	return registry.SchemaContent{}
}

// GetSchemaContentByServiceName get schema content by service name
func (m *RegistryMock) GetSchemaContentByServiceName(svcName, version, appID, env string) []*registry.SchemaContent {
	return []*registry.SchemaContent{}
}

//FindMicroServiceInstances find micro-service instances
func (m *RegistryMock) FindMicroServiceInstances(consumerID, appID, microServiceName, version string) ([]*registry.MicroServiceInstance, error) {
	args := m.Called(consumerID, appID, microServiceName, version)
	return args.Get(0).([]*registry.MicroServiceInstance), args.Error(1)
}

// UnregisterMicroServiceInstance unregistered micro-service instance
func (m *RegistryMock) UnregisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error {
	return nil
}

// WatchMicroService watch micro-service
func (m *RegistryMock) WatchMicroService(selfMicroServiceID string, callback func(*model.MicroServiceInstanceChangedEvent)) {
	return
}

// UpdateMicroServiceInstanceStatus update micro-service instance status
func (m *RegistryMock) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error {
	return nil
}

// UpdateMicroServiceProperties update micro-service properties
func (m *RegistryMock) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	return nil
}

// UpdateMicroServiceInstanceProperties update micro-service instance properties
func (m *RegistryMock) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error {
	return nil
}

// String returns empty string
func (m *RegistryMock) String() string {
	return ""
}

// AutoSync auto sync
func (m *RegistryMock) AutoSync() {}

// AddSchemas add schemas
func (m *RegistryMock) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	return nil
}
