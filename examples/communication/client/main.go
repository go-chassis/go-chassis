package main

import (
	"github.com/ServiceComb/go-chassis"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
	_ "github.com/ServiceComb/go-chassis/examples/plugin/handler"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"golang.org/x/net/context"
	"log"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	// just init client
	chassis.Init()
	// specify chain name
	invoker := core.NewRPCInvoker(core.ChainName("custom"))
	// new response object
	reply := &helloworld.HelloReply{}
	// create context with metadata
	ctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User": "tianxiaoliang",
	})
	err := invoker.Invoke(ctx, "SimpleServer", "HelloServer", "SayHello", &helloworld.HelloRequest{Name: "Peter"}, reply, core.WithEndpoint("127.0.0.1:9901"), core.WithProtocol("highway"))
	if err != nil {
		lager.Logger.Errorf(err, "Invoke failed.")
	}
	log.Println("reply -----------", reply)
}
