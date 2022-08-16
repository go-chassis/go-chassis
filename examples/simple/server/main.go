package main

import (
	"fmt"
	"github.com/go-chassis/go-chassis/v2"
	rf "github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
	"net/http"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/server/

type RestFulHello struct {
}

func (r *RestFulHello) Root(b *rf.Context) {
	b.Write([]byte(fmt.Sprintf("hello %s", b.ReadRequest().RemoteAddr)))
}

// URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/hello", ResourceFunc: r.Root,
			Returns: []*rf.Returns{{Code: 200}}},
	}
}
func main() {
	chassis.RegisterSchema("rest", &RestFulHello{})
	if err := chassis.Init(); err != nil {
		openlog.Fatal("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
