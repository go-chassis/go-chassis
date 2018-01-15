package handler

import (
	/*"errors"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/route"*/
	"github.com/stretchr/testify/assert"
	/*"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
	"gopkg.in/yaml.v2"
	"log"*/
	"testing"
)

var routeFile = []byte(`
sourceTemplate:
  vmall-with-special-header:
    source: vmall
    sourceTags:
      version: v2
    httpHeaders:
      cookie:
        regex: "^(.*?;)?(user=jason)(;.*)?$"
      X-Age:
        exact: "18"
routeRule:
  server:
    - precedence: 2
      route:
      - tags:
          version: 1.2
          app: HelloWorld
        weight: 80
      - tags:
          version: 2.0
        weight: 20
      match:
        source: reviews.default.svc.cluster.local
        httpHeaders:
          test:
            regex: "user"
      fault:
        rest:
          delay:
            fixedDelay: 2
            percent: 100
          abort:
            httpStatus: 451
            percent: 30
        highway:
          delay:
            fixedDelay: 2
            percent: 100
          abort:
            httpStatus: 451
            percent: 30
  ShoppingCart:
    - precedence: 2
      route:
      - tags:
          version: 1.2
          app: HelloWorld
        weight: 80
      - tags:
          version: 2.0
        weight: 20
      match:
        refer: vmall-with-special-header
        source: reviews.default.svc.cluster.local
        sourceTags:
          version: v2
        httpHeaders:
          cookie:
            regex: "^(.*?;)?(user=jason)(;.*)?$"
      restfault:
        delay:
          fixedDelay: 2
          percent: 100
        abort:
          httpStatus: 451
          percent: 40
    - precedence: 1
      route:
      - tags:
          version: v3
        weight: 100
 `)

/*func TestRestFaultHandler_Handle(t *testing.T) {
	t.Log("testing fault-inject handler")
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	config.Init()
	archaius.Init()

	c := Chain{}
	RegisterHandler("fault-inject", FaultHandle)
	c.AddHandler(FaultHandle())

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Consumer["fault-inject"] = "fault-inject"

	inv := &invocation.Invocation{
		MicroServiceName: "ShoppingCart",
	}

	si := &registry.SourceInfo{
		Tags: map[string]string{},
	}
	si.Name = "vmall"
	si.Tags[common.BuildinTagApp] = "app"
	si.Tags[common.BuildinTagVersion] = "v2"

	cc := &config.RouterConfig{}
	if err := yaml.Unmarshal([]byte(routeFile), cc); err != nil {
		t.Error(err)
	}

	router.Init(cc.Destinations, cc.SourceTemplates)

	header := fasthttp.RequestHeader{}
	header.Add("cookie", "user=jason")
	header.Add("X-Age", "18")
	_ = router.Route(header, si, inv)
	c.Next(inv, func(r *invocation.InvocationResponse) error {
		assert.Error(t, errors.New("injecting abort and delay"), r.Err)
		log.Println(r.Result)
		return r.Err
	})
}*/

func TestRestFaultHandler_Names(t *testing.T) {
	restCon := FaultHandle()
	conName := restCon.Name()
	assert.Equal(t, "fault-inject", conName)

}
