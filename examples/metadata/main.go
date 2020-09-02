package main

import (
	"github.com/go-chassis/go-chassis/v2"
	"github.com/go-chassis/go-chassis/v2/examples/metadata/resource"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/discovery/server/
func main() {
	chassis.RegisterSchema("rest", &resource.RestFulHello{})
	//start all server you register in server/schemas.
	if err := chassis.Init(); err != nil {
		panic(err)
		return
	}
	chassis.Run()
}
