package resource

import (
	"net/http"

	"fmt"
	rf "github.com/go-chassis/go-chassis/server/restful"
	"math/rand"
)

var num = rand.Intn(100)

//RestFulHello is a struct used for implementation of restfull hello program
type RestFulHello struct {
}

//Health
func (r *RestFulHello) Health(b *rf.Context) {
	b.Write([]byte(fmt.Sprintf("handler chain set metadata %s,set header %s", b.ReadRestfulRequest().Attribute("auth"), b.ReadRequest().Header.Get("X-Auth"))))
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/health", ResourceFuncName: "Health"},
	}
}
