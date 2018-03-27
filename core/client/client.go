// Package client is an interface for any protocol's client
package client

import "context"

// 无论client 还是server request response都是协议无关的
// request 和 response 都是编解码无关的

// ProtocolClient is the interface to communicate with one kind of ProtocolServer, it is used in transport handler
// rcp protocol client,http protocol client,or you can implement your own
// for example: rpc,it could be any client over any protocol,such as rpc over tcp ,rpc over http etc
type ProtocolClient interface {
	Call(ctx context.Context, addr string, req *Request, rsp interface{}) error
	String() string
}

//NewRequest create common request,you can operate it in you own Call method of a ProtocolClient
func NewRequest(service, schemaID, operationID string, arg interface{}) *Request {
	r := &Request{
		MicroServiceName: service,
		Schema:           schemaID,
		Operation:        operationID,
		Arg:              arg,
	}
	return r
}

// Request is a struct for a protocol request to micro service server
// it includes common attribute for a protocol,
// usually, you should use Arg to wrap your real protocol request
type Request struct {
	ID               int
	MicroServiceName string
	Schema           string
	Operation        string
	Arg              interface{}
	Metadata         map[string]interface{}
}

// Response is a struct of micro service server response
type Response struct {
	ID       int
	Error    string
	Reply    interface{}
	Metadata map[string]interface{}
}
