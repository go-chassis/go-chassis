package resource

import (
	"fmt"
	"net/http"

	rf "github.com/go-chassis/go-chassis/v2/server/restful"
)

//RestFulHello is a struct used for implementation of restful hello program
type RestFulHello struct {
}

//Health
func (r *RestFulHello) Health(b *rf.Context) {
	b.Write([]byte(fmt.Sprintf("handler chain set metadata %s,set header %s", b.ReadRestfulRequest().Attribute("auth"), b.ReadRequest().Header.Get("X-Auth"))))
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/health", ResourceFunc: r.Health},
	}
}
