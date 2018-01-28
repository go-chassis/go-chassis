package restful

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/metrics"
	microServer "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-swagger12"
	"golang.org/x/net/context"
)

// constants for metric path and name
const (
	//Name is a variable of type string which indicates the protocol being used
	Name              = "rest"
	DefaultMetricPath = "metrics"
)

func init() {
	restful.SetCacheReadEntity(false)
	server.InstallPlugin(Name, newRestfulServer)
}

type restfulServer struct {
	microServiceName string
	container        *restful.Container
	ws               *restful.WebService
	opts             microServer.Options
	mux              sync.RWMutex
	exit             chan chan error
	server           *http.Server
}

func newRestfulServer(opts ...microServer.Option) microServer.Server {
	options := newOptions(opts...)
	ws := new(restful.WebService)
	ws.Path("/").Doc("root path").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	if archaius.GetBool("cse.metrics.enable", false) {

		metricPath := archaius.GetString("cse.metrics.apiPath", DefaultMetricPath)
		if !strings.HasPrefix(metricPath, "/") {
			metricPath = "/" + metricPath
		}
		lager.Logger.Info("Enbaled metrics API on " + metricPath)
		ws.Route(ws.GET(metricPath).To(metrics.MetricsHandleFunc))
	}
	return &restfulServer{
		opts:      options,
		container: restful.NewContainer(),
		ws:        ws,
	}
}
func newOptions(opt ...microServer.Option) microServer.Options {
	opts := microServer.Options{
		Metadata: map[string]string{},
	}

	for _, o := range opt {
		o(&opts)
	}
	return opts
}

func (r *restfulServer) Init(opts ...microServer.Option) error {
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, opt := range opts {
		opt(&r.opts)
	}
	lager.Logger.Info("Rest server init success")
	return nil
}

func (r *restfulServer) Options() microServer.Options {
	r.mux.RLock()
	defer r.mux.RUnlock()
	opts := r.opts
	return opts
}

func (r *restfulServer) Register(schema interface{}, options ...microServer.RegisterOption) (string, error) {
	lager.Logger.Info("register rest server")
	opts := microServer.RegisterOptions{}
	r.mux.Lock()
	defer r.mux.Unlock()
	for _, o := range options {
		o(&opts)
	}
	if opts.MicroServiceName == "" {
		opts.MicroServiceName = config.SelfServiceName
	}
	r.microServiceName = opts.MicroServiceName
	routes, err := GetRoutes(schema)
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
			lager.Logger.Errorf(nil, "router func can not find: %s", route.ResourceFuncName)
			return "", fmt.Errorf("router func can not find: %s", route.ResourceFuncName)
		}

		handle := func(req *restful.Request, rep *restful.Response) {
			c, err := handler.GetChain(common.Provider, r.opts.ChainName)
			if err != nil {
				lager.Logger.Errorf(err, "Handler chain init err")
				rep.AddHeader("Content-Type", "text/plain")
				rep.WriteErrorString(http.StatusInternalServerError, err.Error())
				return
			}
			//todo: use it for hystric
			inv := invocation.Invocation{
				MicroServiceName:   r.microServiceName,
				SourceMicroService: req.HeaderParameter(common.HeaderSourceName),
				Args:               req,
				Protocol:           common.ProtocolRest,
				SchemaID:           schemaName,
				OperationID:        method.Name,
			}
			bs := NewBaseServer(context.TODO())
			bs.req = req
			bs.resp = rep
			c.Next(&inv, func(ir *invocation.InvocationResponse) error {
				if ir.Err != nil {
					return ir.Err
				}
				method.Func.Call([]reflect.Value{schemaValue, reflect.ValueOf(bs)})
				if bs.resp.StatusCode() >= http.StatusBadRequest {
					return fmt.Errorf("get err from http handle, get status: %d", bs.resp.StatusCode())
				}
				return nil
			})

		}

		switch route.Method {
		case http.MethodGet:
			r.ws.Route(r.ws.GET(route.Path).To(handle).Doc(route.ResourceFuncName).Operation(route.ResourceFuncName))
		case http.MethodPost:
			r.ws.Route(r.ws.POST(route.Path).To(handle).Doc(route.ResourceFuncName).Operation(route.ResourceFuncName))
		case http.MethodHead:
			r.ws.Route(r.ws.HEAD(route.Path).To(handle).Doc(route.ResourceFuncName).Operation(route.ResourceFuncName))
		case http.MethodPut:
			r.ws.Route(r.ws.PUT(route.Path).To(handle).Doc(route.ResourceFuncName).Operation(route.ResourceFuncName))
		case http.MethodPatch:
			r.ws.Route(r.ws.PATCH(route.Path).To(handle).Doc(route.ResourceFuncName).Operation(route.ResourceFuncName))
		case http.MethodDelete:
			r.ws.Route(r.ws.DELETE(route.Path).To(handle).Doc(route.ResourceFuncName).Operation(route.ResourceFuncName))
		default:
			return "", errors.New("method do not support: " + route.Method)
		}
	}
	return reflect.TypeOf(schema).String(), nil
}

func (r *restfulServer) Start() error {
	config := r.Options()
	r.mux.Lock()
	r.opts.Address = config.Address
	r.mux.Unlock()
	r.container.Add(r.ws)
	if r.opts.TLSConfig != nil {
		r.server = &http.Server{Addr: config.Address, Handler: r.container, TLSConfig: r.opts.TLSConfig}
	} else {
		r.server = &http.Server{Addr: config.Address, Handler: r.container}
	}
	// TODO Choose a suitable strategy of transforming code to contract
	// register to swagger ui
	if val := os.Getenv("GO_CHASSIS_SWAGGERFILEPATH"); val != "" {
		swaggerConfig := swagger.Config{
			WebServices:    r.container.RegisteredWebServices(),
			WebServicesUrl: config.Address,
			ApiPath:        "/apidocs.json",
			// Optionally, specify where the UI is located
			SwaggerPath:     "/apidocs/",
			SwaggerFilePath: val}
		swagger.RegisterSwaggerService(swaggerConfig, r.container)
	}

	var err error
	go func() {
		if r.server.TLSConfig != nil {
			err = r.server.ListenAndServeTLS("", "")
		} else {
			err = r.server.ListenAndServe()
		}
		if err != nil {
			lager.Logger.Error("Can't start http server", err)
			server.ServerErr <- err
		}

	}()

	lager.Logger.Warnf(nil, "Restful server listening on: %s", config.Address)
	return nil
}

func (r *restfulServer) Stop() error {
	//only golang 1.8 is support graceful shutdown.
	if err := r.server.Shutdown(nil); err != nil {
		return err // failure/timeout shutting down the server gracefully
	}
	return nil
}

func (r *restfulServer) String() string {
	return Name
}
