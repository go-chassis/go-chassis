package core

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
	"golang.org/x/net/context"
	"sync"
)

// RPCInvoker is rpc invoker
//one invoker for one microservice
//thread safe
type RPCInvoker struct {
	sync.RWMutex
	opts Options
}

// NewRPCInvoker is gives the object of rpc invoker
func NewRPCInvoker(opt ...Option) *RPCInvoker {
	opts := newOptions(opt...)

	ri := &RPCInvoker{
		opts: opts,
	}
	//clientPluginName := os.Getenv("rpc_client_plugin")
	//clientF := client.GetClientNewFunc(clientPluginName)
	return ri
}

// newOptions is for updating options
func newOptions(options ...Option) Options {
	opts := DefaultOptions

	for _, o := range options {
		o(&opts)
	}
	if opts.ChainName == "" {
		opts.ChainName = common.DefaultChainName
	}
	return opts
}

// Invoke is for to invoke the functions during API calls
func (ri *RPCInvoker) Invoke(ctx context.Context, microServiceName, schemaID, operationID string, arg interface{}, reply interface{}, options ...InvocationOption) error {
	opts := getOpts(microServiceName, options...)
	if opts.Protocol == "" {
		opts.Protocol = common.ProtocolHighway
	}
	if len(opts.Filters) == 0 {
		opts.Filters = ri.opts.Filters
	}

	md, ok := metadata.FromContext(ctx)
	if ok {
		md[common.HeaderSourceName] = config.SelfServiceName
	} else {
		ctx = metadata.NewContext(context.Background(), map[string]string{
			common.HeaderSourceName: config.SelfServiceName,
		})
	}

	i := invocation.CreateInvocation()
	wrapInvocationWithOpts(i, opts)
	i.MicroServiceName = microServiceName
	i.SchemaID = schemaID
	i.OperationID = operationID
	i.Args = arg
	i.Reply = reply
	i.Ctx = ctx
	c, err := handler.GetChain(common.Consumer, ri.opts.ChainName)
	if err != nil {
		lager.Logger.Errorf(err, "Handler chain init err.")
		return err
	}

	c.Next(i, func(ir *invocation.InvocationResponse) error {
		err = ir.Err
		reply = ir.Result
		return err
	})
	return err
}
