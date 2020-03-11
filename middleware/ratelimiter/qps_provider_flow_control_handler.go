package ratelimiter

import (
	"github.com/go-chassis/go-chassis/control"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/rate"
)

// ProviderRateLimiterHandler provider rate limiter handler
type ProviderRateLimiterHandler struct{}

// Handle is to handle provider rateLimiter things
func (rl *ProviderRateLimiterHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	rlc := control.DefaultPanel.GetRateLimiting(*i, common.Provider)
	if !rlc.Enabled {
		chain.Next(i, cb)

		return
	}
	//qps rate <=0
	if rlc.Rate <= 0 {
		r := newErrResponse(i, rlc)
		cb(r)
		return
	}
	if rate.GetRateLimiters().TryAccept(rlc.Key, rlc.Rate) {
		chain.Next(i, cb)
	} else {
		r := newErrResponse(i, rlc)
		cb(r)
	}
	return
}

func newProviderRateLimiterHandler() handler.Handler {
	return &ProviderRateLimiterHandler{}
}

// Name returns the name providerratelimiter
func (rl *ProviderRateLimiterHandler) Name() string {
	return "providerratelimiter"
}
