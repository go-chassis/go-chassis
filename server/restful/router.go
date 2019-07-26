package restful

import (
	"fmt"
	"net/http"
	"reflect"

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
	Method           string        //Method is one of the following: GET,PUT,POST,DELETE. required
	Path             string        //Path contains a path pattern. required
	ResourceFuncName string        //the func this API calls. required
	FuncDesc         string        //tells what this route is all about. Optional.
	Parameters       []*Parameters //Parameters is a slice of request parameters for a single endpoint Optional.
	Returns          []*Returns    //what kind of response this API returns. Optional.
	Read             interface{}   //Read tells what resource type will be read from the request payload. Optional.
	Consumes         []string      //Consumes specifies that this WebService can consume one or more MIME types.
	Produces         []string      //Produces specifies that this WebService can produce one or more MIME types.
}

//Returns describe response doc
type Returns struct {
	Code    int // http response code
	Message string
	Model   interface{} // response body structure
}

//Parameters describe parameters in url path or query params
type Parameters struct {
	Name      string //parameter name
	DataType  string // string, int etc
	ParamType int    //restful.QueryParameterKind or restful.PathParameterKind
	Desc      string
}

//GetRouteSpecs is to return a rest API specification of a go struct
func GetRouteSpecs(schema interface{}) ([]Route, error) {
	rfValue := reflect.ValueOf(schema)
	name := reflect.Indirect(rfValue).Type().Name()
	urlPatternFunc := rfValue.MethodByName("URLPatterns")
	if !urlPatternFunc.IsValid() {
		return []Route{}, fmt.Errorf("<rest.RegisterResource> no 'URLPatterns' function in servant struct `%s`", name)
	}
	vals := urlPatternFunc.Call([]reflect.Value{})
	if len(vals) <= 0 {
		return []Route{}, fmt.Errorf("<rest.RegisterResource> call URLPatterns function failed in struct `%s`", name)
	}

	if !rfValue.CanInterface() {
		return []Route{}, fmt.Errorf("<rest.RegisterResource> result of 'URLPatterns' function not interface type in servant struct `%s`", name)
	}

	if routes, ok := vals[0].Interface().([]Route); ok {
		return routes, nil
	}
	return []Route{}, fmt.Errorf("<rest.RegisterResource> result of 'URLPatterns' function not []*Route type in servant struct `%s`", name)
}

//WrapHandlerChain wrap business handler with handler chain
func WrapHandlerChain(route Route, schemaType reflect.Type, schemaValue reflect.Value, schemaName string,
	opts server.Options) (restful.RouteFunction, error) {
	openlogging.GetLogger().Infof("add route path: [%s] method: [%s] func: [%s]. ", route.Path, route.Method, route.ResourceFuncName)
	method, exist := schemaType.MethodByName(route.ResourceFuncName)
	if !exist {
		openlogging.GetLogger().Errorf("router func can not find: %s", route.ResourceFuncName)
		return nil, fmt.Errorf("router func can not find: %s", route.ResourceFuncName)
	}

	handler := func(req *restful.Request, rep *restful.Response) {
		c, err := handler.GetChain(common.Provider, opts.ChainName)
		if err != nil {
			openlogging.GetLogger().Errorf("handler chain init err [%s]", err.Error())
			rep.AddHeader("Content-Type", "text/plain")
			rep.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}
		inv, err := HTTPRequest2Invocation(req, schemaName, method.Name)
		if err != nil {
			openlogging.GetLogger().Errorf("transfer http request to invocation failed, err [%s]", err.Error())
			return
		}
		//give inv.Ctx to user handlers, modules may inject headers in handler chain
		c.Next(inv, func(ir *invocation.Response) error {
			if ir.Err != nil {
				if rep != nil {
					rep.WriteHeader(ir.Status)
				}
				return ir.Err
			}
			Invocation2HTTPRequest(inv, req)

			bs := NewBaseServer(inv.Ctx)
			bs.Req = req
			bs.Resp = rep
			ir.Status = bs.Resp.StatusCode()
			// check body size
			if opts.BodyLimit > 0 {
				bs.Req.Request.Body = http.MaxBytesReader(bs.Resp, bs.Req.Request.Body, opts.BodyLimit)
			}
			method.Func.Call([]reflect.Value{schemaValue, reflect.ValueOf(bs)})

			if bs.Resp.StatusCode() >= http.StatusBadRequest {
				return fmt.Errorf("get err from http handle, get status: %d", bs.Resp.StatusCode())
			}
			return nil
		})

	}

	return handler, nil
}
