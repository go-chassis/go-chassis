package handler

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/qpslimiter"
)

// ConsumerRateLimiterHandler consumer rate limiter handler
type ConsumerRateLimiterHandler struct{}

// Handle is handles the consumer rate limiter APIs
func (rl *ConsumerRateLimiterHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	if !archaius.GetBool("cse.flowcontrol.Consumer.qps.enabled", true) {
		chain.Next(i, func(r *invocation.InvocationResponse) error {
			return cb(r)
		})

		return
	}

	//get operation meta info ms.schema, ms.schema.operation, ms
	operationMeta := qpslimiter.InitSchemaOperations(i)
	rl.GetOrCreate(operationMeta)

	chain.Next(i, func(r *invocation.InvocationResponse) error {
		return cb(r)
	})
}

func newConsumerRateLimiterHandler() Handler {
	return &ConsumerRateLimiterHandler{}
}

// Name returns consumerratelimiter string
func (rl *ConsumerRateLimiterHandler) Name() string {
	return "consumerratelimiter"
}

// GetOrCreate is for getting or creating qps limiter meta data
func (rl *ConsumerRateLimiterHandler) GetOrCreate(op *qpslimiter.OperationMeta) {

	qpsRate, key := qpslimiter.GetQPSTrafficLimiter().GetQPSRateWithPriority(op)

	qpslimiter.GetQPSTrafficLimiter().ProcessQPSTokenReq(key, qpsRate)

	return
}
