package core

import (
	"strings"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/session"
)

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

// abstract invoker is a common invoker for RPC
type abstractInvoker struct {
	opts Options
}

func (ri *abstractInvoker) invoke(i *invocation.Invocation) error {
	if len(i.Filters) == 0 {
		i.Filters = ri.opts.Filters
	}

	// add self service name into remote context, this value used in provider rate limiter
	i.Ctx = common.WithContext(i.Ctx, common.HeaderSourceName, runtime.ServiceName)

	c, err := handler.GetChain(common.Consumer, ri.opts.ChainName)
	if err != nil {
		lager.Logger.Errorf("Handler chain init err [%s]", err.Error())
		return err
	}

	c.Next(i, func(ir *invocation.Response) error {
		err = ir.Err
		return err
	})
	return err
}

// setCookieToCache   set go-chassisLB cookie to cache when use SessionStickiness strategy
func setCookieToCache(inv invocation.Invocation, namespace string) {
	if inv.Strategy != loadbalancer.StrategySessionStickiness {
		return
	}
	cookie := session.GetSessionIDFromInv(inv, common.LBSessionID)
	if cookie != "" {
		cookies := strings.Split(cookie, "=")
		if len(cookies) > 1 {
			session.AddSessionStickinessToCache(cookies[1], namespace)
		}
	}
}

// getNamespaceFromMetadata get namespace from opts.Metadata
func getNamespaceFromMetadata(metadata map[string]interface{}) string {
	if namespaceTemp, ok := metadata[common.SessionNameSpaceKey]; ok {
		if v, ok := namespaceTemp.(string); ok {
			return v
		}
	}
	return common.SessionNameSpaceDefaultValue
}
