package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/examples/schemas"
	_ "github.com/go-mesh/mesher/plugins/registry/istiov2"
	"github.com/go-mesh/openlogging"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/

func main() {
	chassis.RegisterSchema("rest", &schemas.RestFulHello{}, server.WithSchemaID("RestHelloService"))
	openlogging.SetLogger(lager.Logger)
	if err := chassis.Init(); err != nil {
		lager.Logger.Errorf("Init failed: %s", err.Error())
		return
	}
	chassis.Run()
}
