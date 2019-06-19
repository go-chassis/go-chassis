package handler_test

import (
	_ "github.com/go-chassis/go-chassis/core/router/cse"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/stretchr/testify/assert"
	"testing"
)

type normalAfter struct {
}

func (th *normalAfter) Name() string {
	return "fake"
}

func (th *normalAfter) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	cb(&invocation.Response{})
}
func TestRouterHandler_Handle(t *testing.T) {
	t.Run("decide route with existing route tags, should skip", func(t *testing.T) {
		c := handler.Chain{}
		c.AddHandler(&handler.RouterHandler{})
		c.AddHandler(&normalAfter{})
		i := &invocation.Invocation{
			MicroServiceName: "service1",
			RouteTags:        utiltags.NewDefaultTag("1.0", "appID"),
		}
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
	})

	microContent := `---
service_description:
  name: Client
  version: 0.1`
	var routerContent = `
routeRule:
  service1:  
    - precedence: 1 # 优先级，数字越大优先级越高
      route: #路由规则列表
      - tags:
          version: 1.0
          project: x
        weight: 100 #全重 50%到这里
      match:
        headers:
          os:
            regex: ios
 `

	chassisConf := prepareConfDir(t)
	prepareTestFile(t, chassisConf, "chassis.yaml", "")
	prepareTestFile(t, chassisConf, "microservice.yaml", microContent)
	prepareTestFile(t, chassisConf, "router.yaml", routerContent)

	err := config.Init()
	assert.NoError(t, err)
	err = router.Init()
	assert.NoError(t, err)

	r := &handler.RouterHandler{}
	t.Run("decide route with empty route tags", func(t *testing.T) {
		c := handler.Chain{}
		c.AddHandler(r)
		c.AddHandler(&normalAfter{})
		i := &invocation.Invocation{
			MicroServiceName: "service1",
			Ctx: common.NewContext(map[string]string{
				"os": "ios",
			}),
		}
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			assert.Equal(t, "1.0", i.RouteTags.KV["version"])
			assert.Equal(t, "x", i.RouteTags.KV["project"])
			return r.Err
		})
	})

	assert.Equal(t, "router", r.Name())
}
