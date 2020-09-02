package ratelimiter_test

import (
	"context"
	"github.com/go-chassis/go-chassis/v2/core/governance"
	"github.com/go-chassis/go-chassis/v2/core/marker"
	"net/http"
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/handler"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/middleware/ratelimiter"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{})
}

func TestHandler_Handle(t *testing.T) {
	testName := "api1"
	testMatchPolicy := `
apiPath:
  contains: "api/1"
`
	marker.SaveMatchPolicy(testMatchPolicy, "servicecomb.marker."+testName, testName)

	b := []byte(`
match: api1
rate: 10
burst: 2
`)
	err := governance.ProcessLimiter("servicecomb.rateLimiting.test", string(b))
	assert.NoError(t, err)

	c := handler.Chain{}
	c.AddHandler(&handler.MarkHandler{})
	c.AddHandler(&ratelimiter.Handler{})
	r, _ := http.NewRequest("GET", "/api/1", nil)
	inv := invocation.New(context.TODO())
	inv.Args = r

	c.Next(inv, func(r *invocation.Response) {
		assert.NoError(t, r.Err)
		t.Log(r.Err)
	})
	inv.HandlerIndex = 0
	c.Next(inv, func(r *invocation.Response) {
		assert.Error(t, r.Err)
		t.Log(r.Err)
	})
}
