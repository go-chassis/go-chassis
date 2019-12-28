package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/examples/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rpc/server/

func main() {
	chassis.RegisterSchema("highway", &schemas.HelloServer{}, server.WithSchemaID("HelloService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.")
		return
	}
	chassis.Run()
}
