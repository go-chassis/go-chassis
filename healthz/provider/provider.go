package provider

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/healthz/client"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	rf "github.com/go-chassis/go-chassis/server/restful"
)

var (
	once               sync.Once
	defaultHealthCheck = &HealthCheck{}
	checkResult        []byte
	checkReply         *client.Reply
)

func firstRequest() {
	once.Do(func() {
		checkReply = &client.Reply{
			AppID:       runtime.App,
			ServiceName: runtime.ServiceName,
			Version:     runtime.Version,
		}
		checkResult, _ = json.Marshal(checkReply)
	})
}

// HealthCheck is the struct defines provider health check
type HealthCheck struct {
}

// RestCheck returns status OK and self serviceName
func (hc *HealthCheck) RestCheck(ctx *rf.Context) {
	firstRequest()

	ctx.AddHeader("Content-Type", common.JSON)
	ctx.Write(checkResult)
}

// URLPatterns returns HealthCheck's routes
func (hc *HealthCheck) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/healthz", ResourceFunc: hc.RestCheck},
	}
}

func init() {
	chassis.RegisterSchema("rest", defaultHealthCheck, server.WithSchemaID("_chassis_rest_healthz"))
}
