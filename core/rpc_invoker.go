package core

import (
	"context"
	"sync"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
)

// RPCInvoker is rpc invoker
//one invoker for one microservice
//thread safe
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
	//clientPluginName := os.Getenv("rpc_client_plugin")
	//clientF := client.GetClientNewFunc(clientPluginName)
	return ri
}

// Invoke is for to invoke the functions during API calls
func (ri *RPCInvoker) Invoke(ctx context.Context, microServiceName, schemaID, operationID string, arg interface{}, reply interface{}, options ...InvocationOption) error {
	opts := getOpts(microServiceName, options...)
	if opts.Protocol == "" {
		opts.Protocol = common.ProtocolHighway
	}

	i := invocation.CreateInvocation()
	wrapInvocationWithOpts(i, opts)
	i.MicroServiceName = microServiceName
	i.SchemaID = schemaID
	i.OperationID = operationID
	i.Args = arg
	i.Reply = reply
	i.Ctx = ctx
	return ri.invoke(i)
}
