package config_test

import (
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"testing"
)

var file = []byte(`
sourceTemplate:
  vmall-with-special-header:
    source: vmall
    sourceTags:
      version: v2 #这里也可以通过 k8s api 或者sc得到
    httpHeaders:
      cookie: #多个规则的语义是与
        regex: "^(.*?;)?(user=jason)(;.*)?$"
      X-Age: #多个规则的语义是与
        exact: "18"
routeRule:
  ShoppingCart: #这里就是请求里的host,一般来说推荐直接为sc里的service name，或者k8s的serviceName.namespace.dnsSuffix
    - precedence: 2 #优先权 越大优先级越高
      route:
      - tags:
          version: 1.2 #对接sc如果不填就自动为0.1
          app: HelloWorld #对接sc如果不填就自动为default
        weight: 80 #全重 80%到这里
      - tags:
          version: 1.3
          app: HelloWorld
        weight: 20 #全重 20%到这里
      match:
        refer: vmall-with-special-header
        source: vmall
        sourceTags:
            version: v2 #这里也可以通过 k8s api 或者sc得到
        httpHeaders:
            cookie: #多个规则的语义是与
              regex: "^(.*?;)?(user=jason)(;.*)?$"
            X-Age: #多个规则的语义是与
              exact: "18"
    - precedence: 1 #这个语义表示，
      route:
      - tags:
          version: 1.0
        weight: 100
`)

var file1 = []byte(`
{
  "policyType": "RULE",
  "ruleItems": [
    {
      "groupName": "s1",
      "groupCondition": "version=0.3",
      "policyCondition": "test!=30"
    },
    {
      "groupName": "s2",
      "groupCondition": "version=0.4",
      "policyCondition": "t>3"
    }
  ]
}`)

func TestSetConfig(t *testing.T) {
	c := &config.RouterConfig{}
	if err := yaml.Unmarshal([]byte(file), c); err != nil {
		t.Error(err)
	}
	_, ok := c.Destinations["ShoppingCart"]
	assert.Equal(t, true, ok)
	assert.Equal(t, "vmall", c.SourceTemplates["vmall-with-special-header"].Source)
}

func TestInitRouterInit(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	archaius.Init()
	archaius.AddKeyValue("cse.darklaunch.policy.ShoppingCart", string(file1))
	archaius.AddKeyValue(config.TemplateKey, string(file))
	_ = config.InitRouter()
	v, ok := config.GetRouterConfig().Destinations["ShoppingCart"]
	assert.True(t, ok)
	assert.Equal(t, "30", v[0].Match.HTTPHeaders["test"]["noEqu"])
	assert.Equal(t, "vmall", config.GetRouterConfig().SourceTemplates["vmall-with-special-header"].Source)
}
