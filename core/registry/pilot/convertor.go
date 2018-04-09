package pilot

import (
	"fmt"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/metadata"
	"github.com/ServiceComb/go-chassis/core/registry"
	"strings"
)

const (
	// framework name
	Istio = "Istio"
	// microservice level
	DefaultLevel = "BACK"
	// microservice status
	DefaultStatus = "UP"
)

// ToMicroService assign pilot micro-service to go chassis micro-service
func ToMicroService(scs *service) *registry.MicroService {
	cs := &registry.MicroService{}
	cs.ServiceID = scs.ServiceKey
	cs.ServiceName = scs.ServiceKey
	cs.Version = common.DefaultVersion
	cs.AppID = common.DefaultApp
	cs.Level = DefaultLevel
	cs.Status = DefaultStatus
	cs.Framework = &registry.Framework{
		Name:    Istio,
		Version: common.LatestVersion,
	}
	cs.RegisterBy = metadata.PlatformRegistrationComponent
	return cs
}

// ToMicroServiceInstance assign pilot host parameters to registry micro-service instance parameters
func ToMicroServiceInstance(ins *host) *registry.MicroServiceInstance {
	ipPort := fmt.Sprintf("%s:%d", ins.Address, ins.Port)
	msi := &registry.MicroServiceInstance{}
	msi.InstanceID = strings.Replace(ipPort, ":", "_", 1)
	msi.HostName = msi.InstanceID
	msi.EndpointsMap = map[string]string{
		common.ProtocolRest: ipPort,
	}
	msi.DefaultEndpoint = ipPort
	msi.DefaultProtocol = common.ProtocolRest
	if ins.Tags != nil {
		msi.DataCenterInfo = &registry.DataCenterInfo{
			AvailableZone: ins.Tags.AZ,
		}
	}

	return msi
}
