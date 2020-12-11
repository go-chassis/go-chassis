package config_test

import (
	"github.com/go-chassis/foundation/stringutil"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"strings"
	"testing"
)

func TestRouterConfig(t *testing.T) {
	b := []byte(`
servicecomb:
  routeRule:
    service1: | 
      # test
      - precedence: 1
        route:  
          - tags:
              version: latest  
            weight: 100  
`)
	c := &config.ServiceComb{}
	err := yaml.Unmarshal(b, c)
	assert.NoError(t, err)
	v, ok := c.Prefix.RouteRule["service1"]
	assert.True(t, ok)

	v = strings.TrimSpace(v)
	t.Log(v)
	type AutoGenerated []int
	b2 := stringutil.Str2bytes(v)
	r := &config.OneServiceRule{}
	err = yaml.Unmarshal(b2, r)
	assert.NoError(t, err)
	t.Log(r)
	assert.Equal(t, 1, r.Len())
	t.Log(r.Value()[0].Precedence)
}
