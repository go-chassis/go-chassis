package restful

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/emicklei/go-restful"
	rf "github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/stretchr/testify/assert"
)

var addrHighway = "127.0.0.1:2399"
var addrHighway1 = "127.0.0.1:2330"

func initEnv() {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("CHASSIS_HOME"))
	os.Setenv("GO_CHASSIS_SWAGGERFILEPATH", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("GO_CHASSIS_SWAGGERFILEPATH"))
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.Init()
	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition = &model.GlobalCfg{}
}

func TestRestStart(t *testing.T) {
	t.Log("Testing restful server start function")
	initEnv()
	schema := "schema1"

	//trClient := tcp.NewTransport()

	defaultChain := make(map[string]string)
	defaultChain["default"] = ""

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

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

	config.GlobalDefinition.Cse.Handler.Chain.Provider = defaultChain
	config.GlobalDefinition.Cse.Handler.Chain.Consumer = defaultChain

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
		{Method: http.MethodGet, Path: "/", ResourceFuncName: "Get",
			Returns: []*Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayhello/{userid}", ResourceFuncName: "Post",
			Returns: []*Returns{{Code: 200}}},

		{Method: http.MethodDelete, Path: "/sayhi", ResourceFuncName: "Delete",
			Returns: []*Returns{{Code: 200}}},

		{Method: http.MethodHead, Path: "/sayjson", ResourceFuncName: "Head",
			Returns: []*Returns{{Code: 200}}},
		{Method: http.MethodPatch, Path: "/sayjson", ResourceFuncName: "Patch",
			Returns: []*Returns{{Code: 200}}},
		{Method: http.MethodPut, Path: "/hi", ResourceFuncName: "Put",
			Returns: []*Returns{{Code: 200}}},
	}
}

func TestNoRefreshSchemaConfig(t *testing.T) {
	p := os.Getenv("GOPATH")
	os.Setenv("CHASSIS_HOME", filepath.Join(p, "src", "github.com", "go-chassis", "go-chassis", "examples", "discovery", "server"))
	log.Println(os.Getenv("CHASSIS_HOME"))
	config.Init()
	assert.Equal(t, true, config.GlobalDefinition.Cse.NoRefreshSchema)
	config.GlobalDefinition = &model.GlobalCfg{}
}

type Data struct {
	ID         string `json:"priceID"`
	Category   string `json:"type"`
	Value      string `json:"value"`
	CreateTime string `json:"-"`
}

func TestFillParam(t *testing.T) {
	var rb *rf.RouteBuilder = &rf.RouteBuilder{}
	var routeSpec Route
	p := &Parameters{
		"p", "", rf.QueryParameterKind, "",
	}
	routeSpec.Parameters = append(routeSpec.Parameters, p)
	p1 := &Parameters{
		"p1", "", rf.BodyParameterKind, "",
	}
	routeSpec.Parameters = append(routeSpec.Parameters, p1)
	p2 := &Parameters{
		"p2", "", rf.FormParameterKind, "",
	}
	routeSpec.Parameters = append(routeSpec.Parameters, p2)
	p3 := &Parameters{
		"p3", "", rf.HeaderParameterKind, "",
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
		microServiceName: "rest",
		container:        restful.NewContainer(),
		ws:               new(restful.WebService),
		server:           &http.Server{},
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
