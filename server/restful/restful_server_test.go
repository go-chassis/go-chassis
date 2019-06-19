package restful

import (
	rf "github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/model"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
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
