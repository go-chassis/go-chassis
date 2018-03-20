package main

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	_ "github.com/ServiceComb/go-chassis/examples/plugin/handler"
	"github.com/ServiceComb/go-chassis/examples/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	chassis.RegisterSchema("highway", &schemas.HelloServer{}, server.WithSchemaID("HelloService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}
	chassis.Run()
}
