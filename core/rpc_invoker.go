package core

import (
	"context"
	"sync"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
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

	if ctx == nil {
		ctx = context.WithValue(context.Background(), common.ContextValueKey{}, map[string]string{
			common.HeaderSourceName: config.SelfServiceName,
		})
	} else {
		at, ok := ctx.Value(common.ContextValueKey{}).(map[string]string)
		if ok {
			at[common.HeaderSourceName] = config.SelfServiceName
		} else {
			ctx = context.WithValue(context.Background(), common.ContextValueKey{}, map[string]string{
				common.HeaderSourceName: config.SelfServiceName,
			})
		}
	}

	i := invocation.CreateInvocation()
	wrapInvocationWithOpts(i, opts)
	i.MicroServiceName = microServiceName
	i.SchemaID = schemaID
	i.OperationID = operationID
	i.Args = arg
	i.Reply = reply
	i.Ctx = ctx
	return ri.invoke(i, reply)
}

// abstract invoker is a common invoker for RPC
type abstractInvoker struct {
	opts Options
}

func (ri *abstractInvoker) WithContext(ctx context.Context, key, val string) context.Context {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return metadata.NewContext(ctx, map[string]string{
			key: val,
		})
	}
	md[key] = val
	return ctx
}

func (ri *abstractInvoker) invoke(i *invocation.Invocation, reply interface{}) error {
	if len(i.Filters) == 0 {
		i.Filters = ri.opts.Filters
	}

	// add self service name into remote context, this value used in provider rate limiter
	i.Ctx = ri.WithContext(i.Ctx, common.HeaderSourceName, config.SelfServiceName)

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
