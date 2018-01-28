package handler

import (
	"errors"
	"strings"
	"time"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalance"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/selector"
	"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
	"github.com/cenkalti/backoff"
)

const (
	lbPrefix                                 = "cse.loadbalance"
	propertyStrategyName                     = "strategy.name"
	propertySessionStickinessRuleTimeout     = "SessionStickinessRule.sessionTimeoutInSeconds"
	propertySessionStickinessRuleFailedTimes = "SessionStickinessRule.successiveFailedTimes"
	propertyRetryEnabled                     = "retryEnabled"
	propertyRetryOnNext                      = "retryOnNext"
	propertyRetryOnSame                      = "retryOnSame"
	propertyBackoffKind                      = "backoff.kind"
	propertyBackoffMinMs                     = "backoff.minMs"
	propertyBackoffMaxMs                     = "backoff.maxMs"

	backoffJittered = "jittered"
	backoffConstant = "constant"
	backoffZero     = "zero"
)

// LBHandler loadbalancer handler struct
type LBHandler struct{}

// StrategyName strategy name
func StrategyName(i *invocation.Invocation) string {

	global := archaius.GetString(genKey(lbPrefix, propertyStrategyName), loadbalance.StrategyRoundRobin)
	ms := archaius.GetString(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertyStrategyName), global)
	return ms
}

// StrategySessionTimeout strategy session timeout
func StrategySessionTimeout(i *invocation.Invocation) int {
	global := archaius.GetInt(genKey(lbPrefix, propertySessionStickinessRuleTimeout), 30)
	ms := archaius.GetInt(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertySessionStickinessRuleTimeout), global)

	return ms
}

// StrategySuccessiveFailedTimes strategy successive failed times
func StrategySuccessiveFailedTimes(i *invocation.Invocation) int {
	global := archaius.GetInt(genKey(lbPrefix, propertySessionStickinessRuleFailedTimes), 5)
	ms := archaius.GetInt(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertySessionStickinessRuleFailedTimes), global)

	return ms
}

// retryEnabled retry enabled
func (lb *LBHandler) retryEnabled(i *invocation.Invocation) bool {
	global := archaius.GetBool(genKey(lbPrefix, propertyRetryEnabled), false)
	ms := archaius.GetBool(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertyRetryEnabled), global)
	return ms
}

func (lb *LBHandler) retryOnNext(i *invocation.Invocation) int {
	global := archaius.GetInt(genKey(lbPrefix, propertyRetryOnNext), 0)
	ms := archaius.GetInt(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertyRetryOnNext), global)
	return ms
}

func (lb *LBHandler) retryOnSame(i *invocation.Invocation) int {
	global := archaius.GetInt(genKey(lbPrefix, propertyRetryOnSame), 0)
	ms := archaius.GetInt(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertyRetryOnSame), global)
	return ms
}

func (lb *LBHandler) backoffKind(i *invocation.Invocation) string {
	global := archaius.GetString(genKey(lbPrefix, propertyBackoffKind), backoffZero)
	ms := archaius.GetString(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertyBackoffKind), global)
	return ms
}

func (lb *LBHandler) backoffMinMs(i *invocation.Invocation) int {
	global := archaius.GetInt(genKey(lbPrefix, propertyBackoffMinMs), 0)
	ms := archaius.GetInt(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertyBackoffMinMs), global)
	return ms
}

func (lb *LBHandler) backoffMaxMs(i *invocation.Invocation) int {
	global := archaius.GetInt(genKey(lbPrefix, propertyBackoffMaxMs), 0)
	ms := archaius.GetInt(genMsKey(lbPrefix, i.SourceMicroService, i.MicroServiceName, propertyBackoffMaxMs), global)
	return ms
}

func (lb *LBHandler) getBackOff(i *invocation.Invocation) backoff.BackOff {
	backoffKind := lb.backoffKind(i)
	backMin := lb.backoffMinMs(i)
	backMax := lb.backoffMaxMs(i)
	switch backoffKind {
	case backoffJittered:
		return &backoff.ExponentialBackOff{
			InitialInterval:     time.Duration(backMin) * time.Millisecond,
			RandomizationFactor: backoff.DefaultRandomizationFactor,
			Multiplier:          backoff.DefaultMultiplier,
			MaxInterval:         time.Duration(backMax) * time.Millisecond,
			MaxElapsedTime:      0,
			Clock:               backoff.SystemClock,
		}
	case backoffConstant:
		return backoff.NewConstantBackOff(time.Duration(backMin) * time.Millisecond)
	case backoffZero:
		return &backoff.ZeroBackOff{}
	default:
		lager.Logger.Errorf(nil, "Not support backoff kind: %s, reset to: zero.", backoffKind)
		return &backoff.ZeroBackOff{}
	}
}

func (lb *LBHandler) getEndpoint(i *invocation.Invocation, cb invocation.ResponseCallBack) (string, error) {
	var metadata interface{}
	strategy := i.Strategy
	var strategyFun selector.Strategy
	var err error
	if strategy == "" {
		strategyName := StrategyName(i)
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
	filters := archaius.GetServerListFilters()
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
	if !lb.retryEnabled(i) {
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
	retryOnSame := lb.retryOnSame(i)
	retryOnNext := lb.retryOnNext(i)
	handlerIndex := chain.HandlerIndex
	for j := 0; j < retryOnNext+1; j++ {

		// exchange and retry on the next server
		ep, err := lb.getEndpoint(i, cb)
		if err != nil {
			writeErr(err, cb)
			return
		}
		// retry on the same server
		lbBackoff := lb.getBackOff(i)
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
		value := req.GetHeader("Cookie")
		cookieKey := strings.Split(value, "=")
		if value != "" && (cookieKey[0] == common.SessionID || cookieKey[0] == common.LBSessionID) {
			metadata = cookieKey[1]
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
