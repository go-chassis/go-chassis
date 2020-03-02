package qps_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/qps"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"testing"
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
	qps := qps.GetRateLimiters()
	b := qps.TryAccept("serviceName", 100)
	assert.True(t, b)
	b = qps.TryAccept("serviceName.Schema1.op1", 10)
	assert.True(t, b)
}

func TestGetQpsRateWithPriority(t *testing.T) {
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}
	opMeta := qps.GetConsumerKey(i.SourceMicroService, i.MicroServiceName, i.SchemaID, i.OperationID)

	l := qps.GetRateLimiters()
	rate, key := l.GetQPSRateWithPriority(opMeta.OperationQualifiedName, opMeta.SchemaQualifiedName, opMeta.MicroServiceName)
	t.Log("rate is :", rate)
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1", key)

	i = &invocation.Invocation{
		MicroServiceName: "service1",
	}
	keys := qps.GetProviderKey(i.SourceMicroService)
	rate, key = l.GetQPSRateWithPriority(keys.ServiceOriented, keys.Global)
	assert.Equal(t, "cse.flowcontrol.Provider.qps.global.limit", key)
}

func TestUpdateRateLimit(t *testing.T) {
	l := qps.GetRateLimiters()
	l.UpdateRateLimit("cse.flowcontrol.Consumer.l.limit.Server.Employee", 200)
	l.UpdateRateLimit("cse.flowcontrol.Provider.l.limit.Server", 100)
}

func TestDeleteRateLimit(t *testing.T) {
	qps := qps.GetRateLimiters()
	qps.DeleteRateLimiter("cse.flowcontrol.Consumer.qps.limit.Server.Employee")
}
