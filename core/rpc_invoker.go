package core

import (
	"context"
	"sync"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
)

// RPCInvoker is rpc invoker
// one invoker for one microservice
// thread safe
type RPCInvoker struct {
	*abstractInvoker
	sync.RWMutex
}

// NewRPCInvoker is gives the object of rpc invoker
func NewRPCInvoker(opt ...Option) *RPCInvoker {
	opts := newOptions(opt...)

	ri := &RPCInvoker{
		abstractInvoker: &abstractInvoker{
			opts: opts,
		},
	}
	return ri
}

// Invoke is for to invoke the functions during API calls
func (ri *RPCInvoker) Invoke(ctx context.Context, microServiceName, schemaID, operationID string, arg interface{}, reply interface{}, options ...InvocationOption) error {
	opts := getOpts(options...)
	if opts.Protocol == "" {
		opts.Protocol = common.ProtocolHighway
	}

	i := invocation.New(ctx)
	i.MicroServiceName = microServiceName
	wrapInvocationWithOpts(i, opts)
	i.SchemaID = schemaID
	i.OperationID = operationID
	i.Args = arg
	i.Reply = reply
	err := ri.invoke(i)
	if err == nil {
		setCookieToCache(*i, getNamespaceFromMetadata(opts.Metadata))
	}
	return err
}
