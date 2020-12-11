package restful

import (
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"log"
	"net/http"
	"testing"

	rf "github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/server"
	"github.com/stretchr/testify/assert"
)

var addrHighway = "127.0.0.1:2399"
var addrHighway1 = "127.0.0.1:2330"

func init() {
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.noRefreshSchema", true)
	config.ReadGlobalConfigFromArchaius()
}
func initEnv() {
	defaultChain := make(map[string]string)
	defaultChain["default"] = ""
}

func TestRestStart(t *testing.T) {
	initEnv()
	schema := "schema1"

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(server.Options{
		Address:   addrHighway,
		ChainName: "default",
	})

	_, err = s.Register(&TestSchema{},
		server.WithSchemaID(schema))
	assert.NoError(t, err)

	name := s.String()
	assert.Equal(t, "rest", name)

	err = s.Stop()
	assert.NoError(t, err)
}

func TestRestStartFailure(t *testing.T) {
	t.Log("Testing restful server for start function failure")
	initEnv()
	schema := "schema2"

	//trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.ServiceComb.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.ServiceComb.Handler.Chain.Consumer = defaultChain

	f, err := server.GetServerFunc("rest")
	assert.NoError(t, err)
	s := f(server.Options{
		Address:   addrHighway,
		ChainName: "default",
	})

	_, err = s.Register(TestSchema{},
		server.WithSchemaID(schema))
	assert.Error(t, err)

	err = s.Start()
	assert.NoError(t, err)

	name := s.String()
	assert.Equal(t, "rest", name)

	err = s.Stop()
	assert.NoError(t, err)
}

type TestSchema struct {
}

func (r *TestSchema) Put(b *Context) {
}

func (r *TestSchema) Get(b *Context) {
}

func (r *TestSchema) Delete(b *Context) {
}

func (r *TestSchema) Head(b *Context) {
}
func (r *TestSchema) Patch(b *Context) {
}
func (r *TestSchema) Post(b *Context) {
}

//URLPatterns helps to respond for corresponding API calls
func (r *TestSchema) URLPatterns() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/", ResourceFunc: r.Get,
			Returns: []*Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayhello/{userid}", ResourceFunc: r.Post,
			Returns: []*Returns{{Code: 200}}},

		{Method: http.MethodDelete, Path: "/sayhi", ResourceFunc: r.Delete,
			Returns: []*Returns{{Code: 200}}},

		{Method: http.MethodHead, Path: "/sayjson", ResourceFunc: r.Head,
			Returns: []*Returns{{Code: 200}}},
		{Method: http.MethodPatch, Path: "/sayjson", ResourceFunc: r.Patch,
			Returns: []*Returns{{Code: 200}}},
		{Method: http.MethodPut, Path: "/hi", ResourceFunc: r.Put,
			Returns: []*Returns{{Code: 200}}},
	}
}

func TestFillParam(t *testing.T) {
	var rb = &rf.RouteBuilder{}
	var routeSpec Route
	p := &Parameters{
		"p", "", rf.QueryParameterKind, "", false,
	}
	routeSpec.Parameters = append(routeSpec.Parameters, p)
	p1 := &Parameters{
		"p1", "", rf.BodyParameterKind, "", false,
	}
	routeSpec.Parameters = append(routeSpec.Parameters, p1)
	p2 := &Parameters{
		"p2", "", rf.FormParameterKind, "", false,
	}
	routeSpec.Parameters = append(routeSpec.Parameters, p2)
	p3 := &Parameters{
		"p3", "", rf.HeaderParameterKind, "", false,
	}
	routeSpec.Parameters = append(routeSpec.Parameters, p3)

	rb = fillParam(routeSpec, rb)
	assert.Equal(t, rf.QueryParameterKind, rb.ParameterNamed("p").Kind())
	assert.Equal(t, rf.BodyParameterKind, rb.ParameterNamed("p1").Kind())
	assert.Equal(t, rf.FormParameterKind, rb.ParameterNamed("p2").Kind())
	assert.Equal(t, rf.HeaderParameterKind, rb.ParameterNamed("p3").Kind())

}

var schemaTestProduces = []string{"application/json"}
var schemaTestConsumes = []string{"application/xml"}
var schemaTestRoutes = []Route{
	{
		Method:           http.MethodGet,
		Path:             "none",
		ResourceFuncName: "Handler",
	},
	{
		Method:           http.MethodGet,
		Path:             "with-produces",
		ResourceFuncName: "Handler",
		Produces:         schemaTestProduces,
	},
	{
		Method:           http.MethodGet,
		Path:             "with-consumes",
		ResourceFuncName: "Handler",
		Consumes:         schemaTestConsumes,
	},
	{
		Method:           http.MethodGet,
		Path:             "with-all",
		ResourceFuncName: "Handler",
		Produces:         schemaTestProduces,
		Consumes:         schemaTestConsumes,
	},
}

type SchemaTest struct {
}

func (st SchemaTest) URLPatterns() []Route {
	return schemaTestRoutes
}

func (st SchemaTest) Handler(ctx *Context) {
}

func Test_restfulServer_register2GoRestful(t *testing.T) {
	initEnv()

	rest := &restfulServer{
		container: rf.NewContainer(),
		ws:        new(rf.WebService),
		server:    &http.Server{},
	}

	_, err := rest.Register(&SchemaTest{})
	assert.NoError(t, err)

	routes := rest.ws.Routes()
	assert.Equal(t, 4, len(routes), "there should be %d routes", len(schemaTestRoutes))

	for _, route := range routes {
		switch route.Path {
		case "/none":
			assert.Equal(t, []string{"*/*"}, route.Consumes)
			assert.Equal(t, []string{"*/*"}, route.Produces)
		case "/with-produces":
			assert.Equal(t, schemaTestProduces, route.Produces)
			assert.Equal(t, []string{"*/*"}, route.Consumes)
		case "/with-consumes":
			assert.Equal(t, []string{"*/*"}, route.Produces)
			assert.Equal(t, schemaTestConsumes, route.Consumes)
		case "/with-all":
			assert.Equal(t, schemaTestProduces, route.Produces)
			assert.Equal(t, schemaTestConsumes, route.Consumes)
		default:
			log.Println(route.Path)
		}
	}
}
