package main

import (
	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/examples/plugin/handler"
	_ "github.com/ServiceComb/go-chassis/examples/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	chassis.Init()
	chassis.Run()
}
