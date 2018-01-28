// Package client is an interface for any protocol's client
package client

import (
	"golang.org/x/net/context"
)

// 无论client 还是server request response都是协议无关的
// request 和 response 都是编解码无关的

// Client is the interface used to make requests to services.
// It supports Request/Response via Transport
//rcp client,rest client,or you can implement your own
//for rpc,it could be any client over any protocol,such as rpc over tcp ,rpc over http etc
type Client interface {
	Init(...Option) error
	Options() Options
	NewRequest(service, schemaID, operationID string, arg interface{}, reqOpts ...RequestOption) *Request
	Call(ctx context.Context, addr string, req *Request, rsp interface{}, opts ...CallOption) error
	String() string
}

// Request is a struct for request APIs of micro service
type Request struct {
	ID               int
	MicroServiceName string
	Struct           string
	Method           string
	Arg              interface{}
	Metadata         map[string]interface{}
}

// Response is a struct of microservice response APIs
type Response struct {
	ID       int
	Error    string
	Reply    interface{}
	Metadata map[string]interface{}
}
