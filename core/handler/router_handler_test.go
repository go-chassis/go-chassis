package handler_test

import (
	"github.com/go-chassis/go-archaius"
	_ "github.com/go-chassis/go-chassis/core/router/servicecomb"

	"testing"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/router"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/stretchr/testify/assert"
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
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
	})

	var routerContent = `
      - precedence: 1 # 优先级，数字越大优先级越高
        route: 
          - tags:
              version: 1.0
              project: x
            weight: 100
        match:
          headers:
            Os:
              regex: ios
`

	err := archaius.Init(archaius.WithMemorySource())
	assert.NoError(t, err)
	archaius.Set("servicecomb.routeRule.service1", string(routerContent))
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
				"Os": "ios",
			}),
		}
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
			assert.Equal(t, "1.0", i.RouteTags.KV["version"])
			assert.Equal(t, "x", i.RouteTags.KV["project"])
		})
	})

	assert.Equal(t, "router", r.Name())
}
