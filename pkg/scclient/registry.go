package client

import "github.com/go-chassis/go-chassis/pkg/scclient/proto"

const (
	//EventCreate is a constant of type string
	EventCreate string = "CREATE"
	//EventUpdate is a constant of type string
	EventUpdate string = "UPDATE"
	//EventDelete is a constant of type string
	EventDelete string = "DELETE"
	//EventError is a constant of type string
	EventError string = "ERROR"
	//MicorserviceUp is a constant of type string
	MicorserviceUp string = "UP"
	//MicroserviceDown is a constant of type string
	MicroserviceDown string = "DOWN"
	//MSInstanceUP is a constant of type string
	MSInstanceUP string = "UP"
	//MSIinstanceDown is a constant of type string
	MSIinstanceDown string = "DOWN"
	//CheckByHeartbeat is a constant of type string
	CheckByHeartbeat string = "push"
	//DefaultLeaseRenewalInterval is a constant of type int which declares default lease renewal time
	DefaultLeaseRenewalInterval = 30
	//APIPath is a constant of type string
	APIPath = "/registry/v3"
	//MSI_STARTING     string = "STARTING"
	//MSI_OUTOFSERVICE string = "OUTOFSERVICE"
	//CHECK_BY_PLATFORM             string = "pull"
)

// MicroServiceKey is a struct with key information about Microservice
type MicroServiceKey struct {
	Tenant      string `protobuf:"bytes,1,opt,name=tenant" json:"tenant,omitempty"`
	Project     string `protobuf:"bytes,2,opt,name=project" json:"project,omitempty"`
	AppID       string `protobuf:"bytes,3,opt,name=appId" json:"appId,omitempty"`
	ServiceName string `protobuf:"bytes,4,opt,name=serviceName" json:"serviceName,omitempty"`
	Version     string `protobuf:"bytes,5,opt,name=version" json:"version,omitempty"`
	ins         []proto.MicroServiceInstance
}

// ServicePath is a struct with path and property information
type ServicePath struct {
	Path     string            `protobuf:"bytes,1,opt,name=path" json:"path,omitempty"`
	Property map[string]string `protobuf:"bytes,2,rep,name=property" json:"property,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
}

// Framework is a struct which contains name and version of the Framework
type Framework struct {
	Name    string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Version string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
}

// HealthCheck is struct with contains mode, port and interval of sc from which it needs to poll information
type HealthCheck struct {
	Mode     string `protobuf:"bytes,1,opt,name=mode" json:"mode,omitempty"`
	Port     int32  `protobuf:"varint,2,opt,name=port" json:"port,omitempty"`
	Interval int32  `protobuf:"varint,3,opt,name=interval" json:"interval,omitempty"`
	Times    int32  `protobuf:"varint,4,opt,name=times" json:"times,omitempty"`
}

// DataCenterInfo is a struct with contains the zone information of the data center
type DataCenterInfo struct {
	Name          string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Region        string `protobuf:"bytes,2,opt,name=region" json:"region,omitempty"`
	AvailableZone string `protobuf:"bytes,3,opt,name=availableZone" json:"availableZone,omitempty"`
}

// MicroServiceInstanceChangedEvent is a struct to store the Changed event information
type MicroServiceInstanceChangedEvent struct {
	Action   string                      `protobuf:"bytes,2,opt,name=action" json:"action,omitempty"`
	Key      *MicroServiceKey            `protobuf:"bytes,3,opt,name=key" json:"key,omitempty"`
	Instance *proto.MicroServiceInstance `protobuf:"bytes,4,opt,name=instance" json:"instance,omitempty"`
}

// MicroServiceInstanceKey is a struct to key ID's of the microservice
type MicroServiceInstanceKey struct {
	InstanceID string `protobuf:"bytes,1,opt,name=instanceId" json:"instanceId,omitempty"`
	ServiceID  string `protobuf:"bytes,2,opt,name=serviceId" json:"serviceId,omitempty"`
}

// DependencyMicroService is a struct to keep dependency information for the microservice
type DependencyMicroService struct {
	AppID       string `protobuf:"bytes,1,opt,name=appId" json:"appId,omitempty"`
	ServiceName string `protobuf:"bytes,2,opt,name=serviceName" json:"serviceName,omitempty"`
	Version     string `protobuf:"bytes,3,opt,name=version" json:"version,omitempty"`
}

// MicroServiceDependency is a struct to keep the all the dependency information
type MicroServiceDependency struct {
	Consumer  *DependencyMicroService   `protobuf:"bytes,1,opt,name=consumer" json:"consumer,omitempty"`
	Providers []*DependencyMicroService `protobuf:"bytes,2,rep,name=providers" json:"providers,omitempty"`
}

// GetServicesInfoResponse is a struct to keep all the list of services.
type GetServicesInfoResponse struct {
	AllServicesDetail []*ServiceDetail `protobuf:"bytes,2,rep,name=allServicesDetail" json:"allServicesDetail,omitempty"`
}

// ServiceDetail is a struct to store all the relevant information for a microservice
type ServiceDetail struct {
	MicroService         *proto.MicroService           `protobuf:"bytes,1,opt,name=microService" json:"microService,omitempty"`
	Instances            []*proto.MicroServiceInstance `protobuf:"bytes,2,rep,name=instances" json:"instances,omitempty"`
	Providers            []*proto.MicroService         `protobuf:"bytes,5,rep,name=providers" json:"providers,omitempty"`
	Consumers            []*proto.MicroService         `protobuf:"bytes,6,rep,name=consumers" json:"consumers,omitempty"`
	Tags                 map[string]string             `protobuf:"bytes,7,rep,name=tags" json:"tags,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	MicroServiceVersions []string                      `protobuf:"bytes,8,rep,name=microServiceVersions" json:"microServiceVersions,omitempty"`
}
