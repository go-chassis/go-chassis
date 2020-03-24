package rate_test

import (
	"testing"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/rate"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
	_ = archaius.Init(
		archaius.WithMemorySource())
}

func TestProcessQpsTokenReq(t *testing.T) {
	qps := rate.GetRateLimiters()
	b := qps.TryAccept("serviceName", 100)
	assert.True(t, b)
	b = qps.TryAccept("serviceName.Schema1.op1", 10)
	assert.True(t, b)
}

func TestUpdateRateLimit(t *testing.T) {
	l := rate.GetRateLimiters()
	l.UpdateRateLimit("cse.flowcontrol.Consumer.l.limit.Server.Employee", 200)
	l.UpdateRateLimit("cse.flowcontrol.Provider.l.limit.Server", 100)
}

func TestDeleteRateLimit(t *testing.T) {
	qps := rate.GetRateLimiters()
	qps.DeleteRateLimiter("cse.flowcontrol.Consumer.qps.limit.Server.Employee")
}
