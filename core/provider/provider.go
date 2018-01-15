package provider

import (
	"github.com/ServiceComb/go-chassis/core/invocation"
	"reflect"
)

// Provider is the interface for schema methods
type Provider interface {
	//Register a schema
	Register(schema interface{}) (string, error)
	RegisterName(name string, schema interface{}) error
	//invoke schema function
	Invoke(inv *invocation.Invocation) (interface{}, error)
	GetOperation(schemaID string, operationID string) (Operation, error)
	//if exists in local,for localcall
	Exist(schemaID, operationID string) bool
}

// Operation is the interface for schema parameters
type Operation interface {
	Method() reflect.Method
	Args() []reflect.Type
	Reply() []reflect.Type
}
