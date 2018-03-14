package handler

import (
	"errors"
	"strings"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"

	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
	"github.com/cenkalti/backoff"
)

// LBHandler loadbalancer handler struct
type LBHandler struct{}

func (lb *LBHandler) getEndpoint(i *invocation.Invocation, cb invocation.ResponseCallBack) (string, error) {
	var metadata interface{}
	strategy := i.Strategy
	var strategyFun selector.Strategy
	var err error
	if strategy == "" {
		strategyName := config.GetStrategyName(i.SourceMicroService, i.MicroServiceName)
		i.Strategy = strategyName
		strategyFun, err = loadbalance.GetStrategyPlugin(strategyName)
		if err != nil {
			lager.Logger.Errorf(err, selector.LBError{
				Message: "Get strategy [" + strategyName + "] failed."}.Error())
		}
	} else {
		strategyFun, err = loadbalance.GetStrategyPlugin(strategy)
		if err != nil {
			lager.Logger.Errorf(err, selector.LBError{
				Message: "Get strategy [" + strategy + "] failed."}.Error())
		}
	}
	//append filters in config
	filters := config.GetServerListFilters()
	for _, fName := range filters {
		f := loadbalance.Filters[fName]
		if f != nil {
			i.Filters = append(i.Filters, f)
			continue
		}
	}

	if i.Strategy == loadbalance.StrategySessionStickiness {
		metadata = getSessionID(i)
	}

	if i.Strategy == loadbalance.StrategyLatency {
		metadata = i.MicroServiceName + "/" + i.Protocol
	}

	next, err := loadbalance.DefaultSelector.Select(
		i.MicroServiceName, i.Version,
		selector.WithStrategy(strategyFun),
		selector.WithFilter(i.Filters),
		selector.WithAppID(i.AppID),
		selector.WithConsumerID(i.SourceServiceID),
		selector.WithMetadata(metadata))
	if err != nil {
		writeErr(err, cb)
		return "", err
	}

	ins, err := next()
	if err != nil {
		lbErr := selector.LBError{Message: err.Error()}
		writeErr(lbErr, cb)
		return "", lbErr
	}

	var ep string
	if i.Protocol == "" {
		i.Protocol = archaius.GetString("cse.references."+i.MicroServiceName+".transport", ins.DefaultProtocol)
	}
	if i.Protocol == "" {
		for k := range ins.EndpointsMap {
			i.Protocol = k
			break
		}
	}
	ep, ok := ins.EndpointsMap[i.Protocol]
	if !ok {
		errStr := "No available instance support [" + i.Protocol + "] protocol, msName: " + i.MicroServiceName
		lbErr := selector.LBError{Message: errStr}
		lager.Logger.Errorf(nil, lbErr.Error())
		writeErr(lbErr, cb)
		return "", lbErr
	}
	return ep, nil
}

// Handle to handle the load balancing
func (lb *LBHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	if !config.RetryEnabled(i.SourceMicroService, i.MicroServiceName) {
		lb.handleWithNoRetry(chain, i, cb)
	} else {
		lb.handleWithRetry(chain, i, cb)
	}
}

func (lb *LBHandler) handleWithNoRetry(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	ep, err := lb.getEndpoint(i, cb)
	if err != nil {
		writeErr(err, cb)
		return
	}

	i.Endpoint = ep
	chain.Next(i, func(r *invocation.InvocationResponse) error {
		return cb(r)
	})
}

func (lb *LBHandler) handleWithRetry(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	retryOnSame := config.GetRetryOnSame(i.SourceMicroService, i.MicroServiceName)
	retryOnNext := config.GetRetryOnNext(i.SourceMicroService, i.MicroServiceName)
	handlerIndex := chain.HandlerIndex
	for j := 0; j < retryOnNext+1; j++ {

		// exchange and retry on the next server
		ep, err := lb.getEndpoint(i, cb)
		if err != nil {
			writeErr(err, cb)
			return
		}
		// retry on the same server
		lbBackoff := config.GetBackOff(i.SourceMicroService, i.MicroServiceName)
		callTimes := 0
		operation := func() error {
			if callTimes == retryOnSame+1 {
				return backoff.Permanent(errors.New("Retry time expires"))
			}
			callTimes++
			i.Endpoint = ep
			var respErr error
			var callbackErr error
			chain.HandlerIndex = handlerIndex
			chain.Next(i, func(r *invocation.InvocationResponse) error {
				respErr = r.Err
				callbackErr = cb(r)
				return callbackErr
			})
			if respErr != nil {
				return respErr
			}
			if callbackErr != nil {
				return callbackErr
			}
			return nil
		}
		if err = backoff.Retry(operation, lbBackoff); err == nil {
			return
		}
	}
}

// Name returns loadbalance string
func (lb *LBHandler) Name() string {
	return "loadbalance"
}

func newLBHandler() Handler {
	return &LBHandler{}
}

func getSessionID(i *invocation.Invocation) interface{} {
	var metadata interface{}

	switch i.Args.(type) {
	case *rest.Request:
		req := i.Args.(*rest.Request)
		value := req.GetCookie(common.LBSessionID)
		if value != "" {
			metadata = value
		}
	case *fasthttp.Request:
		req := i.Args.(*fasthttp.Request)
		value := string(req.Header.Peek("Cookie"))
		cookieKey := strings.Split(value, "=")
		if value != "" && (cookieKey[0] == common.SessionID || cookieKey[0] == common.LBSessionID) {
			metadata = cookieKey[1]
		}
	}

	return metadata
}

func genKey(s ...string) string {
	return strings.Join(s, ".")
}

func genMsKey(prefix, src, dest, property string) string {
	if src == "" {
		return genKey(prefix, dest, property)
	}
	return genKey(prefix, src, dest, property)
}
