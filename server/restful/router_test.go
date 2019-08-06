package restful

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetRouteSpecs(t *testing.T) {
	_, err := GetRouteSpecs(WrongSchema{})
	assert.Error(t, err)
	_, err = GetRouteSpecs(&WrongSchema{})
	assert.Error(t, err)
	_, err = GetRouteSpecs(&WrongSchema2{})
	assert.Error(t, err)
}

func TestGetRouteGroup(t *testing.T) {
	gn := GetRouteGroup(&GroupSchema{})
	assert.Equal(t, "HelloGroup", gn)
}

func TestGroupRoutePath(t *testing.T) {
	r := &Route{Path: "/SubRoute"}
	GroupRoutePath(r, &GroupSchema{})
	assert.Equal(t, "HelloGroup/SubRoute", r.Path)
}

func TestGetFunctionName(t *testing.T) {
	ws := &WrongSchema{}
	name := getFunctionName(ws.URLPatterns2)
	assert.Equal(t, "URLPatterns2", name)
	name = getFunctionName(TestGetRouteGroup)
	assert.Equal(t, "TestGetRouteGroup", name)
}

func TestBuildRouteHandler(t *testing.T) {
	schma := &FuncNameSchema{}
	ctx := &Context{Ctx: context.TODO()}

	// FuncName
	route := Route{Path: "/FuncName", ResourceFuncName: "Hello"}
	f, err := BuildRouteHandler(&route, schma)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", route.ResourceFuncName)
	f(ctx)
	assert.Equal(t, "World", ctx.Ctx.Value("Hello"))

	// Func
	ctx = &Context{Ctx: context.TODO()}
	route = Route{Path: "/Func", ResourceFunc: schma.Hello}
	f, err = BuildRouteHandler(&route, schma)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", route.ResourceFuncName)
	f(ctx)
	assert.Equal(t, "World", ctx.Ctx.Value("Hello"))

	// Both
	ctx = &Context{Ctx: context.TODO()}
	route = Route{Path: "/BothFuncAndName", ResourceFunc: schma.Hello, ResourceFuncName: "World"}
	f, err = BuildRouteHandler(&route, schma)
	assert.NoError(t, err)
	assert.Equal(t, "Hello", route.ResourceFuncName)
	f(ctx)
	assert.Equal(t, "World", ctx.Ctx.Value("Hello"))
}

type WrongSchema struct {
}

func (r *WrongSchema) Put(b *Context) {
}

//URLPatterns helps to respond for corresponding API calls
func (r *WrongSchema) URLPatterns2() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/", ResourceFuncName: "Put",
			Returns: []*Returns{{Code: 200}}},
	}
}

type WrongSchema2 struct {
}

//URLPatterns helps to respond for corresponding API calls
func (r *WrongSchema2) URLPatterns() {
}

type GroupSchema struct {
}

func (g *GroupSchema) GroupPath() string {
	return "HelloGroup"
}

type FuncNameSchema struct {
}

func (s *FuncNameSchema) URLPatterns() []Route {
	return []Route{
		{Method: http.MethodGet, Path: "/HelloPath", ResourceFunc: s.Hello, ResourceFuncName: "Hello"},
	}
}

func (s *FuncNameSchema) Hello(ctx *Context) {
	ctx.Ctx = context.WithValue(ctx.Ctx, "Hello", "World")
}
