package handler

import (
	"time"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/session"
	clientOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	microClient "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
)

// TransportHandler transport handler
type TransportHandler struct{}

// Name returns transport string
func (th *TransportHandler) Name() string {
	return "transport"
}
func errNotNill(err error, cb invocation.ResponseCallBack) {
	r := &invocation.InvocationResponse{
		Err: err,
	}
	lager.Logger.Error("GetClient got Error", err)
	cb(r)
	return
}

// Handle is to handle transport related things
func (th *TransportHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	c, err := client.GetClient(i.Protocol, i.MicroServiceName)
	if err != nil {
		errNotNill(err, cb)
	}

	req := c.NewRequest(i.MicroServiceName, i.SchemaID, i.OperationID, i.Args)
	r := &invocation.InvocationResponse{}

	//taking the time elapsed to check for latency aware strategy
	timeBefore := time.Now()
	err = c.Call(i.Ctx, i.Endpoint, req, i.Reply,
		clientOption.WithContentType(i.ContentType),
		clientOption.WithUrlPath(i.URLPathFormat),
		clientOption.WithMethodType(i.MethodType))

	if err != nil {
		r.Err = err
		lager.Logger.Errorf(err, "Call got Error")
		cb(r)
		return
	}

	if i.Strategy == loadbalance.StrategyLatency {
		timeAfter := time.Since(timeBefore)
		loadbalance.SetLatency(timeAfter, i.Endpoint, req.MicroServiceName+"/"+i.Protocol)
	}

	r.Result = i.Reply
	ProcessSpecialProtocol(i, req)

	cb(r)
}

//ProcessSpecialProtocol handles special logic for protocol
func ProcessSpecialProtocol(inv *invocation.Invocation, req *microClient.Request) {
	switch inv.Protocol {
	case common.ProtocolRest:
		if inv.Strategy == loadbalance.StrategySessionStickiness {
			var reply *rest.Response
			//set cookie in the error response so that the next request will go the same instance
			//if we are not setting the session id in the error response then there is no use of keeping
			//successiveFailedTimes attribute
			if inv.Reply != nil && inv.Args != nil {
				reply = inv.Reply.(*rest.Response)
				req := req.Arg.(*rest.Request)
				session.CheckForSessionID(inv, StrategySessionTimeout(inv), reply.GetResponse(), req.GetRequest())
			}
			errCount := loadbalance.GetSuccessiveFailureCount(inv.Endpoint)
			loadbalance.IncreaseSuccessiveFailureCount(inv.Endpoint)
			if errCount == StrategySuccessiveFailedTimes(inv) {
				session.DeletingKeySuccessiveFailure(reply.GetResponse())
				loadbalance.ResetSuccessiveFailureCount(inv.Endpoint)
			}
		}
	}
}
func newTransportHandler() Handler {
	return &TransportHandler{}
}
