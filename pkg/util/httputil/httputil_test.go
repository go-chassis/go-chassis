package httputil_test

import (
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpRequest(t *testing.T) {
	r, err := rest.NewRequest("GET", "http://hello", nil)
	assert.NoError(t, err)
	inv := &invocation.Invocation{Args: r}
	_, err = httputil.HTTPRequest(inv)
	assert.NoError(t, err)
}
