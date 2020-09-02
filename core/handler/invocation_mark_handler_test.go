package handler_test

import (
	"context"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"github.com/go-chassis/go-chassis/v2/core/governance"
	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestMarkHandler_Handle(t *testing.T) {
	t.Log("testing mark handler")

	c := handler.Chain{}
	c.AddHandler(&handler.MarkHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer[handler.TrafficMarker] = handler.TrafficMarker
	t.Run("test no match policy", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/path/test", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "none", i.GetMark())
	})

	archaius.Init(archaius.WithMemorySource())
	var yamlContent = `
headers:
  cookie:
    regex: "^(.*?;)?(user=jason)(;.*)?$"
  user:
    exact: jason
apiPath:
  contains: "path/test"
  exact: "/test2"
method: GET
`
	archaius.Set(strings.Join([]string{governance.KindMatchPrefix, "match-user-json"}, "."), yamlContent)
	governance.Init()
	t.Run("test request one header", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/path/test", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.NotEqual(t, "match-user-json", i.GetMark())
	})

	t.Run("test request all header", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/path/test", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "match-user-json", i.GetMark())
	})

	t.Run("test request path no match", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/test", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "none", i.GetMark())
	})

	t.Run("test request path exact match", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/test2", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "match-user-json", i.GetMark())
	})

}

func TestMarkHandler_Handle2(t *testing.T) {
	t.Log("testing mark handler")

	c := handler.Chain{}
	c.AddHandler(&handler.MarkHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer[handler.TrafficMarker] = handler.TrafficMarker
	archaius.Init(archaius.WithMemorySource())
	var yamlContent = `
method: GET
`
	archaius.Set(strings.Join([]string{governance.KindMatchPrefix, "match-user-json"}, "."), yamlContent)
	governance.Init()
	t.Run("test match", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/test2", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "match-user-json", i.GetMark())
	})

	t.Run("test no match", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodPost, "http://127.0.0.1:9992/test2", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "none", i.GetMark())
	})
}

func TestMarkHandler_HandleMutilePolicy(t *testing.T) {
	t.Log("testing mark handler")

	c := handler.Chain{}
	c.AddHandler(&handler.MarkHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer[handler.TrafficMarker] = handler.TrafficMarker

	archaius.Init(archaius.WithMemorySource())
	var yamlContent = `
headers:
  cookie:
    regex: "^(.*?;)?(user=jason)(;.*)?$"
  user:
    exact: jason
apiPath:
  contains: "path/test"
  exact: "/test2"
method: GET
`
	var yamlContent2 = `
method: POST 
`
	archaius.Set(strings.Join([]string{governance.KindMatchPrefix, "match-user-json"}, "."), yamlContent)
	archaius.Set(strings.Join([]string{governance.KindMatchPrefix, "match-user-json-2"}, "."), yamlContent2)
	governance.Init()
	t.Run("test request one header", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/path/test", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "none", i.GetMark())
	})

	t.Run("test request all header", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/path/test", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "match-user-json", i.GetMark())
	})

	t.Run("test request match2", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodPost, "http://127.0.0.1:9992/test", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "match-user-json-2", i.GetMark())
	})

	t.Run("test request path exact match", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "http://127.0.0.1:9992/test2", nil)
		c.Next(i, func(r *invocation.Response) {
			assert.NoError(t, r.Err)
		})
		assert.Equal(t, "match-user-json", i.GetMark())
	})

}
