package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/loadbalancer"

	"github.com/ServiceComb/go-chassis/session"
	"github.com/cenkalti/backoff"
	"github.com/valyala/fasthttp"
)

// LBHandler loadbalancer handler struct
type LBHandler struct{}

func (lb *LBHandler) getEndpoint(i *invocation.Invocation) (string, error) {
	var metadata interface{}
	strategy := i.Strategy
	var strategyFun func() loadbalancer.Strategy
	var err error
	if strategy == "" {
		strategyName := config.GetStrategyName(i.SourceMicroService, i.MicroServiceName)
		i.Strategy = strategyName
		strategyFun, err = loadbalancer.GetStrategyPlugin(strategyName)
		if err != nil {
			lager.Logger.Errorf(err, loadbalancer.LBError{
				Message: "Get strategy [" + strategyName + "] failed."}.Error())
		}
	} else {
		strategyFun, err = loadbalancer.GetStrategyPlugin(strategy)
		if err != nil {
			lager.Logger.Errorf(err, loadbalancer.LBError{
				Message: "Get strategy [" + strategy + "] failed."}.Error())
		}
	}
	//append filters in config
	filters := config.GetServerListFilters()
	for _, fName := range filters {
		f := loadbalancer.Filters[fName]
		if f != nil {
			i.Filters = append(i.Filters, f)
			continue
		}
	}
	var sessionID string
	if i.Strategy == loadbalancer.StrategySessionStickiness {
		sessionID = getSessionID(i)
	}

	if i.Strategy == loadbalancer.StrategyLatency {
		metadata = i.MicroServiceName + "/" + i.Protocol
	}

	if i.Version == "" {
		i.Version = common.LatestVersion
	}

	s, err := loadbalancer.BuildStrategy(i.SourceServiceID,
		i.MicroServiceName, i.AppID, i.Version, i.Protocol, sessionID, i.Filters, strategyFun(), metadata)
	if err != nil {
		return "", err
	}

	ins, err := s.Pick()
	if err != nil {
		lbErr := loadbalancer.LBError{Message: err.Error()}
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
		errStr := fmt.Sprintf("No available instance support ["+i.Protocol+"] protocol,"+
			" msName: "+i.MicroServiceName+" %s %s %s", i.Version, i.AppID, ins.EndpointsMap)
		lbErr := loadbalancer.LBError{Message: errStr}
		lager.Logger.Errorf(nil, lbErr.Error())
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
	ep, err := lb.getEndpoint(i)
	if err != nil {
		writeErr(err, cb)
		return
	}

	i.Endpoint = ep
	chain.Next(i, cb)
}

func (lb *LBHandler) handleWithRetry(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	retryOnSame := config.GetRetryOnSame(i.SourceMicroService, i.MicroServiceName)
	retryOnNext := config.GetRetryOnNext(i.SourceMicroService, i.MicroServiceName)
	handlerIndex := chain.HandlerIndex
	var invResp *invocation.InvocationResponse
	for j := 0; j < retryOnNext+1; j++ {
		// exchange and retry on the next server
		ep, err := lb.getEndpoint(i)
		if err != nil {
			// if get endpoint failed, no need to retry
			writeErr(err, cb)
			return
		}
		// retry on the same server
		lbBackoff := config.GetBackOff(i.SourceMicroService, i.MicroServiceName)
		callTimes := 0
		operation := func() error {
			if callTimes == retryOnSame+1 {
				return backoff.Permanent(errors.New("retry times expires"))
			}
			callTimes++
			i.Endpoint = ep
			var respErr error
			chain.HandlerIndex = handlerIndex
			chain.Next(i, func(r *invocation.InvocationResponse) error {
				if r != nil {
					invResp = r
					respErr = invResp.Err
					return invResp.Err
				}
				return nil
			})
			return respErr
		}
		if err = backoff.Retry(operation, lbBackoff); err == nil {
			break
		}
	}
	if invResp == nil {
		invResp = &invocation.InvocationResponse{}
	}
	cb(invResp)
}

// Name returns loadbalancer string
func (lb *LBHandler) Name() string {
	return "loadbalancer"
}

func newLBHandler() Handler {
	return &LBHandler{}
}

func getSessionID(i *invocation.Invocation) string {
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
	default:
		value := session.GetContextMetadata(i.Ctx, common.LBSessionID)
		if value != "" {
			cookieKey := strings.Split(string(value), "=")
			if len(cookieKey) > 1 {
				metadata = cookieKey[1]
			}
		}
	}

	if metadata == nil {
		metadata = ""
	}

	return metadata.(string)
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
