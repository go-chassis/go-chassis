package handler

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/qpslimiter"
)

// ProviderRateLimiterHandler provider rate limiter handler
type ProviderRateLimiterHandler struct{}

// constant for provider qps limiter keys
const (
	ProviderQPSLimit       = "cse.flowcontrol.Provider.qps.limit"
	ProviderLimitKeyGlobal = "cse.flowcontrol.Provider.qps.global.limit"
)

// Handle is to handle provider rateLimiter things
func (rl *ProviderRateLimiterHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	if !archaius.GetBool("cse.flowcontrol.Provider.qps.enabled", true) {
		chain.Next(i, func(r *invocation.InvocationResponse) error {
			return cb(r)
		})
		return
	}

	//provider has limiter only on microservice name.
	if i.SourceMicroService != "" {
		//use chassis Invoker will send SourceMicroService through network
		qpsRate, ok := qpslimiter.GetQPSRate(ProviderQPSLimit + "." + i.SourceMicroService)
		if !ok {
			qpsRate, _ = qpslimiter.GetQPSRate(ProviderLimitKeyGlobal)
			qpslimiter.GetQPSTrafficLimiter().ProcessQPSTokenReq(ProviderLimitKeyGlobal, qpsRate)
		} else {
			qpsRate, _ = qpslimiter.GetQPSRate(ProviderQPSLimit + "." + i.SourceMicroService)
			qpslimiter.GetQPSTrafficLimiter().ProcessQPSTokenReq(ProviderQPSLimit+"."+i.SourceMicroService, qpsRate)
		}

	} else {
		qpsRate, _ := qpslimiter.GetQPSRate(ProviderLimitKeyGlobal)
		qpslimiter.GetQPSTrafficLimiter().ProcessQPSTokenReq(ProviderLimitKeyGlobal, qpsRate)
	}

	//call next chain
	chain.Next(i, func(r *invocation.InvocationResponse) error {
		return cb(r)
	})

}

func newProviderRateLimiterHandler() Handler {
	return &ProviderRateLimiterHandler{}
}

// Name returns the name providerratelimiter
func (rl *ProviderRateLimiterHandler) Name() string {
	return "providerratelimiter"
}
