// Package server is an interface for a micro server
package server

// Server interface for the server it implements innit, register, start, and stop the server..
type Server interface {
	Options() Options
	Init(...Option) error
	//register a schema of microservice,return unique schema id,you can specify schema id and microservice name of this schema
	Register(interface{}, ...RegisterOption) (string, error)
	Start() error
	Stop() error
	String() string
}

// Streamer represents a stream established with a client.
// A stream can be bidirectional which is indicated by the request.
// The last error will be left in Error().
// EOF indicated end of the stream.
//type Streamer interface {
//	Context() context.Context
//	Request() invocation.Invocation
//	Send(interface{}) error
//	Recv(interface{}) error
//	Error() error
//	Close() error
//}
