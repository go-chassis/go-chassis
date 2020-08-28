package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/examples/schemas"
	"github.com/go-chassis/openlog"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/server/

func main() {
	chassis.RegisterSchema("rest", &schemas.Hello{})
	chassis.RegisterSchema("rest-legacy", &schemas.Legacy{})
	chassis.RegisterSchema("rest-admin", &schemas.Admin{})
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
