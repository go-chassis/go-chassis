package handler

import (
	"crypto/tls"
	"time"

	"github.com/go-chassis/go-chassis/core/client"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	chassisTLS "github.com/go-chassis/go-chassis/core/tls"
	"github.com/go-chassis/go-chassis/session"
	"github.com/go-mesh/openlogging"
	"net/http"
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
	openlogging.Error("GetClient got Error: " + err.Error())
	cb(r)
	return
}

// Handle is to handle transport related things
func (th *TransportHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {

	c, err := client.GetClient(i.Protocol, i.MicroServiceName, i.GenEndpoint())
	if err != nil {
		errNotNil(err, cb)
		return
	}
	if err := setClientTLSConfig(i, c); err != nil {
		errNotNil(err, cb)
		return
	}

	r := &invocation.Response{}

	//taking the time elapsed to check for latency aware strategy
	timeBefore := time.Now()
	err = c.Call(i.Ctx, i.Endpoint, i, i.Reply)
	if resp, ok := i.Reply.(*http.Response); ok {
		r.Status = resp.StatusCode
	}
	if err != nil {
		r.Err = err
		if err != client.ErrCanceled {
			openlogging.GetLogger().Errorf("Call got Error, err [%s]", err.Error())
		}
		if i.Strategy == loadbalancer.StrategySessionStickiness {
			ProcessSpecialProtocol(i)
			ProcessSuccessiveFailure(i)
		}
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

func setClientTLSConfig(i *invocation.Invocation, c client.ProtocolClient) error {
	var tlsConfig *tls.Config
	tlsConfig, err := client.ReadTLSConfig(i.MicroServiceName, i.Protocol)
	//it will set tls config when provider's endpoint has sslEnable=true suffix or
	// consumer had set provider tls config
	if i.SSLEnable && err != nil {
		return err
	} else if !i.SSLEnable && !chassisTLS.IsSSLConfigNotExist(err) {
		return err
	}
	opts := c.GetOptions()
	opts.TLSConfig = tlsConfig
	c.ReloadConfigs(opts)
	return nil
}

//ProcessSpecialProtocol handles special logic for protocol
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

//ProcessSuccessiveFailure handles special logic for protocol
func ProcessSuccessiveFailure(i *invocation.Invocation) {
	var cookie string
	var reply *http.Response

	switch i.Protocol {
	case common.ProtocolRest:
		if i.Reply != nil && i.Args != nil {
			reply = i.Reply.(*http.Response)
		}
		cookie = session.GetSessionCookie(nil, reply)
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

func newTransportHandler() Handler {
	return &TransportHandler{}
}
