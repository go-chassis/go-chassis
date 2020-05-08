package restful

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-chassis/go-chassis/server/restful/api"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/emicklei/go-restful"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/core/common"
	globalconfig "github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/config/schema"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/registry"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/pkg/profile"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/pkg/util/iputil"
	swagger "github.com/go-chassis/go-restful-swagger20"
	"github.com/go-mesh/openlogging"
)

// constants for metric path and name
const (
	//Name is a variable of type string which indicates the protocol being used
	Name                    = "rest"
	DefaultMetricPath       = "metrics"
	DefaultProfilePath      = "profile"
	ProfileRouteRuleSubPath = "route-rule"
	ProfileDiscoverySubPath = "discovery"
	MimeFile                = "application/octet-stream"
	MimeMult                = "multipart/form-data"
)

const openTLS = "?sslEnabled=true"

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
	if archaius.GetBool("cse.metrics.enable", false) {
		metricPath := archaius.GetString("cse.metrics.apiPath", DefaultMetricPath)
		if !strings.HasPrefix(metricPath, "/") {
			metricPath = "/" + metricPath
		}
		openlogging.Info("Enabled metrics API on " + metricPath)
		ws.Route(ws.GET(metricPath).To(api.PrometheusHandleFunc))
	}
	addProfileRoutes(ws)
	return &restfulServer{
		opts:      opts,
		container: restful.NewContainer(),
		ws:        ws,
	}
}

func addProfileRoutes(ws *restful.WebService) {
	if !archaius.GetBool("cse.profile.enable", false) {
		return
	}
	profilePath := archaius.GetString("cse.profile.apiPath", DefaultProfilePath)
	if !strings.HasPrefix(profilePath, "/") {
		profilePath = "/" + profilePath
	}

	openlogging.Info("Enabled profile API on " + profilePath)
	ws.Route(ws.GET(profilePath).To(profile.HTTPHandleProfileFunc))

	profileRouteRulePath := profilePath + "/" + ProfileRouteRuleSubPath
	openlogging.Info("Enabled profile route-rule API on " + profileRouteRulePath)
	ws.Route(ws.GET(profileRouteRulePath).To(profile.HTTPHandleRouteRuleFunc))

	profileDiscoveryPath := profilePath + "/" + ProfileDiscoverySubPath
	openlogging.Info("Enabled profile discovery API on " + profileDiscoveryPath)
	ws.Route(ws.GET(profileDiscoveryPath).To(profile.HTTPHandleDiscoveryFunc))
}

// HTTPRequest2Invocation convert http request to uniform invocation data format
func HTTPRequest2Invocation(req *restful.Request, schema, operation string, resp *restful.Response) (*invocation.Invocation, error) {
	inv := &invocation.Invocation{
		MicroServiceName:   runtime.ServiceName,
		SourceMicroService: common.GetXCSEContext(common.HeaderSourceName, req.Request),
		Args:               req,
		Reply:              resp,
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
	openlogging.Info("register rest server")
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

	var schemaName string
	tokens := strings.Split(reflect.TypeOf(schema).String(), ".")
	if len(tokens) >= 1 {
		schemaName = tokens[len(tokens)-1]
	}
	openlogging.GetLogger().Infof("schema registered is [%s]", schemaName)
	for k := range routes {
		GroupRoutePath(&routes[k], schema)
		handler, err := WrapHandlerChain(&routes[k], schema, schemaName, r.opts)
		if err != nil {
			return "", err
		}
		if err := Register2GoRestful(routes[k], r.ws, handler); err != nil {
			return "", err
		}
	}
	return reflect.TypeOf(schema).String(), nil
}

// Invocation2HTTPRequest convert invocation back to http request, set down all meta data
func Invocation2HTTPRequest(inv *invocation.Invocation, req *restful.Request) {
	for k, v := range inv.Metadata {
		req.SetAttribute(k, v.(string))
	}
	m := common.FromContext(inv.Ctx)
	for k, v := range m {
		req.Request.Header.Set(k, v)
	}

}

//Register2GoRestful register http handler to go-restful framework
func Register2GoRestful(routeSpec Route, ws *restful.WebService, handler restful.RouteFunction) error {
	var rb *restful.RouteBuilder
	switch routeSpec.Method {
	case http.MethodGet:
		rb = ws.GET(routeSpec.Path)
	case http.MethodPost:
		rb = ws.POST(routeSpec.Path)
	case http.MethodHead:
		rb = ws.HEAD(routeSpec.Path)
	case http.MethodPut:
		rb = ws.PUT(routeSpec.Path)
	case http.MethodPatch:
		rb = ws.PATCH(routeSpec.Path)
	case http.MethodDelete:
		rb = ws.DELETE(routeSpec.Path)
	default:
		return errors.New("method [" + routeSpec.Method + "] do not support")
	}
	rb = fillParam(routeSpec, rb)
	for k, v := range routeSpec.Metadata {
		rb = rb.Metadata(k, v)
	}
	for _, r := range routeSpec.Returns {
		rb = rb.ReturnsWithHeaders(r.Code, r.Message, r.Model, r.Headers)
	}
	if routeSpec.Read != nil {
		rb = rb.Reads(routeSpec.Read)
	}

	if len(routeSpec.Consumes) > 0 {
		rb = rb.Consumes(routeSpec.Consumes...)
	} else {
		rb = rb.Consumes("*/*")
	}
	if len(routeSpec.Produces) > 0 {
		rb = rb.Produces(routeSpec.Produces...)
	} else {
		rb = rb.Produces("*/*")
	}
	ws.Route(rb.To(handler).Doc(routeSpec.FuncDesc).Operation(routeSpec.ResourceFuncName))

	return nil
}

//fillParam is for handle parameter by type
func fillParam(routeSpec Route, rb *restful.RouteBuilder) *restful.RouteBuilder {
	for _, param := range routeSpec.Parameters {
		p := &restful.Parameter{}
		switch param.ParamType {
		case restful.QueryParameterKind:
			p = restful.QueryParameter(param.Name, param.Desc)
		case restful.PathParameterKind:
			p = restful.PathParameter(param.Name, param.Desc)
		case restful.HeaderParameterKind:
			p = restful.HeaderParameter(param.Name, param.Desc)
		case restful.BodyParameterKind:
			p = restful.BodyParameter(param.Name, param.Desc)
		case restful.FormParameterKind:
			p = restful.FormParameter(param.Name, param.Desc)
		}
		rb = rb.Param(p.Required(param.Required).DataType(param.DataType))

	}
	return rb
}
func (r *restfulServer) Start() error {
	var err error
	config := r.opts
	r.mux.Lock()
	r.opts.Address = config.Address
	r.mux.Unlock()
	r.container.Add(r.ws)
	sslFlag := ""
	r.server = &http.Server{
		Addr:         config.Address,
		Handler:      r.container,
		ReadTimeout:  r.opts.Timeout,
		WriteTimeout: r.opts.Timeout,
		IdleTimeout:  r.opts.Timeout,
	}
	if r.opts.HeaderLimit > 0 {
		r.server.MaxHeaderBytes = r.opts.HeaderLimit
	}
	if r.opts.TLSConfig != nil {
		r.server.TLSConfig = r.opts.TLSConfig
		sslFlag = openTLS
	}
	// create schema
	err = r.CreateLocalSchema(config)
	if err != nil {
		return err
	}
	l, lIP, lPort, err := iputil.StartListener(config.Address, config.TLSConfig)

	if err != nil {
		return fmt.Errorf("failed to start listener: %s", err.Error())
	}

	registry.InstanceEndpoints[config.ProtocolServerName] = net.JoinHostPort(lIP, lPort) + sslFlag

	go func() {
		err = r.server.Serve(l)
		if err != nil {
			openlogging.Error("http server err: " + err.Error())
			server.ErrRuntime <- err
		}

	}()

	openlogging.GetLogger().Infof("http server is listening at %s", registry.InstanceEndpoints[config.ProtocolServerName])
	return nil
}

//register to swagger ui,Whether to create a schema, you need to refer to the configuration.
func (r *restfulServer) CreateLocalSchema(opts server.Options) error {
	var path string
	if path = schema.GetSchemaPath(runtime.ServiceName); path == "" {
		return errors.New("schema path is empty")
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to generate swagger doc: %s", err.Error())
	}
	if err := os.MkdirAll(path, 0760); err != nil {
		return fmt.Errorf("failed to generate swagger doc: %s", err.Error())
	}
	swagger.LogInfo = func(format string, v ...interface{}) {
		openlogging.GetLogger().Infof(format, v...)
	}
	swaggerConfig := swagger.Config{
		WebServices:     r.container.RegisteredWebServices(),
		WebServicesUrl:  opts.Address,
		ApiPath:         "/apidocs.json",
		SwaggerPath:     "/swagger/",
		FileStyle:       "yaml",
		OpenService:     true,
		SwaggerFilePath: "./swagger-ui/dist/",
	}
	if globalconfig.GlobalDefinition.Cse.NoRefreshSchema {
		openlogging.Info("will not create schema file. if you want to change it, please update chassis.yaml->NoRefreshSchema=true")
	} else {
		swaggerConfig.OutFilePath = filepath.Join(path, runtime.ServiceName+".yaml")
	}
	sws := swagger.RegisterSwaggerService(swaggerConfig, r.container)
	openlogging.Info("contract has been created successfully. path:" + path)
	//set schema information when create local schema file
	err := schema.SetSchemaInfo(sws)
	if err != nil {
		return fmt.Errorf("set schema information,%s", err.Error())
	}
	return nil
}

func (r *restfulServer) Stop() error {
	if r.server == nil {
		openlogging.Info("http server never started")
		return nil
	}
	//only golang 1.8 support graceful shutdown.
	if err := r.server.Shutdown(context.TODO()); err != nil {
		openlogging.Warn("http shutdown error: " + err.Error())
		return err // failure/timeout shutting down the server gracefully
	}
	return nil
}

func (r *restfulServer) String() string {
	return Name
}
