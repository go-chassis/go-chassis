package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/examples/metadata/resource"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/discovery/server/
func main() {
	chassis.RegisterRestSchema(&resource.RestFulHello{})
	//start all server you register in server/schemas.
	if err := chassis.Init(); err != nil {
		panic(err)
		return
	}
	chassis.Run()
}
