package registry

import (
	"errors"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/schema"
	"github.com/go-chassis/go-chassis/core/metadata"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-mesh/openlogging"
)

var errEmptyServiceIDFromRegistry = errors.New("got empty serviceID from registry")

// InstanceEndpoints instance endpoints
var InstanceEndpoints = make(map[string]string)

// RegisterService register micro-service
func RegisterService() error {
	service := config.MicroserviceDefinition
	if e := service.Environment; e != "" {
		openlogging.GetLogger().Infof("Microservice environment: [%s]", e)
	}
	var err error
	runtime.Schemas, err = schema.GetSchemaIDs(runtime.ServiceName)
	if err != nil || len(runtime.Schemas) == 0 {
		openlogging.GetLogger().Warnf("no schemas file for microservice [%s].", runtime.ServiceName)
		runtime.Schemas = make([]string, 0)
		// from yaml setting
		if len(service.Schemas) != 0 {
			runtime.Schemas = service.Schemas
		}
	}
	if service.ServicesStatus == "" {
		service.ServicesStatus = common.DefaultStatus
	}
	if service.Properties == nil {
		service.Properties = make(map[string]string)
	}
	framework := metadata.NewFramework()

	svcPaths := service.ServicePaths
	var regpaths []ServicePath
	for _, svcPath := range svcPaths {
		var regpath ServicePath
		regpath.Path = svcPath.Path
		regpath.Property = svcPath.Property
		regpaths = append(regpaths, regpath)
	}
	microservice := &MicroService{
		ServiceID:   runtime.ServiceID,
		AppID:       runtime.App,
		ServiceName: service.Name,
		Version:     service.Version,
		Paths:       regpaths,
		Environment: service.Environment,
		Status:      service.ServicesStatus,
		Schemas:     runtime.Schemas,
		Framework: &Framework{
			Version: framework.Version,
			Name:    framework.Name,
		},
		RegisterBy: framework.Register,
		Metadata:   make(map[string]string),
		// TODO allows to customize microservice alias
		Alias: "",
	}
	//update metadata
	if len(microservice.Alias) == 0 {
		// if the microservice is allowed to be called by consumers with different appId,
		// this means that the governance configuration of the consumer side needs to
		// support key format with appid, like 'servicecomb.loadbalance.{alias}.strategy.name'.
		microservice.Alias = microservice.AppID + ":" + microservice.ServiceName
	}
	if config.GetRegistratorScope() == common.ScopeFull {
		microservice.Metadata["allowCrossApp"] = common.TRUE
		service.Properties["allowCrossApp"] = common.TRUE
	} else {
		service.Properties["allowCrossApp"] = common.FALSE
	}
	openlogging.GetLogger().Debugf("update micro service properties%v", service.Properties)
	openlogging.GetLogger().Infof("framework registered is [ %s:%s ]", framework.Name, framework.Version)
	openlogging.GetLogger().Infof("micro service registered by [ %s ]", framework.Register)

	sid, err := DefaultRegistrator.RegisterService(microservice)
	if err != nil {
		openlogging.GetLogger().Errorf("register service [%s] failed: %s", microservice.ServiceName, err)
		return err
	}
	if sid == "" {
		openlogging.Error(errEmptyServiceIDFromRegistry.Error())
		return errEmptyServiceIDFromRegistry
	}
	runtime.ServiceID = sid
	openlogging.GetLogger().Infof("register service success:[%s] ", runtime.ServiceID)

	return nil
}

// RegisterServiceInstances register micro-service instances
func RegisterServiceInstances() error {
	var err error
	service := config.MicroserviceDefinition
	runtime.Schemas, err = schema.GetSchemaIDs(service.Name)
	if err != nil || len(runtime.Schemas) == 0 {
		runtime.Schemas = make([]string, 0)
		// from yaml setting
		if len(service.Schemas) != 0 {
			runtime.Schemas = service.Schemas
		}
	}

	for _, schemaID := range runtime.Schemas {
		schemaInfo := schema.GetContent(schemaID)
		err := DefaultRegistrator.AddSchemas(runtime.ServiceID, schemaID, schemaInfo)
		if err != nil {
			openlogging.Warn("upload contract to registry failed: " + err.Error())
		}
		openlogging.Debug("upload contract to registry, " + schemaID)
	}
	openlogging.Debug("start to register instance.")
	eps, err := MakeEndpointMap(config.GlobalDefinition.ServiceComb.Protocols)
	if err != nil {
		return err
	}
	openlogging.GetLogger().Infof("service support protocols %v", config.GlobalDefinition.ServiceComb.Protocols)
	if len(InstanceEndpoints) != 0 {
		eps = make(map[string]*Endpoint, len(InstanceEndpoints))
		for m, ep := range InstanceEndpoints {
			epObj, err := NewEndPoint(ep)
			if err != nil {
				openlogging.GetLogger().Errorf("parser instance protocol %s endpoint error %s", m, err)
				continue
			}
			eps[m] = epObj
		}
	}
	if service.ServicesStatus == "" {
		service.ServicesStatus = common.DefaultStatus
	}
	microServiceInstance := &MicroServiceInstance{
		InstanceID:   runtime.InstanceID,
		EndpointsMap: eps,
		HostName:     runtime.HostName,
		Status:       service.ServicesStatus,
		Metadata:     map[string]string{"nodeIP": runtime.NodeIP},
	}
	var dInfo = new(DataCenterInfo)
	if config.GlobalDefinition.DataCenter.Name != "" && config.GlobalDefinition.DataCenter.AvailableZone != "" {
		dInfo.Name = config.GlobalDefinition.DataCenter.Name
		dInfo.Region = config.GlobalDefinition.DataCenter.Name
		dInfo.AvailableZone = config.GlobalDefinition.DataCenter.AvailableZone
		microServiceInstance.DataCenterInfo = dInfo
	}
	instanceID, err := DefaultRegistrator.RegisterServiceInstance(runtime.ServiceID, microServiceInstance)
	if err != nil {
		openlogging.GetLogger().Errorf("register instance failed, serviceID: %s, err %s", runtime.ServiceID, err.Error())
		return err
	}
	//Set to runtime
	runtime.InstanceID = instanceID
	runtime.InstanceStatus = runtime.StatusRunning
	if service.InstanceProperties != nil {
		if err := DefaultRegistrator.UpdateMicroServiceInstanceProperties(runtime.ServiceID, instanceID, service.InstanceProperties); err != nil {
			openlogging.GetLogger().Errorf("UpdateMicroServiceInstanceProperties failed, microServiceID/instanceID = %s/%s.", runtime.ServiceID, instanceID)
			return err
		}
		runtime.InstanceMD = service.InstanceProperties
		openlogging.GetLogger().Debugf("UpdateMicroServiceInstanceProperties success, microServiceID/instanceID = %s/%s.", runtime.ServiceID, instanceID)
	}
	openlogging.GetLogger().Infof("register instance success, instanceID: %s.", instanceID)
	return nil
}
