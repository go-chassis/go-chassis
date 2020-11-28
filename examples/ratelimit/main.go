package main

import (
	"github.com/go-chassis/go-chassis/v2"
	_ "github.com/go-chassis/go-chassis/v2/middleware/ratelimiter"
	"github.com/go-chassis/go-chassis/v2/server/restful"
	"github.com/go-chassis/openlog"
	"net/http"
)

type DemoResource struct {
}

func (r *DemoResource) Limit(b *restful.Context) {
	b.ReadResponseWriter().WriteHeader(http.StatusOK)
	b.ReadResponseWriter().Write([]byte("ok"))
}

// URLPatterns returns routes
func (r *DemoResource) URLPatterns() []restful.Route {
	return []restful.Route{
		{Method: http.MethodGet, Path: "/limit", ResourceFunc: r.Limit},
	}
}

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/{project_root}/

func main() {
	chassis.RegisterSchema("rest", &DemoResource{})
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
