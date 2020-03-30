package match_test

import (
	"context"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/governance"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/match"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestMarkHandler_Handle(t *testing.T) {
	t.Log("testing mark handler")

	c := handler.Chain{}
	c.AddHandler(&match.MarkHandler{})

	config.GlobalDefinition = &model.GlobalCfg{}
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = make(map[string]string)
	config.GlobalDefinition.Cse.Handler.Chain.Consumer[match.TrafficMarker] = match.TrafficMarker
	t.Run("test no match policy", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.Args, _ = rest.NewRequest(http.MethodGet, "cse://127.0.0.1:9992/path/test", nil)
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
		assert.Equal(t, "", i.GetMark())
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
		i.Args, _ = rest.NewRequest(http.MethodGet, "cse://127.0.0.1:9992/path/test", nil)
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
		assert.NotEqual(t, "match-user-json", i.GetMark())
	})

	t.Run("test request all header", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "cse://127.0.0.1:9992/path/test", nil)
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
		assert.Equal(t, "match-user-json", i.GetMark())
	})

	t.Run("test request path no match", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "cse://127.0.0.1:9992/test", nil)
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
		assert.Equal(t, "", i.GetMark())
	})

	t.Run("test request path exact match", func(t *testing.T) {
		i := invocation.New(context.Background())
		i.Metadata = make(map[string]interface{})
		i.SetHeader("user", "jason")
		i.SetHeader("cookie", "asdfojjsdof;user=jason;sfaoabc")
		i.Args, _ = rest.NewRequest(http.MethodGet, "cse://127.0.0.1:9992/test2", nil)
		c.Next(i, func(r *invocation.Response) error {
			assert.NoError(t, r.Err)
			return r.Err
		})
		assert.Equal(t, "match-user-json", i.GetMark())
	})

}
