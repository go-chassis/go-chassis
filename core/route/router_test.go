package router_test

import (
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/registry"
	"github.com/ServiceComb/go-chassis/core/route"
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
  server: #这里就是请求里的host,一般来说推荐直接为sc里的service name，或者k8s的serviceName.namespace.dnsSuffix
    - precedence: 2 #优先权 越大优先级越高
      route:
      - tags:
          version: 1.2 #对接sc如果不填就自动为0.1
          app: HelloWorld #对接sc如果不填就自动为default
        weight: 80 #全重 80%到这里
      - tags:
          version: 2.0
        weight: 20 #全重 20%到这里
      match:
        source: reviews.default.svc.cluster.local
        httpHeaders:
          test: #多个规则的语义是与
            regex: "user"
  ShoppingCart: #这里就是请求里的host,一般来说推荐直接为sc里的service name，或者k8s的serviceName.namespace.dnsSuffix
    - precedence: 2 #优先权 越大优先级越高
      route:
      - tags:
          version: 1.2 #对接sc如果不填就自动为0.1
          app: HelloWorld #对接sc如果不填就自动为default
        weight: 80 #全重 80%到这里
      - tags:
          version: 2.0
        weight: 20 #全重 20%到这里
      match:
        refer: vmall-with-special-header
        source: reviews.default.svc.cluster.local
        sourceTags:
          version: v2 #这里也可以通过 k8s api 或者sc得到
        httpHeaders:
          cookie: #多个规则的语义是与
            regex: "^(.*?;)?(user=jason)(;.*)?$"
    - precedence: 1 #这个语义表示，
      route:
      - tags:
          version: v3
        weight: 100
 `)
var file2 = []byte(`
routeRule:
  catalogue:
    - precedence: 2
      route:
      - tags:
          version: 0.0.1
          app: sockshop
        weight: 100
  orders:
    - precedence: 2
      route:
      - tags:
          version: 0.0.1
          app: sockshop
        weight: 100
  carts:
    - precedence: 2
      route:
      - tags:
          version: 0.0.1
          app: sockshop
        weight: 100
  `)

var rpcfile = []byte(`
routeRule:
  RPCServer: #这里就是请求里的host,一般来说推荐直接为sc里的service name，或者k8s的serviceName.namespace.dnsSuffix
    - precedence: 2 #优先权 越大优先级越高
      route:
      - tags:
          version: v2 #对接sc如果不填就自动为0.1
        weight: 100 #全重 100%到这里
      match:
        headers:
          test: #多个规则的语义是与
            exact: "user"
    - precedence: 1 #优先权 越大优先级越高
      route:
      - tags:
          version: v3
        weight: 100
 `)

/*var file1 = []byte(`
{
  "policyType": "RULE",
  "ruleItems": [
    {
      "groupName": "s1",
      "groupCondition": "version=0.1,0.2,0.3",
      "policyCondition": "test=30"
    },
    {
      "groupName": "s2",
      "groupCondition": "version=0.4",
      "policyCondition": "test=40"
    }
  ]
}`)*/

/*
func initialize() {
	os.Setenv("CHASSIS_HOME", "/tmp/")
	chassisConf := filepath.Join("/tmp/", "conf")
	os.MkdirAll(chassisConf, 0600)
	os.Create(filepath.Join(chassisConf, "chassis.yaml"))
	os.Create(filepath.Join(chassisConf, "microservice.yaml"))
}*/

func TestInit(t *testing.T) {
	r := &config.RouteRule{}
	r.Precedence = 2
	match := config.Match{}
	match.Source = "source.service"
	r.Match = match
	router.Init(map[string][]*config.RouteRule{"server": {r}}, map[string]*config.Match{"template1": &match})
	rr := router.GetRouteRule()
	template := router.GetTemplates()
	assert.Equal(t, "source.service", template["template1"].Source)
	assert.Equal(t, 2, rr["server"][0].Precedence)
	assert.Equal(t, "source.service", rr["server"][0].Match.Source)
}

func TestRPCRoute(t *testing.T) {
	si := &registry.SourceInfo{
		Tags: map[string]string{},
	}
	si.Tags[common.BuildinTagVersion] = "v2"

	c := &config.RouterConfig{}
	if err := yaml.Unmarshal([]byte(rpcfile), c); err != nil {
		t.Error(err)
	}
	router.Init(c.Destinations, c.SourceTemplates)

	header := map[string]string{
		"cookie": "user=jason",
		"X-Age":  "18",
		"test":   "user",
	}

	inv := new(invocation.Invocation)
	inv.MicroServiceName = "RPCServer"
	err := router.Route(header, si, inv)
	assert.Nil(t, err, "")
	assert.Equal(t, "default", inv.AppID)
	assert.Equal(t, "v2", inv.Version)
	assert.Equal(t, "RPCServer", inv.MicroServiceName)
}

func TestRoute(t *testing.T) {
	si := &registry.SourceInfo{
		Tags: map[string]string{},
	}
	si.Name = "vmall"
	si.Tags[common.BuildinTagApp] = "HelloWorld"
	si.Tags[common.BuildinTagVersion] = "v2"

	c := &config.RouterConfig{}
	if err := yaml.Unmarshal([]byte(file), c); err != nil {
		t.Error(err)
	}
	router.Init(c.Destinations, c.SourceTemplates)

	header := map[string]string{
		"cookie": "user=jason",
		"X-Age":  "18",
	}

	inv := new(invocation.Invocation)
	inv.MicroServiceName = "ShoppingCart"

	err := router.Route(header, si, inv)
	assert.Nil(t, err, "")
	assert.Equal(t, "HelloWorld", inv.AppID)
	assert.Equal(t, "1.2", inv.Version)
	assert.Equal(t, "ShoppingCart", inv.MicroServiceName)

	si.Name = "source"
	err = router.Route(header, si, inv)
	assert.Equal(t, "v3", inv.Version)
	assert.Equal(t, "HelloWorld", inv.AppID)

	inv.Version = ""
	inv.MicroServiceName = "server"
	header["test"] = "test"
	si.Name = "reviews.default.svc.cluster.local"
	err = router.Route(header, si, inv)
	assert.Error(t, err, "")

	inv.Version = ""
	inv.AppID = ""
	inv.MicroServiceName = "notexist"
	err = router.Route(header, nil, inv)
	assert.Equal(t, common.DefaultVersion, inv.Version)

}

// TODO
// Comment buggy test cases : Already raised an issue to trace this https://github.com/ServiceComb/go-chassis/issues/5
/*
func TestRoute1(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	initialize()

	err := filepath.InitConfigDir()
	if err != nil {
		t.Error(err)
	}
	err = archaius.Init()
	if err != nil {
		t.Error(err)
	}

	archaius.AddKeyValue("cse.darklaunch.policy.ShoppingCart", string(file1))
	_ = config.InitRouter()
	si := &registry.SourceInfo{
		Tags: map[string]string{},
	}
	si.Name = "vmall"
	si.Tags[common.BuildinTagApp] = "HelloWorld"
	si.Tags[common.BuildinTagVersion] = "v2"

	router.Init(config.GetRouterConfig().Destinations, config.GetRouterConfig().SourceTemplates)

	header := fasthttp.RequestHeader{}
	header.Add("test", "40")
	inv := new(invocation.Invocation)
	inv.MicroServiceName = "ShoppingCart"

	err = router.Route(header, si, inv)
	assert.Nil(t, err, "")
	assert.Equal(t, "HelloWorld", inv.AppID)
	assert.Equal(t, "0.4", inv.Version)
	assert.Equal(t, "ShoppingCart", inv.MicroServiceName)

	header.Set("test", "30")
	err = router.Route(header, si, inv)
	assert.Equal(t, "0.1", inv.Version)
	err = router.Route(header, si, inv)
	assert.Equal(t, "0.2", inv.Version)
	err = router.Route(header, si, inv)
	assert.Equal(t, "0.3", inv.Version)
	err = router.Route(header, si, inv)
	assert.Equal(t, "0.1", inv.Version)
}
*/

func TestRoute2(t *testing.T) {
	c := &config.RouterConfig{}
	if err := yaml.Unmarshal([]byte(file2), c); err != nil {
		t.Error(err)
	}
	router.SetRouteRule(c.Destinations)

	header := map[string]string{}
	inv := new(invocation.Invocation)
	inv.MicroServiceName = "carts"

	err := router.Route(header, nil, inv)
	assert.Nil(t, err, "")
	t.Log(inv.AppID)
	t.Log(inv.Version)
	assert.Equal(t, "sockshop", inv.AppID)
	assert.Equal(t, "0.0.1", inv.Version)

}

func TestMatch(t *testing.T) {
	si := &registry.SourceInfo{
		Tags: map[string]string{},
	}
	si.Name = "service"
	si.Tags[common.BuildinTagApp] = "app"
	si.Tags[common.BuildinTagVersion] = "0.1"
	match := InitMatch("service", "((abc.)def.)ghi", map[string]string{"tag1": "v1"})
	headers := map[string]string{}
	headers["cookie"] = "abc-def-ghi"
	assert.Equal(t, false, router.Match(match, headers, si))
	si.Tags["tag1"] = "v1"
	assert.Equal(t, false, router.Match(match, headers, si))
	headers["age"] = "15"
	assert.Equal(t, true, router.Match(match, headers, si))
	match.HTTPHeaders["noEqual"] = map[string]string{"noEqu": "e"}
	assert.Equal(t, true, router.Match(match, headers, si))
	headers["noEqual"] = "noe"
	assert.Equal(t, true, router.Match(match, headers, si))
	match.HTTPHeaders["noLess"] = map[string]string{"noLess": "100"}
	headers["noLess"] = "120"
	assert.Equal(t, true, router.Match(match, headers, si))
	match.HTTPHeaders["noGreater"] = map[string]string{"noGreater": "100"}
	headers["noGreater"] = "120"
	assert.Equal(t, false, router.Match(match, headers, si))

	si.Name = "error"
	assert.Equal(t, false, router.Match(match, headers, si))

	headers["cookie"] = "7gh"
	si.Name = "service"
	assert.Equal(t, false, router.Match(match, headers, si))
}

func TestFitRate(t *testing.T) {
	tags := InitTags("0.1", "0.2")
	tag, _ := router.FitRate(tags, "service") //0,0
	assert.Equal(t, "0.1", tag.Tags["version"])
	tag, _ = router.FitRate(tags, "service") //100%, 0
	assert.Equal(t, "0.2", tag.Tags["version"])
	tag, _ = router.FitRate(tags, "service") //50%, 50%
	assert.Equal(t, "0.1", tag.Tags["version"])

	count := 100
	for count > 0 {
		go fit()
		count--
	}
}

func fit() {
	tags := InitTags("0.1", "0.2")
	router.FitRate(tags, "service")
}

func TestSortRules(t *testing.T) {
	router.SetRouteRule(InitDests())
	assert.Equal(t, 20, router.SortRules("test")[3].Precedence)
}

func InitDests() map[string][]*config.RouteRule {
	r1 := &config.RouteRule{}
	r2 := &config.RouteRule{}
	r3 := &config.RouteRule{}
	r4 := &config.RouteRule{}
	r5 := &config.RouteRule{}
	r1.Precedence = 20
	r2.Precedence = 30
	r3.Precedence = 50
	r4.Precedence = 40
	r5.Precedence = 10
	r1.Routes = InitTags("0.11", "0.2")
	r2.Routes = InitTags("1.1", "1.2")
	r3.Routes = InitTags("2.1", "2.2")
	match1 := InitMatch("source", "((abc.)def.)ghi", map[string]string{"tag1": "v1"})
	match2 := InitMatch("source", "notmatch", map[string]string{"tag1": "v1"})
	match3 := InitMatch("source1", "((abc.)def.)ghi", map[string]string{"tag1": "v1"})
	r2.Match = match2
	r1.Match = match1
	r3.Match = match3
	rules := []*config.RouteRule{r1, r2, r3, r4, r5}
	dests := map[string][]*config.RouteRule{"test": rules}
	return dests
}

func InitTags(v1 string, v2 string) []*config.RouteTag {
	tag1 := new(config.RouteTag)
	tag2 := new(config.RouteTag)
	tag1.Weight = 50
	tag2.Weight = 50
	tag1.Tags = map[string]string{"version": v1}
	tag2.Tags = map[string]string{"version": v2}
	tags := []*config.RouteTag{tag1, tag2}
	return tags
}

func InitMatch(source string, pat string, tags map[string]string) config.Match {
	match := config.Match{}
	match.Source = source
	match.SourceTags = tags
	regex := map[string]string{"regex": pat}
	greater := map[string]string{"greater": "10"}
	match.HTTPHeaders = map[string]map[string]string{"cookie": regex, "age": greater}
	return match
}
