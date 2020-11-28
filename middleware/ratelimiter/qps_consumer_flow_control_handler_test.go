package ratelimiter_test

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/v2/control"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/examples/schemas/helloworld"
	"github.com/go-chassis/go-chassis/v2/middleware/ratelimiter"
	"github.com/go-chassis/go-chassis/v2/pkg/util/fileutil"
	"github.com/stretchr/testify/assert"

	_ "github.com/go-chassis/go-chassis/v2/control/servicecomb"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
func prepareConfDir(t *testing.T) string {
	wd, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", wd)
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join(wd, "conf")
	logConf := filepath.Join(wd, "log")
	err := os.MkdirAll(chassisConf, 0700)
	assert.NoError(t, err)
	err = os.MkdirAll(logConf, 0700)
	assert.NoError(t, err)
	return chassisConf
}
func prepareTestFile(t *testing.T, confDir, file, content string) {
	fullPath := filepath.Join(confDir, file)
	err := os.Remove(fullPath)
	f, err := os.Create(fullPath)
	assert.NoError(t, err)
	_, err = io.WriteString(f, content)
	assert.NoError(t, err)
}
func TestCBInit(t *testing.T) {
	f := prepareConfDir(t)
	microContent := `---
servicecomb:
  service:
    name: Client
    version: 0.1`

	prepareTestFile(t, f, "chassis.yaml", "")
	prepareTestFile(t, f, "microservice.yaml", microContent)
	err := config.Init()
	assert.NoError(t, err)
	opts := control.Options{
		Infra: config.GlobalDefinition.Panel.Infra,
	}
	err = control.Init(opts)
	assert.NoError(t, err)
}

func TestConsumerRateLimiterDisable(t *testing.T) {
	t.Log("testing consumerratelimiter handler with qps enabled as false")
	gopath := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", gopath+"/src/github.com/go-chassis/go-chassis/v2/examples/discovery/server/")

	config.Init()
	opts := control.Options{
		Infra:   config.GlobalDefinition.Panel.Infra,
		Address: config.GlobalDefinition.Panel.Settings["address"],
	}
	err := control.Init(opts)
	assert.NoError(t, err)
	c := handler.Chain{}
	c.AddHandler(&ratelimiter.ConsumerRateLimiterHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.FlowControl.Consumer.QPS.Enabled = false
	i := &invocation.Invocation{
		SourceMicroService: "service1",
		SchemaID:           "schema1",
		OperationID:        "SayHello",
		Args:               &helloworld.HelloRequest{Name: "peter"},
	}
	c.Next(i, func(r *invocation.Response) {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
	})

}
func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
func TestConsumerRateLimiterHandler_Handle(t *testing.T) {
	t.Log("testing consumerratelimiter handler with qps enabled as true")

	config.Init()

	c := handler.Chain{}
	c.AddHandler(&ratelimiter.ConsumerRateLimiterHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.FlowControl.Consumer.QPS.Enabled = true
	i := &invocation.Invocation{
		MicroServiceName: "service1",
		SchemaID:         "schema1",
		OperationID:      "SayHello",
		Args:             &helloworld.HelloRequest{Name: "peter"},
	}

	c.Next(i, func(r *invocation.Response) {
		assert.NoError(t, r.Err)
		log.Println(r.Result)
	})
}
