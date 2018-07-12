package kuberegistry

import (
	"strconv"

	"github.com/ServiceComb/go-chassis/core/registry"
	v1 "k8s.io/api/core/v1"
)

func toMicroService(ss *v1.Service) *registry.MicroService {
	return &registry.MicroService{
		ServiceName: ss.Name,
		ServiceID:   string(ss.UID),
		Metadata:    ss.Spec.Selector,
		RegisterBy:  KubeRegistry,
	}
}

func toMicroServiceInstances(ep *v1.Endpoints) []*registry.MicroServiceInstance {
	ins := []*registry.MicroServiceInstance{}
	for _, ss := range ep.Subsets {
		for _, as := range ss.Addresses {
			ins = append(ins, &registry.MicroServiceInstance{
				InstanceID:   string(as.TargetRef.UID),
				ServiceID:    ep.Name + "." + ep.Namespace,
				HostName:     as.Hostname,
				EndpointsMap: torotocolMap(as, ss.Ports),
			})
		}
	}
	return ins
}

func torotocolMap(address v1.EndpointAddress, ports []v1.EndpointPort) map[string]string {
	ret := map[string]string{}
	for _, port := range ports {
		if _, ok := ret[port.Name]; !ok {
			ret[port.Name] = address.IP + ":" + strconv.Itoa(int(port.Port))
			continue
		}
	}
	return ret
}
