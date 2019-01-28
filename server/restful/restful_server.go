package restful

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/metrics"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-chassis/core/config/schema"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-restful-swagger20"
	"github.com/go-mesh/openlogging"
	"os"
	"path/filepath"
)

// constants for metric path and name
const (
	//Name is a variable of type string which indicates the protocol being used
	Name              = "rest"
	DefaultMetricPath = "metrics"
	MimeFile          = "application/octet-stream"
	MimeMult          = "multipart/form-data"
)

func init() {
	server.InstallPlugin(Name, newRestfulServer)
}

type restfulServer struct {
	microServiceName string
	container        *restful.Container
	ws               *restful.WebService
	opts             server.Options
	mux              sync.RWMutex
	exit             chan chan error
	server           *http.Server
}

func newRestfulServer(opts server.Options) server.ProtocolServer {
	ws := new(restful.WebService)
	ws.Path("/").Doc("root path").
		Consumes(restful.MIME_XML, restful.MIME_JSON, MimeFile, MimeMult).
		Produces(restful.MIME_JSON, restful.MIME_XML, MimeFile, MimeMult) // you can specify this per route as well

	if archaius.GetBool("cse.metrics.enable", false) {

		metricPath := archaius.GetString("cse.metrics.apiPath", DefaultMetricPath)
		if !strings.HasPrefix(metricPath, "/") {
			metricPath = "/" + metricPath
		}
		lager.Logger.Info("Enabled metrics API on " + metricPath)
		ws.Route(ws.GET(metricPath).To(metrics.HTTPHandleFunc))
	}
	return &restfulServer{
		opts:      opts,
		container: restful.NewContainer(),
		ws:        ws,
	}
}
func httpRequest2Invocation(req *restful.Request, schema, operation string) (*invocation.Invocation, error) {

	inv := &invocation.Invocation{
		MicroServiceName:   runtime.ServiceName,
		SourceMicroService: req.HeaderParameter(common.HeaderSourceName),
		Args:               req,
		Protocol:           common.ProtocolRest,
		SchemaID:           schema,
		OperationID:        operation,
		URLPathFormat:      req.Request.URL.Path,
		Metadata: map[string]interface{}{
			common.RestMethod: req.Request.Method,
		},
	}
	//set headers to Ctx, then user do not  need to consider about protocol in handlers
	m := make(map[string]string, 0)
	inv.Ctx = context.WithValue(context.Background(), common.ContextHeaderKey{}, m)
	for k := range req.Request.Header {
		m[k] = req.Request.Header.Get(k)
	}
	return inv, nil
}
func (r *restfulServer) Register(schema interface{}, options ...server.RegisterOption) (string, error) {
	lager.Logger.Info("register rest server")
	opts := server.RegisterOptions{}
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, o := range options {
		o(&opts)
	}

	routes, err := GetRouteSpecs(schema)
	if err != nil {
		return "", err
	}
	schemaType := reflect.TypeOf(schema)
	schemaValue := reflect.ValueOf(schema)
	var schemaName string
	tokens := strings.Split(schemaType.String(), ".")
	if len(tokens) >= 1 {
		schemaName = tokens[len(tokens)-1]
	}
	lager.Logger.Infof("schema registered is [%s]", schemaName)
	for _, route := range routes {
		lager.Logger.Infof("Add route path: [%s] Method: [%s] Func: [%s]. ", route.Path, route.Method, route.ResourceFuncName)
		method, exist := schemaType.MethodByName(route.ResourceFuncName)
		if !exist {
			lager.Logger.Errorf("router func can not find: %s", route.ResourceFuncName)
			return "", fmt.Errorf("router func can not find: %s", route.ResourceFuncName)
		}

		handler := func(req *restful.Request, rep *restful.Response) {
			c, err := handler.GetChain(common.Provider, r.opts.ChainName)
			if err != nil {
				lager.Logger.Errorf("Handler chain init err [%s]", err.Error())
				rep.AddHeader("Content-Type", "text/plain")
				rep.WriteErrorString(http.StatusInternalServerError, err.Error())
				return
			}
			inv, err := httpRequest2Invocation(req, schemaName, method.Name)
			if err != nil {
				lager.Logger.Errorf("transfer http request to invocation failed, err [%s]", err.Error())
				return
			}
			//give inv.Ctx to user handlers, modules may inject headers in handler chain

			c.Next(inv, func(ir *invocation.Response) error {
				if ir.Err != nil {
					return ir.Err
				}
				transfer(inv, req)
				bs := NewBaseServer(inv.Ctx)
				bs.req = req
				bs.resp = rep
				ir.Status = bs.resp.StatusCode()
				method.Func.Call([]reflect.Value{schemaValue, reflect.ValueOf(bs)})
				if bs.resp.StatusCode() >= http.StatusBadRequest {
					return fmt.Errorf("get err from http handle, get status: %d", bs.resp.StatusCode())
				}
				return nil
			})

		}

		if err := r.register2GoRestful(route, handler); err != nil {
			return "", err
		}
	}
	return reflect.TypeOf(schema).String(), nil
}
func transfer(inv *invocation.Invocation, req *restful.Request) {
	for k, v := range inv.Metadata {
		req.SetAttribute(k, v.(string))
	}
	m := common.FromContext(inv.Ctx)
	for k, v := range m {
		req.Request.Header.Set(k, v)
	}

}
func (r *restfulServer) register2GoRestful(routeSpec Route, handler restful.RouteFunction) error {
	var rb *restful.RouteBuilder
	switch routeSpec.Method {
	case http.MethodGet:
		rb = r.ws.GET(routeSpec.Path)
	case http.MethodPost:
		rb = r.ws.POST(routeSpec.Path)
	case http.MethodHead:
		rb = r.ws.HEAD(routeSpec.Path)
	case http.MethodPut:
		rb = r.ws.PUT(routeSpec.Path)
	case http.MethodPatch:
		rb = r.ws.PATCH(routeSpec.Path)
	case http.MethodDelete:
		rb = r.ws.DELETE(routeSpec.Path)
	default:
		return errors.New("method [" + routeSpec.Method + "] do not support")
	}
	for _, param := range routeSpec.Parameters {
		switch param.ParamType {
		case restful.QueryParameterKind:
			rb = rb.Param(restful.QueryParameter(param.Name, param.Desc).DataType(param.DataType))
		case restful.PathParameterKind:
			rb = rb.Param(restful.PathParameter(param.Name, param.Desc).DataType(param.DataType))
		}

	}
	for _, r := range routeSpec.Returns {
		rb = rb.Returns(r.Code, r.Message, r.Model)
	}
	if routeSpec.Read != nil {
		rb = rb.Reads(routeSpec.Read)
	}
	r.ws.Route(rb.To(handler).Doc(routeSpec.FuncDesc).Operation(routeSpec.FuncDesc))
	return nil
}
func (r *restfulServer) Start() error {
	var err error
	config := r.opts
	r.mux.Lock()
	r.opts.Address = config.Address
	r.mux.Unlock()
	r.container.Add(r.ws)
	if r.opts.TLSConfig != nil {
		r.server = &http.Server{Addr: config.Address, Handler: r.container, TLSConfig: r.opts.TLSConfig}
	} else {
		r.server = &http.Server{Addr: config.Address, Handler: r.container}
	}
	// register to swagger ui
	var path string
	if path = schema.GetSchemaPath(runtime.ServiceName); path == "" {
		return errors.New("schema path is empty")
	}
	if err = os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to generate swagger doc: %s", err.Error())
	}
	if err = os.MkdirAll(path, 0760); err != nil {
		return fmt.Errorf("failed to generate swagger doc: %s", err.Error())
	}
	swagger.LogInfo = func(format string, v ...interface{}) {
		openlogging.GetLogger().Infof(format, v...)
	}
	swaggerConfig := swagger.Config{
		WebServices:     r.container.RegisteredWebServices(),
		WebServicesUrl:  config.Address,
		ApiPath:         "/apidocs.json",
		FileStyle:       "yaml",
		SwaggerFilePath: filepath.Join(path, runtime.ServiceName+".yaml")}
	swagger.RegisterSwaggerService(swaggerConfig, r.container)
	go func() {
		if r.server.TLSConfig != nil {
			err = r.server.ListenAndServeTLS("", "")
		} else {
			err = r.server.ListenAndServe()
		}
		if err != nil {
			openlogging.GetLogger().Error("http server err: " + err.Error())
			server.ErrRuntime <- err
		}

	}()

	lager.Logger.Infof("Restful server listening on: %s", config.Address)
	return nil
}

func (r *restfulServer) Stop() error {
	if r.server == nil {
		openlogging.Info("http server never started")
		return nil
	}
	//only golang 1.8 support graceful shutdown.
	if err := r.server.Shutdown(nil); err != nil {
		return err // failure/timeout shutting down the server gracefully
	}
	return nil
}

func (r *restfulServer) String() string {
	return Name
}
