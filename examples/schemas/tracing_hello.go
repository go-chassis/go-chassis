package schemas

import (
	"net/http"

	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
	rf "github.com/go-chassis/go-chassis/server/restful"
	"log"
)

//TracingHello is a struct
type TracingHello struct {
}

//Trace is a method
func (r *TracingHello) Trace(b *rf.Context) {
	log.Println("tracing===", b.Ctx)
	req, err := rest.NewRequest("GET", "http://RESTServerB/sayhello/world", nil)
	if err != nil {
		b.WriteError(500, err)
		return
	}

	resp, err := core.NewRestInvoker().ContextDo(b.Ctx, req)
	if err != nil {
		b.WriteError(500, err)
		return
	}
	defer resp.Body.Close()
	b.Write(httputil.ReadBody(resp))
}

//URLPatterns helps to respond for corresponding API calls
func (r *TracingHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{http.MethodGet, "/trace", "Trace"},
	}
}
