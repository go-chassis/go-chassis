package main

import (
	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/examples/schemas"
	_ "github.com/go-chassis/go-chassis/healthz/provider"
	_ "github.com/go-chassis/go-chassis/middleware/monitoring"
	"github.com/go-mesh/openlogging"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/discovery/server/
func main() {
	chassis.RegisterSchema("rest", &schemas.RestFulHello{})
	chassis.RegisterSchema("rest", &schemas.RestFulMessage{})
	//start all server you register in server/schemas.
	if err := chassis.Init(); err != nil {
		openlogging.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
