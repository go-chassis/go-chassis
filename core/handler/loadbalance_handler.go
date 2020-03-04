package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-chassis/go-chassis/core/registry"
	"io/ioutil"
	"net/http"

	"github.com/cenkalti/backoff"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/loadbalancer"
	"github.com/go-chassis/go-chassis/core/status"
	backoffUtil "github.com/go-chassis/go-chassis/pkg/backoff"
	"github.com/go-chassis/go-chassis/pkg/util"
	"github.com/go-mesh/openlogging"
)

// LBHandler loadbalancer handler struct
type LBHandler struct{}

func (lb *LBHandler) getEndpoint(i *invocation.Invocation, lbConfig control.LoadBalancingConfig) (*registry.Endpoint, error) {
	var strategyFun func() loadbalancer.Strategy
	var err error
	if i.Strategy == "" {
		i.Strategy = lbConfig.Strategy
		strategyFun, err = loadbalancer.GetStrategyPlugin(i.Strategy)
		if err != nil {
			openlogging.GetLogger().Errorf("lb error [%s] because of [%s]", loadbalancer.LBError{
				Message: "Get strategy [" + i.Strategy + "] failed."}.Error(), err.Error())
		}
	} else {
		strategyFun, err = loadbalancer.GetStrategyPlugin(i.Strategy)
		if err != nil {
			openlogging.GetLogger().Errorf("lb error [%s] because of [%s]", loadbalancer.LBError{
				Message: "Get strategy [" + i.Strategy + "] failed."}.Error(), err.Error())
		}
	}
	if len(i.Filters) == 0 {
		i.Filters = lbConfig.Filters
	}

	s, err := loadbalancer.BuildStrategy(i, strategyFun())
	if err != nil {
		return nil, err
	}

	ins, err := s.Pick()
	if err != nil {
		lbErr := loadbalancer.LBError{Message: err.Error()}
		return nil, lbErr
	}

	if i.Protocol == "" {
		i.Protocol = archaius.GetString("cse.references."+i.MicroServiceName+".transport", ins.DefaultProtocol)
	}
	if i.Protocol == "" {
		for k := range ins.EndpointsMap {
			i.Protocol = k
			break
		}
	}
	protocolServer := util.GenProtoEndPoint(i.Protocol, i.Port)
	ep, ok := ins.EndpointsMap[protocolServer]
	if !ok {
		errStr := fmt.Sprintf(
			"No available instance for protocol server [%s] , microservice: %s has %v",
			protocolServer, i.MicroServiceName, ins.EndpointsMap)
		lbErr := loadbalancer.LBError{Message: errStr}
		openlogging.GetLogger().Errorf(lbErr.Error())
		return nil, lbErr
	}
	return ep, nil
}

// Handle to handle the load balancing
func (lb *LBHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	lbConfig := control.DefaultPanel.GetLoadBalancing(*i)
	if !lbConfig.RetryEnabled {
		lb.handleWithNoRetry(chain, i, lbConfig, cb)
	} else {
		lb.handleWithRetry(chain, i, lbConfig, cb)
	}
}

func (lb *LBHandler) handleWithNoRetry(chain *Chain, i *invocation.Invocation, lbConfig control.LoadBalancingConfig, cb invocation.ResponseCallBack) {
	ep, err := lb.getEndpoint(i, lbConfig)
	if err != nil {
		WriteBackErr(err, status.Status(i.Protocol, status.ServiceUnavailable), cb)
		return
	}

	i.Endpoint = ep.Host()
	i.SSLEnable = ep.IsSSLEnable()
	chain.Next(i, cb)
}

func (lb *LBHandler) handleWithRetry(chain *Chain, i *invocation.Invocation, lbConfig control.LoadBalancingConfig, cb invocation.ResponseCallBack) {
	retryOnSame := lbConfig.RetryOnSame
	retryOnNext := lbConfig.RetryOnNext
	handlerIndex := i.HandlerIndex
	var invResp *invocation.Response
	var reqBytes []byte
	if req, ok := i.Args.(*http.Request); ok {
		if req != nil {
			if req.Body != nil {
				reqBytes, _ = ioutil.ReadAll(req.Body)
			}
		}
	}
	// get retry func
	lbBackoff := backoffUtil.GetBackOff(lbConfig.BackOffKind, lbConfig.BackOffMin, lbConfig.BackOffMax)
	callTimes := 0

	ep, err := lb.getEndpoint(i, lbConfig)
	if err != nil {
		// if get endpoint failed, no need to retry
		WriteBackErr(err, status.Status(i.Protocol, status.ServiceUnavailable), cb)
		return
	}
	operation := func() error {
		i.Endpoint = ep.Host()
		i.SSLEnable = ep.IsSSLEnable()
		callTimes++
		var respErr error
		i.HandlerIndex = handlerIndex

		if _, ok := i.Args.(*http.Request); ok {
			i.Args.(*http.Request).Body = ioutil.NopCloser(bytes.NewBuffer(reqBytes))
		}

		chain.Next(i, func(r *invocation.Response) error {
			if r != nil {
				invResp = r
				respErr = invResp.Err
				return invResp.Err
			}
			return nil
		})

		if callTimes >= retryOnSame+1 {
			if retryOnNext <= 0 {
				return backoff.Permanent(errors.New("retry times expires"))
			}
			ep, err = lb.getEndpoint(i, lbConfig)
			if err != nil {
				// if get endpoint failed, no need to retry
				return backoff.Permanent(err)
			}
			callTimes = 0
			retryOnNext--
		}
		return respErr
	}
	if err := backoff.Retry(operation, lbBackoff); err != nil {
		openlogging.GetLogger().Errorf("stop retry , error : %v", err)
	}

	if invResp == nil {
		invResp = &invocation.Response{}
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
