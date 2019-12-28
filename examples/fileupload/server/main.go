package main

import (
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core/lager"
	example "github.com/go-chassis/go-chassis/examples/fileupload/server/schemas"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/fileupload/server/
func main() {
	chassis.RegisterRestSchema(&example.RestFulUpload{})

	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
