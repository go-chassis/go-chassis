// Package client is an interface for any protocol's client
package client

import (
	"context"
	"errors"

	"github.com/go-chassis/go-chassis/v2/core/invocation"
)

// ErrCanceled means Request is canceled by context management
var ErrCanceled = errors.New("request cancelled")

// TransportFailure is caused by client call failure
// for example:  resp, err = client.Do(req)
// if err is not nil then should wrap original error with TransportFailure
type TransportFailure struct {
	Message string
}

// Error return error message
func (e TransportFailure) Error() string {
	return e.Message
}

// ProtocolClient is a interface to communicate with one kind of ProtocolServer, it is used in transport handler.
// this handler orchestrate client implementation.
// gRPC protocol client, http protocol client, or you can implement your own.
type ProtocolClient interface {
	// TODO use invocation.Response as rsp
	// Call is the key function you must implement
	Call(ctx context.Context, addr string, inv *invocation.Invocation, rsp interface{}) error
	// if your protocol has response status(such as http return 200, 500 status code),
	// you need to return it according to response
	Status(rsp interface{}) (status int, err error)
	String() string
	Close() error
	// if you want to reload client settings on-fly, such as timeout, TLS config,
	// you need to implement it
	ReloadConfigs(Options)
	GetOptions() Options
}
