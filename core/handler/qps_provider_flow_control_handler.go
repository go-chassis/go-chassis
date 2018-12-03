package handler

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/qpslimiter"
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
		chain.Next(i, cb)
		return
	}

	//provider has limiter only on microservice name.
	key := ProviderLimitKeyGlobal
	rate := qpslimiter.DefaultRate
	ok := false
	if i.SourceMicroService != "" {
		//use chassis Invoker will send SourceMicroService through network
		key = ProviderQPSLimit + "." + i.SourceMicroService
		rate, ok = qpslimiter.GetQPSRate(key);
		if !ok {
			key = ProviderLimitKeyGlobal
			rate, _ = qpslimiter.GetQPSRate(ProviderLimitKeyGlobal)
		}

	} else {
		key = ProviderLimitKeyGlobal
		rate, _ = qpslimiter.GetQPSRate(key)
	}

	//qps rate <=0
	if rate <= 0 {
		switch i.Reply.(type) {
		case *http.Response:
			resp := i.Reply.(*http.Response)
			resp.StatusCode = http.StatusServiceUnavailable
		}

		r := &invocation.Response{}
		r.Status = http.StatusServiceUnavailable
		r.Err = errors.New(fmt.Sprintf("%s | %v", key, rate))
		cb(r)
		return
	}
	qpslimiter.GetQPSTrafficLimiter().ProcessQPSTokenReq(key, rate)
	//call next chain
	chain.Next(i, cb)

}

func newProviderRateLimiterHandler() Handler {
	return &ProviderRateLimiterHandler{}
}

// Name returns the name providerratelimiter
func (rl *ProviderRateLimiterHandler) Name() string {
	return "providerratelimiter"
}
