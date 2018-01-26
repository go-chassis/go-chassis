package main

import (
	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"golang.org/x/net/context"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}
	//declare reply struct
	reply := &helloworld.HelloReply{}
	//Invoke with microservice name, schema ID and operation ID
	if err := core.NewRPCInvoker().Invoke(context.Background(), "RPCServer", "HelloService", "SayHello",
		&helloworld.HelloRequest{Name: "Peter"}, reply); err != nil {
		lager.Logger.Error("error", err)
	}
	lager.Logger.Info(reply.Message)
}
