package main

import (
	"context"
	"log"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/server"
	"github.com/go-chassis/go-chassis/examples/schemas"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/communication/client/
func main() {
	// just init client
	chassis.RegisterSchema("highway", &schemas.HelloServer{}, server.WithSchemaID("HelloService"))
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	// specify chain name
	invoker := core.NewRPCInvoker()
	// new response object
	reply := &helloworld.HelloReply{}
	// create context with metadata
	ctx := context.WithValue(context.Background(), common.ContextHeaderKey{}, map[string]string{
		"X-User": "tianxiaoliang",
	})
	err := invoker.Invoke(ctx, "SimpleServer", "HelloService", "SayHello", &helloworld.HelloRequest{Name: "Peter"}, reply, core.WithEndpoint("127.0.0.1:9901"), core.WithProtocol("highway"))
	if err != nil {
		lager.Logger.Errorf("Invoke failed.")
	}
	log.Println("reply -----------", reply)
}
