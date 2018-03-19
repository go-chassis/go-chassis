//Package server is a package for protocol of a micro service
package server

// Server interface for the protocol server, a server should implement init, register, start, and stop
type Server interface {
	//Register a schema of microservice,return unique schema id,you can specify schema id and microservice name of this schema
	Register(interface{}, ...RegisterOption) (string, error)
	Start() error
	Stop() error
	String() string
}
