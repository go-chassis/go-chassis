package main

import (
	"github.com/go-chassis/go-chassis"
	rf "github.com/go-chassis/go-chassis/server/restful"
	"github.com/go-mesh/openlogging"
	"net/http"

	_ "github.com/go-chassis/go-chassis/middleware/ratelimiter"
)

type Hello struct{}

//Hello
func (r *Hello) Hello(b *rf.Context) { b.Write([]byte("hi from hello")) }

//URLPatterns helps to respond for corresponding API calls
func (r *Hello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/hello", ResourceFunc: r.Hello},
	}
}

func main() {
	chassis.RegisterSchema("rest", &Hello{})
	if err := chassis.Init(); err != nil {
		openlogging.Fatal("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
