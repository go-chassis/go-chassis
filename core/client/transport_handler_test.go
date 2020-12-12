package client_test

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core/client"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	"github.com/go-chassis/go-chassis/v2/pkg/util/fileutil"

	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/examples/schemas/helloworld"

	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.registry.address", "http://127.0.0.1:30100")
	archaius.Set("servicecomb.service.name", "Client")
	runtime.HostName = "localhost"
	config.MicroserviceDefinition = &model.ServiceSpec{}
	archaius.UnmarshalConfig(config.MicroserviceDefinition)
	config.ReadGlobalConfigFromArchaius()
}

func TestTransportHandler_HandleRest(t *testing.T) {
	microContent := `---
#微服务的私有属性
servicecomb:
  service:
	  name: Client
	  version: 0.1`

	yamlContent := `---
servicecomb:
  registry:
      address: http://127.0.0.1:30100
  protocols:
    rest:
      listenAddress: 127.0.0.1:5001
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
	i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/path/test", nil)
	i.Reply = rest.NewResponse()
	i.Ctx = context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
		"user": "test",
	})

	i.Endpoint = "127.0.0.1:9992"
	i.Protocol = "rest"

	h := &client.TransportHandler{}
	c.Handlers = append(c.Handlers, h)

	c.Next(i, func(r *invocation.Response) {
		t.Log("chain start")
		t.Logf("%#v", r.Result)
		t.Logf("%#v", r.Err)
		assert.Equal(t, nil, r.Result)
	})

}
