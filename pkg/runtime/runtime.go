package runtime

//Status
const (
	StatusRunning = "UP"
	StatusDown    = "DOWN"
)

//HostName is the host name of service host
var HostName string

//ServiceID is the service id in registry service
var ServiceID string

//ServiceName represent self name
var ServiceName string

//InstanceID is the instance id in registry service
var InstanceID string

//InstanceStatus is the current status of instance
var InstanceStatus string

// Init runtime information
func Init() error {
	return nil
}
