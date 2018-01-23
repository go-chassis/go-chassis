package registry

import (
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/schema"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/metadata"
	"github.com/ServiceComb/go-sc-client/model"
	"os"
)

// microServiceDependencies micro-service dependencies
var microServiceDependencies *MicroServiceDependency

// InstanceEndpoints instance endpoints
var InstanceEndpoints map[string]string

// RegisterMicroservice register micro-service
func RegisterMicroservice() error {
	microServiceDependencies = &MicroServiceDependency{}
	service := config.MicroserviceDefinition
	schemas, err := schema.GetSchemaIDs(service.ServiceDescription.Name)
	if err != nil {
		lager.Logger.Warnf(err, "Get schemas failed, microservice = %s.", service.ServiceDescription.Name)
		schemas = make([]string, 0)
	}
	if service.ServiceDescription.Level == "" {
		service.ServiceDescription.Level = "BACK"
	}
	framework := metadata.NewFramework()

	microservice := &MicroService{
		AppID:       config.GlobalDefinition.AppID,
		ServiceName: service.ServiceDescription.Name,
		Version:     service.ServiceDescription.Version,
		Status:      model.MicorserviceUp,
		Level:       service.ServiceDescription.Level,
		Schemas:     schemas,
		Framework: &Framework{
			Version: framework.Version,
			Name:    framework.Name,
		},
		RegisterBy: framework.Register,
	}
	lager.Logger.Infof("Framework registered is [ %s:%s ]", framework.Name, framework.Version)
	lager.Logger.Infof("Microservice registered by [ %s ]", framework.Register)

	sid, err := RegistryService.RegisterService(microservice)
	config.SelfServiceID = sid
	if err != nil {
		lager.Logger.Errorf(err, "Register microservice [%s] failed", microservice.ServiceName)
		return err
	}
	lager.Logger.Warnf(nil, "Register microservice [%s] success", microservice.ServiceName)

	for _, schemaID := range schemas {
		schemaInfo := schema.DefaultSchemaIDsMap[schemaID]
		RegistryService.AddSchemas(sid, schemaID, schemaInfo)
	}
	if service.ServiceDescription.Properties == nil {
		service.ServiceDescription.Properties = make(map[string]string)
	}

	//update metadata
	if config.GlobalDefinition.Cse.Service.Registry.Scope == "full" {
		service.ServiceDescription.Properties["allowCrossApp"] = "true"
	} else {
		service.ServiceDescription.Properties["allowCrossApp"] = "false"
	}
	if err := RegistryService.UpdateMicroServiceProperties(sid, service.ServiceDescription.Properties); err != nil {
		lager.Logger.Errorf(err, "Update microservice properties failed, serviceID = %s.", sid)
		return err
	}
	lager.Logger.Debugf("Update microservice properties success, serviceID = %s.", sid)

	return refreshDependency(microservice)
}

// refreshDependency refresh dependency
func refreshDependency(service *MicroService) error {
	providersDependencyMicroService := make([]*MicroService, 0)
	if len(config.GlobalDefinition.Cse.References) == 0 {
		lager.Logger.Info("Don't need add dependency")
		return nil
	}
	for k, v := range config.GlobalDefinition.Cse.References {
		providerDependencyMicroService := &MicroService{
			AppID:       config.GlobalDefinition.AppID,
			ServiceName: k,
			Version:     v.Version,
		}
		providersDependencyMicroService = append(providersDependencyMicroService, providerDependencyMicroService)
	}
	microServiceDependency := &MicroServiceDependency{
		Consumer:  service,
		Providers: providersDependencyMicroService,
	}
	microServiceDependencies = microServiceDependency

	return RegistryService.AddDependencies(microServiceDependencies)
}

// RegisterMicroserviceInstances register micro-service instances
func RegisterMicroserviceInstances() error {
	lager.Logger.Warn("Start to register instances.", nil)
	hostname, err := os.Hostname()
	if err != nil {
		lager.Logger.Error("Get hostname failed.", err)
		return err
	}
	stage := config.Stage
	service := config.MicroserviceDefinition
	sid, err := RegistryService.GetMicroServiceID(config.GlobalDefinition.AppID, service.ServiceDescription.Name, service.ServiceDescription.Version)
	if err != nil {
		lager.Logger.Errorf(err, "Get service failed, key: %s:%s:%s",
			config.GlobalDefinition.AppID,
			service.ServiceDescription.Name,
			service.ServiceDescription.Version)
		return err
	}
	eps := MakeEndpointMap(config.GlobalDefinition.Cse.Protocols)
	if InstanceEndpoints != nil {
		eps = InstanceEndpoints
	}

	microServiceInstance := &MicroServiceInstance{
		EndpointsMap: eps,
		HostName:     hostname,
		Status:       model.MSInstanceUP,
		Environment:  stage,
		Metadata:     map[string]string{"nodeIP": config.NodeIP},
	}

	var dInfo = new(DataCenterInfo)
	if config.GlobalDefinition.DataCenter.Name != "" && config.GlobalDefinition.DataCenter.AvailableZone != "" {
		dInfo.Name = config.GlobalDefinition.DataCenter.Name
		dInfo.Region = config.GlobalDefinition.DataCenter.Name
		dInfo.AvailableZone = config.GlobalDefinition.DataCenter.AvailableZone
		microServiceInstance.DataCenterInfo = dInfo
	}

	instanceID, err := RegistryService.RegisterServiceInstance(sid, microServiceInstance)
	if err != nil {
		lager.Logger.Errorf(err, "Register instance failed, serviceID: %s.", sid)
		return err
	}
	if service.ServiceDescription.InstanceProperties != nil {
		if err := RegistryService.UpdateMicroServiceInstanceProperties(sid, instanceID, service.ServiceDescription.InstanceProperties); err != nil {
			lager.Logger.Errorf(nil, "UpdateMicroServiceInstanceProperties failed, microServiceID/instanceID = %s/%s.", sid, instanceID)
			return err
		}
		lager.Logger.Debugf("UpdateMicroServiceInstanceProperties success, microServiceID/instanceID = %s/%s.", sid, instanceID)
	}

	value, _ := SelfInstancesCache.Get(microServiceInstance.ServiceID)
	instanceIDs, _ := value.([]string)
	var isRepeat bool
	for _, va := range instanceIDs {
		if va == instanceID {
			isRepeat = true
		}
	}
	if !isRepeat {
		instanceIDs = append(instanceIDs, instanceID)
	}
	SelfInstancesCache.Set(sid, instanceIDs, 0)
	lager.Logger.Warnf(nil, "Register instance success, serviceID/instanceID: %s/%s.", sid, instanceID)
	return nil
}
