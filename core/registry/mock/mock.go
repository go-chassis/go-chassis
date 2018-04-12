package mock

import (
	"github.com/ServiceComb/go-chassis/core/registry"

	"github.com/ServiceComb/go-sc-client/model"
	"github.com/stretchr/testify/mock"
)

// RegistratorMock struct for register mock
type RegistratorMock struct {
	mock.Mock
}

// Close to close the registry mock
func (m *RegistratorMock) Close() error {
	return nil
}

// RegisterServiceInstance register service center instance
func (m *RegistratorMock) RegisterServiceInstance(sid string, instance *registry.MicroServiceInstance) (string, error) {
	return "", nil
}

// RegisterService register service
func (m *RegistratorMock) RegisterService(microservice *registry.MicroService) (string, error) {
	return "", nil
}

// RegisterServiceAndInstance register service and instance
func (m *RegistratorMock) RegisterServiceAndInstance(microService *registry.MicroService, instance *registry.MicroServiceInstance) (string, string, error) {
	return "", "", nil
}

// Heartbeat heart beat
func (m *RegistratorMock) Heartbeat(microServiceID, microServiceInstanceID string) (bool, error) {
	args := m.Called(microServiceID, microServiceInstanceID)
	return args.Bool(0), args.Error(1)
}

// AddDependencies add dependencies
func (m *RegistratorMock) AddDependencies(request *registry.MicroServiceDependency) error {
	return nil
}

// UpdateMicroServiceInstanceStatus update micro-service instance status
func (m *RegistratorMock) UpdateMicroServiceInstanceStatus(microServiceID, microServiceInstanceID, status string) error {
	return nil
}

// UpdateMicroServiceProperties update micro-service properties
func (m *RegistratorMock) UpdateMicroServiceProperties(microServiceID string, properties map[string]string) error {
	return nil
}

// UpdateMicroServiceInstanceProperties update micro-service instance properties
func (m *RegistratorMock) UpdateMicroServiceInstanceProperties(microServiceID, microServiceInstanceID string, properties map[string]string) error {
	return nil
}

// UnRegisterMicroServiceInstance unregistered micro-service instance
func (m *RegistratorMock) UnRegisterMicroServiceInstance(microServiceID, microServiceInstanceID string) error {
	return nil
}

// AddSchemas add schemas
func (m *RegistratorMock) AddSchemas(microServiceID, schemaName, schemaInfo string) error {
	return nil
}

// DiscoveryMock struct for disco mock
type DiscoveryMock struct {
	mock.Mock
}

// GetMicroServiceID get micro-service id
func (m *DiscoveryMock) GetMicroServiceID(appID, microServiceName, version, env string) (string, error) {
	args := m.Called(appID, microServiceName, version, env)
	return args.String(0), args.Error(1)
}

// GetAllMicroServices get all microservices
func (m *DiscoveryMock) GetAllMicroServices() ([]*registry.MicroService, error) {
	return []*registry.MicroService{
		{},
	}, nil
}

// GetAllApplications get all applications
func (m *DiscoveryMock) GetAllApplications() ([]string, error) {
	return []string{}, nil
}

// GetMicroService get micro service
func (m *DiscoveryMock) GetMicroService(microServiceID string) (*registry.MicroService, error) {
	return &registry.MicroService{}, nil
}

// GetMicroServiceInstances get micro-service instances
func (m *DiscoveryMock) GetMicroServiceInstances(consumerID, providerID string) ([]*registry.MicroServiceInstance, error) {
	args := m.Called(consumerID, providerID)
	return args.Get(0).([]*registry.MicroServiceInstance), args.Error(1)
}

//FindMicroServiceInstances find micro-service instances
func (m *DiscoveryMock) FindMicroServiceInstances(consumerID, appID, microServiceName, version, env string) ([]*registry.MicroServiceInstance, error) {
	args := m.Called(consumerID, appID, microServiceName, version, env)
	return args.Get(0).([]*registry.MicroServiceInstance), args.Error(1)
}

// WatchMicroService watch micro-service
func (m *DiscoveryMock) WatchMicroService(selfMicroServiceID string, callback func(*model.MicroServiceInstanceChangedEvent)) {
	return
}

// AutoSync auto sync
func (m *DiscoveryMock) AutoSync() {}

// Close mock
func (m *DiscoveryMock) Close() error { return nil }
