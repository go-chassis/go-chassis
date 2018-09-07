package main

import (
	"fmt"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/examples/schemas"
	_ "github.com/go-mesh/mesher/plugins/registry/istiov2"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/

func main() {
	chassis.RegisterSchema("rest", &schemas.RestFulHello{}, server.WithSchemaID("RestHelloService"))
	if err := chassis.Init(); err != nil {
		fmt.Println("Init failed.", err)
		return
	}
	chassis.Run()
}
