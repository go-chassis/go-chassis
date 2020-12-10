package servicecenter

import (
	scregistry "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/go-chassis/v2/core/registry"
	"github.com/go-chassis/sc-client"
)

// ToMicroService assign sc micro-service to go chassis micro-service
func ToMicroService(scs *scregistry.MicroService) *registry.MicroService {
	cs := &registry.MicroService{}
	cs.ServiceID = scs.ServiceId
	cs.ServiceName = scs.ServiceName
	cs.Version = scs.Version
	cs.AppID = scs.AppId
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
func ToSCService(cs *registry.MicroService) *scregistry.MicroService {
	scs := &scregistry.MicroService{}
	scs.ServiceId = cs.ServiceID
	scs.ServiceName = cs.ServiceName
	scs.Version = cs.Version
	scs.AppId = cs.AppID
	scs.Environment = cs.Environment
	scs.Properties = cs.Metadata
	scs.Schemas = cs.Schemas
	scs.Level = cs.Level
	scs.Status = cs.Status
	svcPaths := cs.Paths
	regpaths := []*scregistry.ServicePath{}
	for _, svcPath := range svcPaths {
		var regpath scregistry.ServicePath
		regpath.Path = svcPath.Path
		regpath.Property = svcPath.Property
		regpaths = append(regpaths, &regpath)
	}
	scs.Paths = regpaths
	if cs.Framework != nil {
		scs.Framework = &scregistry.FrameWorkProperty{}
		scs.Framework.Version = cs.Framework.Version
		scs.Framework.Name = cs.Framework.Name
	}
	scs.RegisterBy = cs.RegisterBy
	scs.Alias = cs.Alias
	return scs
}

// ToMicroServiceInstance assign model micro-service instance parameters to registry micro-service instance parameters
func ToMicroServiceInstance(ins *scregistry.MicroServiceInstance) *registry.MicroServiceInstance {
	msi := &registry.MicroServiceInstance{
		InstanceID: ins.InstanceId,
		Metadata:   ins.Properties,
		Status:     ins.Status,
		Version:    ins.Version,
	}
	m, p := registry.GetProtocolMap(ins.Endpoints)
	msi.EndpointsMap = m
	if len(m) != 0 {
		msi.DefaultEndpoint = m[p].GenEndpoint()
		msi.DefaultProtocol = p
	}
	if ins.DataCenterInfo != nil {
		msi.DataCenterInfo = &registry.DataCenterInfo{
			Name:          ins.DataCenterInfo.Name,
			AvailableZone: ins.DataCenterInfo.AvailableZone,
			Region:        ins.DataCenterInfo.Region,
		}
	}
	if msi.Metadata == nil {
		msi.Metadata = make(map[string]string)
	}
	msi.Metadata["version"] = ins.Version
	return msi
}

// ToSCInstance assign registry micro-service instance parameters to model micro-service instance parameters
func ToSCInstance(msi *registry.MicroServiceInstance) *scregistry.MicroServiceInstance {
	si := &scregistry.MicroServiceInstance{}
	eps := registry.GetProtocolList(msi.EndpointsMap)
	si.InstanceId = msi.InstanceID
	si.Endpoints = eps
	si.Properties = msi.Metadata
	si.HostName = msi.HostName
	si.Status = msi.Status
	if msi.DataCenterInfo != nil {
		si.DataCenterInfo = &scregistry.DataCenterInfo{}
		si.DataCenterInfo.Name = msi.DataCenterInfo.Name
		si.DataCenterInfo.AvailableZone = msi.DataCenterInfo.AvailableZone
		si.DataCenterInfo.Region = msi.DataCenterInfo.Region
	}

	return si
}

//ToSCOptions convert registry opstions into sc client options
func ToSCOptions(options registry.Options) sc.Options {
	sco := sc.Options{}
	sco.Timeout = options.Timeout
	sco.TLSConfig = options.TLSConfig
	sco.Endpoints = options.Addrs
	sco.Compressed = options.Compressed
	sco.EnableSSL = options.EnableSSL
	sco.Verbose = options.Verbose
	sco.Version = options.Version
	return sco
}
