package main

import (
	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	serverOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	chassis.RegisterSchema("highway", &schemas.HelloServer{}, serverOption.WithSchemaID("HelloServer"))
	chassis.RegisterSchema("highway", &schemas.EmployServer{}, serverOption.WithSchemaID("EmployServer"), serverOption.WithMicroServiceName("Server"))
	chassis.RegisterSchema("rest", &schemas.RestFulHello{})
	chassis.RegisterSchema("rest", &schemas.RestFulMessage{})
	//start all server you register in server/schemas.
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}
	chassis.Run()
}
