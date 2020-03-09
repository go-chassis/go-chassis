package handler_test

import (
	"context"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/fileutil"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"

	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}

func TestTransportHandler_HandleRest(t *testing.T) {
	t.Log("testing transport handler with rest protocol")
	microContent := `---
#微服务的私有属性
service_description:
  name: Client
  version: 0.1`

	yamlContent := `---
cse:
  service:
    registry:
      address: http://127.0.0.1:30100
  protocols:
    rest:
      listenAddress: 127.0.0.1:5001
      workerNumber: 1
      advertiseAddress: 127.0.0.1:5001`

	wd, _ := fileutil.GetWorkDir()
	os.Setenv("CHASSIS_HOME", wd)
	defer os.Unsetenv("CHASSIS_HOME")
	chassisConf := filepath.Join(wd, "conf")
	err := os.MkdirAll(chassisConf, 0700)
	assert.NoError(t, err)
	chassisyaml := filepath.Join(chassisConf, "chassis.yaml")
	microserviceyaml := filepath.Join(chassisConf, "microservice.yaml")
	f1, err := os.Create(chassisyaml)
	assert.NoError(t, err)
	f2, err := os.Create(microserviceyaml)
	defer os.RemoveAll(chassisConf)
	assert.NoError(t, err)
	_, err = io.WriteString(f1, yamlContent)
	assert.NoError(t, err)
	_, err = io.WriteString(f2, microContent)

	err = config.Init()
	assert.Nil(t, err)
	t.Logf("%#v", config.GlobalDefinition)

	//dial
	c := &handler.Chain{}
	i := &invocation.Invocation{}
	i.Reply = &helloworld.HelloReply{}
	i.Args, _ = rest.NewRequest(http.MethodGet, "cse://127.0.0.1:9992/path/test", nil)
	i.Reply = rest.NewResponse()
	i.Ctx = context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
		"user": "test",
	})

	i.Endpoint = "127.0.0.1:9992"
	i.Protocol = "rest"

	h := &handler.TransportHandler{}
	c.Handlers = append(c.Handlers, h)

	c.Next(i, func(r *invocation.Response) error {
		t.Log("chain start")
		t.Logf("%#v", r.Result)
		t.Logf("%#v", r.Err)
		assert.Equal(t, nil, r.Result)
		return r.Err
	})

}
