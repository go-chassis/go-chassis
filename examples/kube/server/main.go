package main

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/examples/schemas"

	_ "github.com/ServiceComb/go-chassis/bootstrap"
	_ "github.com/go-chassis/go-chassis-plugins/registry/kube"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/

func main() {
	chassis.RegisterSchema("rest", &schemas.Hello{}, server.WithSchemaID("HelloService"))
	chassis.RegisterSchema("rest-legacy", &schemas.Legacy{}, server.WithSchemaID("LegacyService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}
	chassis.Run()
}
