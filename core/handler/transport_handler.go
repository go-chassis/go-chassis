package handler

import (
	"strconv"
	"strings"
	"time"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/core/session"
	clientOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
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

	/*	if err != nil {
		r := &invocation.InvocationResponse{
			Err: err,
		}
		lager.Logger.Errorf(err, "GetClient got Error")
		cb(r)
		return
	}*/

	req := c.NewRequest(i.MicroServiceName, i.SchemaID, i.OperationID, i.Args)
	r := &invocation.InvocationResponse{}

	//taking the time elapsed to check for latency aware strategy
	timeBefore := time.Now()
	err = c.Call(i.Ctx, i.Endpoint, req, i.Reply,
		clientOption.WithContentType(i.ContentType),
		clientOption.WithUrlPath(i.URLPathFormat),
		clientOption.WithMethodType(i.MethodType))

	if err != nil {
		if i.Protocol == common.ProtocolRest && i.Strategy == loadbalance.StrategySessionStickiness {
			var reply *rest.Response
			//set cookie in the error response so that the next request will go the same instance
			//if we are not setting the session id in the error response then there is no use of keeping
			//successiveFailedTimes attribute
			if i.Reply != nil && req.Arg != nil {
				reply = i.Reply.(*rest.Response)
				req := req.Arg.(*rest.Request)
				if r != nil {
					session.CheckForSessionID(i, StrategySessionTimeout(i), reply.GetResponse(), req.GetRequest())
				}
			}
			var statusCode int
			//process the error string to retrieve the response code
			actualerrorIs := err.Error()
			errValues := strings.Split(actualerrorIs, ":")
			if len(errValues) == 3 {
				code := strings.Split(errValues[1], " ")
				statusCode, _ = strconv.Atoi(code[1])
			}

			// Only in the following cases of errors the successiveFailedTimes count should be increased
			if statusCode >= 500 || err == fasthttp.ErrConnectionClosed ||
				err == fasthttp.ErrTimeout || err == fasthttp.ErrNoFreeConns {
				successiveFailedTimesIs := StrategySuccessiveFailedTimes(i)
				errCount, ok := loadbalance.SuccessiveFailureCount[i.Endpoint]
				if ok {
					errCount++
					if errCount == successiveFailedTimesIs {
						session.DeletingKeySuccessiveFailure(reply.GetResponse())
						loadbalance.SuccessiveFailureCount[i.Endpoint] = 0
					} else {

						loadbalance.SuccessiveFailureCount[i.Endpoint] = errCount
					}
				} else {
					loadbalance.SuccessiveFailureCount[i.Endpoint] = 1
				}
			}
		} else {
			loadbalance.SuccessiveFailureCount[i.Endpoint] = 0
		}
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
	switch i.Protocol {
	case common.ProtocolRest:
		if i.Strategy == loadbalance.StrategySessionStickiness {
			if i.Reply != nil && req.Arg != nil {
				r := i.Reply.(*rest.Response)
				req := req.Arg.(*rest.Request)
				if r != nil {
					session.CheckForSessionID(i, StrategySessionTimeout(i), r.GetResponse(), req.GetRequest())
				}

			}

		}
	}
	cb(r)
}

func newTransportHandler() Handler {
	return &TransportHandler{}
}
