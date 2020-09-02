package servicecomb

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/router"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRouter(t *testing.T) {
	err := archaius.Init(
		archaius.WithMemorySource())
	assert.NoError(t, err)

	json := []byte(`
{"policyType":"RATE","ruleItems":[{"groupName":"s0"},{"groupName":"s1"}]}
`)
	yaml := []byte(`
      - precedence: 1 # 优先级，数字越大优先级越高
        route: #路由规则列表
        - tags:
            version: 1.0
            project: x
          weight: 50 #全重 50%到这里
        - tags:
            version: 2.0
            project: z
          weight: 50 #全重 50%到这里
`)
	archaius.Set(DarkLaunchPrefix+"order", string(json))
	archaius.Set(DarkLaunchPrefixV2+"web", string(yaml))
	routeconf := map[string][]*config.RouteRule{
		"payment": {{
			Precedence: 2,
			Routes: []*config.RouteTag{{
				Tags: map[string]string{
					"version": "0.2",
				},
				Weight: 100,
			}},
		}, {
			Precedence: 1,
			Routes: []*config.RouteTag{{
				Tags: map[string]string{
					"version": "0.1",
				},
				Weight: 100,
			}},
		}},
	}
	r, err := newRouter()
	r.SetRouteRule(routeconf)
	assert.NoError(t, err)
	err = r.Init(router.Options{})
	assert.NoError(t, err)

	t.Run("fetch rule by name ", func(t *testing.T) {
		rr := r.FetchRouteRuleByServiceName("web")
		assert.Equal(t, 1, len(rr))
	})

	t.Run("fire create event ", func(t *testing.T) {
		yaml := []byte(`
      - precedence: 1 # 优先级，数字越大优先级越高
        route: #路由规则列表
        - tags:
            version: 1.0
            project: x
          weight: 50 #全重 50%到这里
        - tags:
            version: 2.0
            project: z
          weight: 50 #全重 50%到这里
`)
		archaius.Set(DarkLaunchPrefixV2+"test", string(yaml))
		time.Sleep(1 * time.Second)
		rr := r.FetchRouteRuleByServiceName("test")
		assert.Equal(t, 1, len(rr))
	})

	t.Run("fire update event ", func(t *testing.T) {
		yaml := []byte(`
      - precedence: 1 # 优先级，数字越大优先级越高
        route: #路由规则列表
        - tags:
            version: 1.0
            project: x
          weight: 100
      - precedence: 2 # 优先级，数字越大优先级越高
        route: #路由规则列表
        - tags:
            version: 2.0
            project: x
          weight: 100
`)
		archaius.Set(DarkLaunchPrefixV2+"web", string(yaml))
		time.Sleep(1 * time.Second)
		rr := r.FetchRouteRuleByServiceName("web")
		assert.Equal(t, 2, len(rr))
	})

	t.Run("fire delete event ", func(t *testing.T) {
		archaius.Delete(DarkLaunchPrefixV2 + "web")
		time.Sleep(1 * time.Second)
		rr := r.FetchRouteRuleByServiceName("web")
		assert.Equal(t, 0, len(rr))
	})
}
