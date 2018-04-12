package servicecenter

import (
	"github.com/ServiceComb/go-chassis/core/registry"

	"github.com/ServiceComb/go-sc-client"
	"github.com/ServiceComb/go-sc-client/model"
)

// ToMicroService assign sc micro-service to go chassis micro-service
func ToMicroService(scs *model.MicroService) *registry.MicroService {
	cs := &registry.MicroService{}
	cs.ServiceID = scs.ServiceID
	cs.ServiceName = scs.ServiceName
	cs.Version = scs.Version
	cs.AppID = scs.AppID
	cs.Metadata = scs.Properties
	cs.Schemas = scs.Schemas
	cs.Level = scs.Level
	cs.Status = scs.Status
	if scs.Framework != nil {
		cs.Framework = &registry.Framework{
			Name:    scs.Framework.Name,
			Version: scs.Framework.Version,
		}
	}
	return cs
}

// ToSCService assign go chassis micro-service to the sc micro-service
func ToSCService(cs *registry.MicroService) *model.MicroService {
	scs := &model.MicroService{}
	scs.ServiceID = cs.ServiceID
	scs.ServiceName = cs.ServiceName
	scs.Version = cs.Version
	scs.AppID = cs.AppID
	scs.Environment = cs.Environment
	scs.Properties = cs.Metadata
	scs.Schemas = cs.Schemas
	scs.Level = cs.Level
	scs.Status = cs.Status
	if cs.Framework != nil {
		scs.Framework = &model.Framework{}
		scs.Framework.Version = cs.Framework.Version
		scs.Framework.Name = cs.Framework.Name
	}
	scs.RegisterBy = cs.RegisterBy
	return scs
}

// ToMicroServiceInstance assign model micro-service instance parameters to registry micro-service instance parameters
func ToMicroServiceInstance(ins *model.MicroServiceInstance) *registry.MicroServiceInstance {
	msi := &registry.MicroServiceInstance{}
	m, p := registry.GetProtocolMap(ins.Endpoints)
	msi.InstanceID = ins.InstanceID
	msi.EndpointsMap = m
	msi.DefaultEndpoint = m[p]
	msi.DefaultProtocol = p
	msi.Metadata = ins.Properties
	if ins.DataCenterInfo != nil {
		msi.DataCenterInfo = &registry.DataCenterInfo{}
		msi.DataCenterInfo.Name = ins.DataCenterInfo.Name
		msi.DataCenterInfo.AvailableZone = ins.DataCenterInfo.AvailableZone
		msi.DataCenterInfo.Region = ins.DataCenterInfo.Region
	}

	return msi
}

// ToSCInstance assign registry micro-service instance parameters to model micro-service instance parameters
func ToSCInstance(msi *registry.MicroServiceInstance) *model.MicroServiceInstance {
	si := &model.MicroServiceInstance{}
	eps := registry.GetProtocolList(msi.EndpointsMap)
	si.InstanceID = msi.InstanceID
	si.Endpoints = eps
	si.Properties = msi.Metadata
	si.HostName = msi.HostName
	si.Status = msi.Status
	if msi.DataCenterInfo != nil {
		si.DataCenterInfo = &model.DataCenterInfo{}
		si.DataCenterInfo.Name = msi.DataCenterInfo.Name
		si.DataCenterInfo.AvailableZone = msi.DataCenterInfo.AvailableZone
		si.DataCenterInfo.Region = msi.DataCenterInfo.Region
	}

	return si
}

// ToSCDependency assign registry micro-service dependencies to model micro-service dependencies
func ToSCDependency(dep *registry.MicroServiceDependency) *model.MircroServiceDependencyRequest {
	scDep := &model.MircroServiceDependencyRequest{
		Dependencies: make([]*model.MicroServiceDependency, 1),
	}
	scDep.Dependencies[0] = &model.MicroServiceDependency{}
	scDep.Dependencies[0].Consumer = &model.DependencyMicroService{
		AppID:       dep.Consumer.AppID,
		ServiceName: dep.Consumer.ServiceName,
		Version:     dep.Consumer.Version,
	}
	for _, p := range dep.Providers {
		scP := &model.DependencyMicroService{
			AppID:       p.AppID,
			ServiceName: p.ServiceName,
			Version:     p.Version,
		}
		scDep.Dependencies[0].Providers = append(scDep.Dependencies[0].Providers, scP)
	}
	return scDep
}

//ToSCOptions convert registry opstions into sc client options
func ToSCOptions(options registry.Options) client.Options {
	sco := client.Options{}
	sco.Timeout = options.Timeout
	sco.TLSConfig = options.TLSConfig
	sco.Addrs = options.Addrs
	sco.Compressed = options.Compressed
	sco.ConfigTenant = options.Tenant
	sco.EnableSSL = options.EnableSSL
	sco.Verbose = options.Verbose
	sco.Version = options.Version
	return sco
}
