package client

import "github.com/go-chassis/go-chassis/pkg/scclient/proto"

// ExistenceIDResponse is a structure for microservice with serviceID, schemaID and InstanceID
type ExistenceIDResponse struct {
	ServiceID  string `json:"serviceId,omitempty"`
	SchemaID   string `json:"schemaId,omitempty"`
	InstanceID string `json:"instanceId,omitempty"`
}

// MicroServiceInstancesResponse is a struct with instances information
type MicroServiceInstancesResponse struct {
	Instances []*proto.MicroServiceInstance `json:"instances,omitempty"`
}

// MicroServiceProvideresponse is a struct with provider information
type MicroServiceProvideresponse struct {
	Services []*proto.MicroService `json:"providers,omitempty"`
}

// AppsResponse is a struct with list of app ID's
type AppsResponse struct {
	AppIds []string `json:"appIds,omitempty"`
}
