package main

import (
	"context"
	"time"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/lager"

	_ "github.com/go-chassis/go-chassis-plugins/registry/kube"
	_ "github.com/go-chassis/go-chassis/bootstrap"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}
	for {
		req, err := rest.NewRequest("GET", "cse://kubeserver/sayhello/world")
		if err != nil {
			lager.Logger.Error("new request failed.", err)
		}
		defer req.Close()

		resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
		if err != nil {
			lager.Logger.Error("do request failed.", err)
		}
		defer resp.Close()
		lager.Logger.Info("REST Server sayhello[GET]: " + string(resp.ReadBody()))
		time.Sleep(1 * time.Second)
	}
}
