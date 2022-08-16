package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/loadbalancer"
	"github.com/go-chassis/go-chassis/v2/session"
	"github.com/go-chassis/openlog"
)

// TransportHandler transport handler
type TransportHandler struct{}

// Name returns transport string
func (th *TransportHandler) Name() string {
	return "transport"
}
func errNotNil(err error, cb invocation.ResponseCallBack) {
	r := &invocation.Response{
		Err: err,
	}
	openlog.Error("GetClient got Error: " + err.Error())
	cb(r)
}

// Handle is to handle transport related things
func (th *TransportHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {

	c, err := GetClient(i)
	if err != nil {
		errNotNil(err, cb)
		return
	}

	r := &invocation.Response{}

	//taking the time elapsed to check for latency aware strategy
	timeBefore := time.Now()
	err = c.Call(i.Ctx, i.Endpoint, i, i.Reply)
	if err != nil {
		r.Err = err
		if !errors.Is(err, ErrCanceled) {
			openlog.Error(fmt.Sprintf("call err [%s]", err.Error()))
		}
		if i.Strategy == loadbalancer.StrategySessionStickiness {
			ProcessSpecialProtocol(i)
			ProcessSuccessiveFailure(i)
		}
		r.Status, _ = c.Status(i.Reply)
		cb(r)
		return
	}
	r.Status, err = c.Status(i.Reply)
	if err != nil {
		r.Err = err
		cb(r)
		return
	}
	if i.Strategy == loadbalancer.StrategyLatency {
		timeAfter := time.Since(timeBefore)
		loadbalancer.SetLatency(timeAfter, i.Endpoint, i.MicroServiceName, i.RouteTags, i.Protocol)
	}

	if i.Strategy == loadbalancer.StrategySessionStickiness {
		ProcessSpecialProtocol(i)
	}

	r.Result = i.Reply
	cb(r)
}

// ProcessSpecialProtocol handles special logic for protocol
func ProcessSpecialProtocol(inv *invocation.Invocation) {
	switch inv.Protocol {
	case common.ProtocolRest:
		var reply *http.Response
		if inv.Reply != nil && inv.Args != nil {
			reply = inv.Reply.(*http.Response)
			req := inv.Args.(*http.Request)
			session.SaveSessionIDFromHTTP(inv.Endpoint, config.GetSessionTimeout(inv.SourceMicroService, inv.MicroServiceName), reply, req)
		}
	case common.ProtocolHighway:
		inv.Ctx = session.SaveSessionIDFromContext(inv.Ctx, inv.Endpoint, config.GetSessionTimeout(inv.SourceMicroService, inv.MicroServiceName))
	}
}

// ProcessSuccessiveFailure handles special logic for protocol
func ProcessSuccessiveFailure(i *invocation.Invocation) {
	var cookie string
	var reply *http.Response

	switch i.Protocol {
	case common.ProtocolRest:
		if i.Reply != nil && i.Args != nil {
			reply = i.Reply.(*http.Response)
		}
		cookie = session.GetSessionCookie(context.TODO(), reply)
		if cookie != "" {
			loadbalancer.IncreaseSuccessiveFailureCount(cookie)
			errCount := loadbalancer.GetSuccessiveFailureCount(cookie)
			if errCount == config.StrategySuccessiveFailedTimes(i.SourceServiceID, i.MicroServiceName) {
				session.DeletingKeySuccessiveFailure(reply)
				loadbalancer.DeleteSuccessiveFailureCount(cookie)
			}
		}
	default:
		cookie = session.GetSessionCookie(i.Ctx, nil)
		if cookie != "" {
			loadbalancer.IncreaseSuccessiveFailureCount(cookie)
			errCount := loadbalancer.GetSuccessiveFailureCount(cookie)
			if errCount == config.StrategySuccessiveFailedTimes(i.SourceServiceID, i.MicroServiceName) {
				session.DeletingKeySuccessiveFailure(nil)
				loadbalancer.DeleteSuccessiveFailureCount(cookie)
			}
		}
	}
}

func newTransportHandler() handler.Handler {
	return &TransportHandler{}
}
func init() {
	err := handler.RegisterHandler(handler.Transport, newTransportHandler)
	if err != nil {
		openlog.Fatal("can not init chassis" + err.Error())
	}
}
