package handler_test

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	_ "github.com/go-chassis/go-chassis/initiator"
	"github.com/go-chassis/go-chassis/pkg/util/tags"
	"github.com/stretchr/testify/assert"
	"net/http"
	"path/filepath"
	"testing"
)

var yamlContent = `---
servicecomb:
  governance:
    Consumer:
      service1:
        policy:
          fault:
            protocols:
              rest:
                abort:
                  httpStatus: 500
                  percent: 100
 `

func TestRestFaultHandler_Names(t *testing.T) {
	restCon := &handler.FaultHandler{}
	conName := restCon.Name()
	assert.Equal(t, "fault-inject", conName)

	microContent := `---
servicecomb:
  service:
    name: Client
    version: 0.1`
	f := prepareConfDir(t)
	prepareTestFile(t, f, "chassis.yaml", "")
	prepareTestFile(t, f, "microservice.yaml", microContent)
	prepareTestFile(t, f, "fault_injection.yaml", yamlContent)

	err := config.Init()
	assert.NoError(t, err)
	archaius.AddFile(filepath.Join(f, "fault_injection.yaml"))
	c := handler.Chain{}
	c.AddHandler(&handler.FaultHandler{})
	c.AddHandler(&normalAfter{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer[handler.FaultInject] = handler.FaultInject

	t.Run("unknown protocol", func(t *testing.T) {
		inv := &invocation.Invocation{
			MicroServiceName: "ShoppingCart",
			Protocol:         "unknown",
		}

		c.Next(inv, func(r *invocation.Response) {
			t.Log(r.Err)
			assert.Error(t, r.Err)
		})

	})
	t.Run("rest protocol to service1", func(t *testing.T) {
		m := map[string]string{
			common.BuildinTagVersion: "0.1",
			common.BuildinTagApp:     "default"}
		inv := &invocation.Invocation{
			MicroServiceName: "service1",
			Protocol:         "rest",
			Reply:            &http.Response{},
			RouteTags: utiltags.Tags{
				KV:    m,
				Label: utiltags.LabelOfTags(m),
			},
		}

		c.Next(inv, func(r *invocation.Response) {
			t.Log(r.Err)
			assert.Error(t, r.Err)
		})
	})
	t.Run("rest protocol to other service", func(t *testing.T) {
		m := map[string]string{
			common.BuildinTagVersion: "0.1",
			common.BuildinTagApp:     "default"}
		inv := &invocation.Invocation{
			MicroServiceName: "service2",
			Protocol:         "rest",
			Reply:            &http.Response{},
			RouteTags: utiltags.Tags{
				KV:    m,
				Label: utiltags.LabelOfTags(m),
			},
		}

		c.Next(inv, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
	})
}
