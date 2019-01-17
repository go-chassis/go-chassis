package client

import "github.com/go-chassis/go-chassis/pkg/scclient/proto"

// MicroServiceRequest is a struct with microservice information
type MicroServiceRequest struct {
	Service *proto.MicroService `json:"service"`
}

// MicroServiceInstanceRequest is struct with microservice instance information
type MicroServiceInstanceRequest struct {
	Instance *proto.MicroServiceInstance `json:"instance"`
}

// MircroServiceDependencyRequest is a struct with dependencies request
type MircroServiceDependencyRequest struct {
	Dependencies []*MicroServiceDependency `json:"dependencies"`
}
