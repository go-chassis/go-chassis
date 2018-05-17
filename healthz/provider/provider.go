package provider

import (
	"context"
	"encoding/json"
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/healthz/client"
	rf "github.com/ServiceComb/go-chassis/server/restful"
	"net/http"
	"sync"
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
			AppId:       config.GlobalDefinition.AppID,
			ServiceName: config.SelfServiceName,
			Version:     config.SelfVersion,
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

// HighwayCheck returns status OK and self serviceName
func (hc *HealthCheck) HighwayCheck(_ context.Context, _ *client.Request) (*client.Reply, error) {
	firstRequest()

	return checkReply, nil
}

// URLPatterns returns HealthCheck's routes
func (hc *HealthCheck) URLPatterns() []rf.Route {
	return []rf.Route{
		{http.MethodGet, "/healthz", "RestCheck"},
	}
}

func init() {
	chassis.RegisterSchema(common.ProtocolRest, defaultHealthCheck, server.WithSchemaID("_chassis_rest_healthz"))
	chassis.RegisterSchema(common.ProtocolHighway, defaultHealthCheck, server.WithSchemaID("_chassis_highway_healthz"))
}
