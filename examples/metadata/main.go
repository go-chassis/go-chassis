package main

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/examples/metadata/resource"
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
