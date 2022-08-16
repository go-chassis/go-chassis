package main

import (
	"github.com/go-chassis/go-chassis/v2"
	example "github.com/go-chassis/go-chassis/v2/examples/fileupload/server/schemas"
	"github.com/go-chassis/openlog"
)

// if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/fileupload/server/
func main() {
	chassis.RegisterSchema("rest", &example.RestFulUpload{})

	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	chassis.Run()
}
