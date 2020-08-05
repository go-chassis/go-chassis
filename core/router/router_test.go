package router_test

import (
	"context"
	"github.com/go-chassis/go-chassis/core/marker"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chassis/go-chassis/core/lager"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/core/router"
	_ "github.com/go-chassis/go-chassis/core/router/servicecomb"
	"github.com/stretchr/testify/assert"
)

func init() {
	lager.Init(&lager.Options{
		LoggerLevel:   "INFO",
		RollingPolicy: "size",
	})
}

func TestBuildRouter(t *testing.T) {
	path := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(path, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))

	config.Init()
	router.BuildRouter("cse")

	err := router.BuildRouter("fake")
	assert.Error(t, err)
	err = router.BuildRouter("cse")
	assert.NoError(t, err)
	assert.NotNil(t, router.DefaultRouter)
}
func TestRPCRoute(t *testing.T) {
	si := &registry.SourceInfo{
		Tags: map[string]string{},
	}
	si.Tags[common.BuildinTagVersion] = "v2"

	d := map[string][]*config.RouteRule{
		"RPCServer": {
			{
				Precedence: 2,
				Match: config.Match{
					Headers: map[string]map[string]string{
						"test": {"exact": "user"},
					},
				},
				Routes: []*config.RouteTag{
					{Weight: 100, Tags: map[string]string{"version": "v2"}},
				},
			},
			{
				Precedence: 1,
				Routes: []*config.RouteTag{
					{Weight: 100, Tags: map[string]string{"version": "v3"}},
				},
			},
		},
	}
	router.BuildRouter("cse")
	router.DefaultRouter.SetRouteRule(d)

	header := map[string]string{
		"cookie": "user=jason",
		"X-Age":  "18",
		"test":   "user",
	}

	inv := new(invocation.Invocation)
	inv.MicroServiceName = "RPCServer"
	err := router.Route(header, si, inv)
	assert.Nil(t, err, "")
	assert.Equal(t, "v2", inv.RouteTags.Version())
	assert.Equal(t, "RPCServer", inv.MicroServiceName)
}

func TestRoute(t *testing.T) {
	si := &registry.SourceInfo{
		Tags: map[string]string{},
	}
	si.Name = "vmall"
	si.Tags[common.BuildinTagApp] = "HelloWorld"
	si.Tags[common.BuildinTagVersion] = "v2"
	d := map[string][]*config.RouteRule{
		"server": {
			{
				Precedence: 2,
				Match: config.Match{
					Headers: map[string]map[string]string{
						"test": {"regex": "user"},
					},
				},
				Routes: []*config.RouteTag{
					{Weight: 80, Tags: map[string]string{"version": "1.2", "app": "HelloWorld"}},
					{Weight: 20, Tags: map[string]string{"version": "2.0", "app": "HelloWorld"}},
				},
			},
		},
		"ShoppingCart": {
			{
				Precedence: 2,
				Match: config.Match{
					Headers: map[string]map[string]string{
						"cookie": {"regex": "^(.*?;)?(user=jason)(;.*)?$"},
					},
				},
				Routes: []*config.RouteTag{
					{Weight: 80, Tags: map[string]string{"version": "1.2", "app": "HelloWorld"}},
					{Weight: 20, Tags: map[string]string{"version": "2.0", "app": "HelloWorld"}},
				},
			}, {
				Precedence: 1,
				Routes: []*config.RouteTag{
					{Weight: 100, Tags: map[string]string{"version": "v3", "app": "HelloWorld"}},
				},
			},
		},
	}
	router.BuildRouter("cse")
	router.DefaultRouter.SetRouteRule(d)

	header := map[string]string{
		"cookie": "user=jason",
		"X-Age":  "18",
	}

	inv := new(invocation.Invocation)
	inv.MicroServiceName = "ShoppingCart"

	err := router.Route(header, si, inv)
	assert.Nil(t, err, "")
	assert.Equal(t, "HelloWorld", inv.RouteTags.AppID())
	assert.Equal(t, "1.2", inv.RouteTags.Version())
	assert.Equal(t, "ShoppingCart", inv.MicroServiceName)

	inv.MicroServiceName = "server"
	header["test"] = "test"
	si.Name = "reviews.default.svc.cluster.local"
	err = router.Route(header, si, inv)
	assert.Nil(t, err, "")

	inv.MicroServiceName = "notexist"
	err = router.Route(header, nil, inv)
	assert.Nil(t, err, "")
}

func TestRoute2(t *testing.T) {

	d := map[string][]*config.RouteRule{
		"catalogue": {
			{
				Precedence: 2,
				Routes: []*config.RouteTag{
					{Weight: 100, Tags: map[string]string{"version": "0.0.1", "app": "sockshop"}},
				},
			},
		},
		"orders": {
			{
				Precedence: 2,
				Routes: []*config.RouteTag{
					{Weight: 100, Tags: map[string]string{"version": "0.0.1", "app": "sockshop"}},
				},
			},
		},
		"carts": {
			{
				Precedence: 2,
				Routes: []*config.RouteTag{
					{Weight: 100, Tags: map[string]string{"version": "0.0.1", "app": "sockshop"}},
				},
			},
		},
	}
	router.BuildRouter("cse")
	router.DefaultRouter.SetRouteRule(d)

	header := map[string]string{}
	inv := new(invocation.Invocation)
	inv.MicroServiceName = "carts"

	err := router.Route(header, nil, inv)
	assert.Nil(t, err, "")
	t.Log(inv.RouteTags.AppID())
	t.Log(inv.RouteTags.Version())
	assert.Equal(t, "sockshop", inv.RouteTags.AppID())
	assert.Equal(t, "0.0.1", inv.RouteTags.Version())
}

func TestMatch(t *testing.T) {
	si := &registry.SourceInfo{
		Tags: map[string]string{},
	}
	si.Name = "service"
	si.Tags[common.BuildinTagApp] = "app"
	si.Tags[common.BuildinTagVersion] = "0.1"
	matchConf := initMatch("service", "((abc.)def.)ghi", map[string]string{"tag1": "v1"})
	headers := map[string]string{}
	headers["cookie"] = "abc-def-ghi"
	assert.Equal(t, false, router.Match(nil, matchConf, headers, si))
	si.Tags["tag1"] = "v1"
	assert.Equal(t, false, router.Match(nil, matchConf, headers, si))
	headers["age"] = "15"
	assert.Equal(t, true, router.Match(nil, matchConf, headers, si))
	matchConf.HTTPHeaders["noEqual"] = map[string]string{"noEqu": "e"}
	assert.Equal(t, true, router.Match(nil, matchConf, headers, si))
	headers["noEqual"] = "noe"
	assert.Equal(t, true, router.Match(nil, matchConf, headers, si))
	matchConf.HTTPHeaders["noLess"] = map[string]string{"noLess": "100"}
	headers["noLess"] = "120"
	assert.Equal(t, true, router.Match(nil, matchConf, headers, si))
	matchConf.HTTPHeaders["noGreater"] = map[string]string{"noGreater": "100"}
	headers["noGreater"] = "120"
	assert.Equal(t, false, router.Match(nil, matchConf, headers, si))

	si.Name = "error"
	assert.Equal(t, false, router.Match(nil, matchConf, headers, si))

	headers["cookie"] = "7gh"
	si.Name = "service"
	assert.Equal(t, false, router.Match(nil, matchConf, headers, si))
}

func TestMatchRefer(t *testing.T) {
	m := config.Match{}
	m.Refer = "testMarker"
	inv := invocation.New(context.TODO())
	b := router.Match(inv, m, nil, nil)
	assert.False(t, b)
	inv.Args, _ = http.NewRequest("GET", "some/api", nil)
	inv.Metadata = make(map[string]interface{})
	testMatchPolicy := `
apiPath:
  contains: "some/api"
method: GET
`
	marker.SaveMatchPolicy(testMatchPolicy, "servicecomb.marker."+m.Refer, m.Refer)
	b = router.Match(inv, m, nil, nil)
	assert.True(t, b)
}

func TestFitRate(t *testing.T) {
	tags := InitTags("0.1", "0.2")
	tag := router.FitRate(tags, "service") //0,0
	assert.Equal(t, "0.1", tag.Tags["version"])
	tag = router.FitRate(tags, "service") //100%, 0
	assert.Equal(t, "0.2", tag.Tags["version"])
	tag = router.FitRate(tags, "service") //50%, 50%
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
	router.BuildRouter("cse")
	router.DefaultRouter.SetRouteRule(InitDests())
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
	match1 := initMatch("source", "((abc.)def.)ghi", map[string]string{"tag1": "v1"})
	match2 := initMatch("source", "notmatch", map[string]string{"tag1": "v1"})
	match3 := initMatch("source1", "((abc.)def.)ghi", map[string]string{"tag1": "v1"})
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

func initMatch(source string, pat string, tags map[string]string) config.Match {
	match := config.Match{}
	match.Source = source
	match.SourceTags = tags
	regex := map[string]string{"regex": pat}
	greater := map[string]string{"greater": "10"}
	match.HTTPHeaders = map[string]map[string]string{"cookie": regex, "age": greater}
	return match
}
