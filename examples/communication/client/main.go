package main

import (
	"context"
	"log"

	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/examples/schemas"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	// just init client
	chassis.RegisterSchema("highway", &schemas.HelloServer{}, server.WithSchemaID("HelloService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}
	// specify chain name
	invoker := core.NewRPCInvoker()
	// new response object
	reply := &helloworld.HelloReply{}
	// create context with metadata
	ctx := context.WithValue(context.Background(), common.ContextValueKey{}, map[string]string{
		"X-User": "tianxiaoliang",
	})
	err := invoker.Invoke(ctx, "SimpleServer", "HelloService", "SayHello", &helloworld.HelloRequest{Name: "Peter"}, reply, core.WithEndpoint("127.0.0.1:9901"), core.WithProtocol("highway"))
	if err != nil {
		lager.Logger.Errorf(err, "Invoke failed.")
	}
	log.Println("reply -----------", reply)
}
