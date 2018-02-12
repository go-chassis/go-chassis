// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
// Package server is an interface for a micro server
package server

// Server interface for the server it implements innit, register, start, and stop the server..
type Server interface {
	Options() Options
	//Register a schema of microservice,return unique schema id,you can specify schema id and microservice name of this schema
	Register(interface{}, ...RegisterOption) (string, error)
	Start() error
	Stop() error
	String() string
}
