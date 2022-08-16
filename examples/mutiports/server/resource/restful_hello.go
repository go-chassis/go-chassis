package resource

import (
	"net/http"

	rf "github.com/go-chassis/go-chassis/v2/server/restful"
)

// RestFulHello is a struct used for implementation of restfull hello program
type RestFulHello struct {
}

// Hello
func (r *RestFulHello) Hello(b *rf.Context) {
	b.Write([]byte("hi from hello"))
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/hello", ResourceFunc: r.Hello},
	}
}

// Legacy is a struct
type Legacy struct {
}

// Do
func (r *Legacy) Do(b *rf.Context) {
	b.Write([]byte("hello from legacy"))
}

// URLPatterns helps to respond for corresponding API calls
func (r *Legacy) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/legacy", ResourceFunc: r.Do},
	}
}

// Legacy is a struct
type Admin struct {
}

// Do
func (r *Admin) Do(b *rf.Context) {
	b.Write([]byte("hello from admin"))
}

// URLPatterns helps to respond for corresponding API calls
func (r *Admin) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/admin", ResourceFunc: r.Do},
	}
}
