package main

import (
	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/bootstrap"
	_ "github.com/go-chassis/go-chassis/configserver"
	"github.com/go-chassis/go-chassis/examples/schemas"
	"github.com/go-chassis/openlog"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	chassis.RegisterSchema("rest", &schemas.RestFulHello{})
	chassis.RegisterSchema("rest", &schemas.RestFulMessage{})
	//start all server you register in server/schemas.
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
