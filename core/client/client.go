// Package client is an interface for any protocol's client
package client

import (
	"context"
	"github.com/ServiceComb/go-chassis/core/invocation"
)

// ProtocolClient is the interface to communicate with one kind of ProtocolServer, it is used in transport handler
// rcp protocol client,http protocol client,or you can implement your own
type ProtocolClient interface {
	Call(ctx context.Context, addr string, inv *invocation.Invocation, rsp interface{}) error
	String() string
}

// Response is a struct of micro service server response
type Response struct {
	ID       int
	Error    string
	Reply    interface{}
	Metadata map[string]interface{}
}
