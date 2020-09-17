package rate_test

import (
	"fmt"
	"github.com/go-chassis/go-chassis/v2/resilience/rate"
	"testing"
	"time"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "DEBUG",
	})
	_ = archaius.Init(
		archaius.WithMemorySource())
}

func TestProcessQpsTokenReq(t *testing.T) {
	qps := rate.GetRateLimiters()
	b := qps.TryAccept("serviceName", 100, 10)
	assert.True(t, b)
	b = qps.TryAccept("serviceName.Schema1.op1", 10, 10)
	assert.True(t, b)
}

func TestUpdateRateLimit(t *testing.T) {
	l := rate.GetRateLimiters()
	l.UpdateRateLimit("cse.flowcontrol.Consumer.l.limit.Server.Employee", 200, 1)
	l.UpdateRateLimit("cse.flowcontrol.Provider.l.limit.Server", 100, 1)
}

func TestDeleteRateLimit(t *testing.T) {
	qps := rate.GetRateLimiters()
	qps.DeleteRateLimiter("cse.flowcontrol.Consumer.qps.limit.Server.Employee")
}

func TestLimiters_TryAccept(t *testing.T) {
	after := time.After(5 * time.Second)
	count := 0
	stop := false
	for !stop {
		select {
		case <-after:
			fmt.Println(count)
			stop = true
		default:
			pass := rate.GetRateLimiters().TryAccept("serviceName", 100, 2)
			if pass {
				count++
			}
		}
	}
}
