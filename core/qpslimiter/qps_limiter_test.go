package qpslimiter_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/qpslimiter"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"go.uber.org/ratelimit"
)

func initialize() {
	os.Setenv("CHASSIS_HOME", "/tmp/")
	chassisConf := filepath.Join("/tmp/", "conf")
	os.MkdirAll(chassisConf, 0600)
	os.Create(filepath.Join(chassisConf, "chassis.yaml"))
	os.Create(filepath.Join(chassisConf, "microservice.yaml"))
}

func TestProcessQpsTokenReq(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	bucketsize := 100
	r := ratelimit.New(bucketsize)
	qps := qpslimiter.GetQPSTrafficLimiter()
	qps.KeyMap["serviceName"] = r

	log.Println("serviceName key and bucket already exist")
	qps.ProcessQPSTokenReq("serviceName", 100)

	log.Println("rate limit for new operation")
	qps.ProcessQPSTokenReq("serviceName.Schema1.op1", 10)

}

func TestGetQpsRateWithPriority(t *testing.T) {
	initialize()
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	opMeta := qpslimiter.GetConsumerKey(i.SourceMicroService, i.MicroServiceName, i.SchemaID, i.OperationID)

	qps := qpslimiter.GetQPSTrafficLimiter()
	rate, key := qps.GetQPSRateWithPriority(opMeta.OperationQualifiedName, opMeta.SchemaQualifiedName, opMeta.MicroServiceName)
	t.Log("rate is :", rate)
	assert.Equal(t, "cse.flowcontrol.Consumer.qps.limit.service1", key)

	i = &invocation.Invocation{
		MicroServiceName: "service1",
	}
	keys := qpslimiter.GetProviderKey(i.SourceMicroService)
	rate, key = qps.GetQPSRateWithPriority(keys.ServiceOriented, keys.Global)
	assert.Equal(t, "cse.flowcontrol.Provider.qps.global.limit", key)
}

func TestUpdateRateLimit(t *testing.T) {
	qps := qpslimiter.GetQPSTrafficLimiter()
	qps.UpdateRateLimit("cse.flowcontrol.Consumer.qps.limit.Server.Employee", 200)
	qps.UpdateRateLimit("cse.flowcontrol.Provider.qps.limit.Server", 100)
}

func TestDeleteRateLimit(t *testing.T) {
	qps := qpslimiter.GetQPSTrafficLimiter()
	qps.DeleteRateLimiter("cse.flowcontrol.Consumer.qps.limit.Server.Employee")
}
