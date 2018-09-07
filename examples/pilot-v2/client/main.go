package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"

	_ "github.com/go-chassis/go-chassis/bootstrap"
	_ "github.com/go-mesh/mesher/plugins/registry/istiov2"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		fmt.Println("Init failed.", err)
		return
	}
	for {
		req, err := rest.NewRequest("GET", "cse://pilotv2server/sayhello/world")
		if err != nil {
			fmt.Println("new request failed.", err)
		}
		defer req.Close()

		resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
		if err != nil {
			fmt.Println("do request failed.", err)
		}
		defer resp.Close()
		fmt.Println("REST Server sayhello[GET]: " + string(resp.ReadBody()))
		time.Sleep(5 * time.Second)
	}
}
