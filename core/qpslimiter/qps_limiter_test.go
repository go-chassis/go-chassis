package qpslimiter_test

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/qpslimiter"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
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
	r := qpslimiter.New(bucketsize)
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
	archaius.Init()
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	opMeta := qpslimiter.InitSchemaOperations(i)

	qps := qpslimiter.GetQPSTrafficLimiter()
	rate, key := qps.GetQPSRateWithPriority(opMeta)
	log.Println("rate is :", rate)
	assert.Equal(t, key, "cse.flowcontrol.Consumer.qps.limit.service1")
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
