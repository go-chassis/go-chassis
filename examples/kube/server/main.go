package main

import (
	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis-plugins/registry/kube"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/examples/schemas"

	_ "github.com/go-chassis/go-chassis-plugins/registry/kube"
	_ "github.com/go-chassis/go-chassis/bootstrap"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/

func main() {
	chassis.RegisterSchema("rest", &schemas.Hello{}, server.WithSchemaID("HelloService"))
	chassis.RegisterSchema("rest-legacy", &schemas.Legacy{}, server.WithSchemaID("LegacyService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
