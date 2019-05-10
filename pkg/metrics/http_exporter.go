package metrics

import (
	"github.com/emicklei/go-restful"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTPHandleFunc is a go-restful handler which can expose metrics in http server
func HTTPHandleFunc(req *restful.Request, rep *restful.Response) {
	promhttp.HandlerFor(GetSystemPrometheusRegistry(), promhttp.HandlerOpts{}).ServeHTTP(rep.ResponseWriter, req.Request)
}
