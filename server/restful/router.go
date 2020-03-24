package restful

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-mesh/openlogging"
)

//const for doc
const (
	Path  = "path"
	Query = "query"
)

//Route describe http route path and swagger specifications for API
type Route struct {
	Method           string                 //Method is one of the following: GET,PUT,POST,DELETE. required
	Path             string                 //Path contains a path pattern. required
	ResourceFunc     func(ctx *Context)     //the func this API calls. you must set this field or ResourceFunc, if you set both, ResourceFunc will be used
	ResourceFuncName string                 //the func this API calls. you must set this field or ResourceFunc
	FuncDesc         string                 //tells what this route is all about. Optional.
	Parameters       []*Parameters          //Parameters is a slice of request parameters for a single endpoint Optional.
	Returns          []*Returns             //what kind of response this API returns. Optional.
	Read             interface{}            //Read tells what resource type will be read from the request payload. Optional.
	Consumes         []string               //Consumes specifies that this WebService can consume one or more MIME types.
	Produces         []string               //Produces specifies that this WebService can produce one or more MIME types.
	Metadata         map[string]interface{} //Metadata adds or updates a key=value pair to api
}

//Returns describe response doc
type Returns struct {
	Code    int // http response code
	Message string
	Model   interface{} // response body structure
	Headers map[string]restful.Header
}

//Parameters describe parameters in url path or query params
type Parameters struct {
	Name      string //parameter name
	DataType  string // string, int etc
	ParamType int    //restful.QueryParameterKind or restful.PathParameterKind
	Desc      string
	Required  bool
}

//Router is to define how route the request
type Router interface {
	//URLPatterns returns route
	URLPatterns() []Route
}

//RouteGroup is to define the route group name
type RouteGroup interface {
	//GroupPath if return non-zero-value, it would be appended to route as prefix
	GroupPath() string
}

//GetRouteGroup is to return a router group path
func GetRouteGroup(schema interface{}) string {
	v, ok := schema.(RouteGroup)
	if !ok {
		return ""
	}

	return v.GroupPath()
}

//GetRouteSpecs is to return a rest API specification of a go struct
func GetRouteSpecs(schema interface{}) ([]Route, error) {
	v, ok := schema.(Router)
	if !ok {
		return []Route{}, fmt.Errorf("can not register APIs to server: %s", reflect.TypeOf(schema).String())
	}
	return v.URLPatterns(), nil
}

//WrapHandlerChain wrap business handler with handler chain
func WrapHandlerChain(route *Route, schema interface{}, schemaName string, opts server.Options) (restful.RouteFunction, error) {
	handleFunc, err := BuildRouteHandler(route, schema)
	if err != nil {
		return nil, err
	}
	restHandler := func(req *restful.Request, resp *restful.Response) {
		defer func() {
			if r := recover(); r != nil {
				var stacktrace = GetTrace()
				openlogging.Error("handle request panic.", openlogging.WithTags(openlogging.Tags{
					"path":  route.Path,
					"panic": r,
					"stack": stacktrace,
				}))
				if err := resp.WriteErrorString(http.StatusInternalServerError, "server got a panic, plz check log."); err != nil {
					openlogging.Error("write response failed when handler panic.", openlogging.WithTags(openlogging.Tags{
						"err": err.Error(),
					}))
				}
			}
		}()
		originChain := &handler.Chain{}
		if opts.ChainName != "" {
			originChain, err = handler.GetChain(common.Provider, opts.ChainName)
			if err != nil {
				openlogging.Error("handler chain init err.", openlogging.WithTags(openlogging.Tags{
					"err": err.Error(),
				}))
				resp.AddHeader("Content-Type", "text/plain")
				resp.WriteErrorString(http.StatusInternalServerError, err.Error())
				return
			}
		}

		inv, err := HTTPRequest2Invocation(req, schemaName, route.ResourceFuncName, resp)
		if err != nil {
			openlogging.Error("transfer http request to invocation failed.", openlogging.WithTags(openlogging.Tags{
				"err": err.Error(),
			}))
			return
		}
		bs := NewBaseServer(inv.Ctx)
		bs.Req = req
		bs.Resp = resp
		//create a new chain for each resource handler
		c := &handler.Chain{}
		*c = *originChain
		c.AddHandler(newHandler(handleFunc, bs, opts))
		//give inv.Ctx to user handlers, modules may inject headers in handler chain
		c.Next(inv, func(ir *invocation.Response) error {
			if ir.Err != nil {
				if resp != nil {
					resp.WriteHeader(ir.Status)
				}
				return ir.Err
			}
			return nil
		})

	}

	openlogging.Info("add route path.", openlogging.WithTags(openlogging.Tags{
		"path":      route.Path,
		"method":    route.Method,
		"func_name": route.ResourceFuncName,
	}))
	return restHandler, nil
}

// GroupRoutePath add group route path to route
func GroupRoutePath(route *Route, schema interface{}) {
	groupPath := GetRouteGroup(schema)
	if groupPath != "" {
		route.Path = groupPath + route.Path
	}
}

//BuildRouteHandler build handler func from ResourceFunc or ResourceFuncName
func BuildRouteHandler(route *Route, schema interface{}) (func(ctx *Context), error) {
	if route.ResourceFunc != nil {
		if route.ResourceFuncName == "" {
			route.ResourceFuncName = getFunctionName(route.ResourceFunc)
		}

		return func(ctx *Context) {
			route.ResourceFunc(ctx)
		}, nil
	}

	method, exist := reflect.TypeOf(schema).MethodByName(route.ResourceFuncName)
	if !exist {
		openlogging.GetLogger().Errorf("router func can not find: %s", route.ResourceFuncName)
		return nil, fmt.Errorf("router func can not find: %s", route.ResourceFuncName)
	}

	return func(ctx *Context) {
		method.Func.Call([]reflect.Value{reflect.ValueOf(schema), reflect.ValueOf(ctx)})
	}, nil
}

//getFunctionName get method name from func
func getFunctionName(i interface{}) string {
	metaName := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	metaNameArr := strings.Split(metaName, ".")
	funcName := metaNameArr[len(metaNameArr)-1]

	// replace suffix "-fm" if function is bounded to struct
	reg := regexp.MustCompile("-fm$")
	return reg.ReplaceAllString(funcName, "")
}

//GetTrace get trace
func GetTrace() string {
	var stacktrace string
	for i := 1; ; i++ {
		_, f, l, got := runtime.Caller(i)
		if !got {
			break
		}
		stacktrace += fmt.Sprintf("%s:%d\n", f, l)
	}
	return stacktrace
}
