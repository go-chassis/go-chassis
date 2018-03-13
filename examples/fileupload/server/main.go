package main

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core/lager"
	example "github.com/ServiceComb/go-chassis/examples/fileupload/server/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	chassis.RegisterSchema("rest", &example.RestFulUpload{})

	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}
	chassis.Run()
}
