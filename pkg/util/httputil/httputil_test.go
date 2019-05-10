package httputil_test

import (
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpRequest(t *testing.T) {

	t.Run("convert invocation with ctx to http request,should success",
		func(t *testing.T) {
			r, err := rest.NewRequest("GET", "http://hello", nil)
			assert.NoError(t, err)
			inv := &invocation.Invocation{
				Args: r,
				Ctx: common.NewContext(map[string]string{
					"os":   "mac",
					"user": "peter",
				})}
			r2, err := httputil.HTTPRequest(inv)
			assert.NoError(t, err)
			assert.Equal(t, "mac", r2.Header.Get("os"))
			httputil.SetURI(r, "http://example.com")
			assert.Equal(t, "http://example.com", r.URL.String())

		})
	t.Run("set wrong type to invocation,should fail",
		func(t *testing.T) {
			inv := &invocation.Invocation{
				Ctx: common.NewContext(map[string]string{
					"os":   "mac",
					"user": "peter",
				})}
			_, err := httputil.HTTPRequest(inv)
			assert.Error(t, err)
		})
}
