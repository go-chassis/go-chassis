package main

import (
	"context"
	"time"

	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/lager"

	_ "github.com/go-chassis/go-chassis/bootstrap"
	_ "github.com/go-mesh/mesher/plugins/registry/istiov2"
	"github.com/go-mesh/openlogging"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Errorf("Init failed: %s", err.Error())
		return
	}
	openlogging.SetLogger(lager.Logger)
	for {
		req, err := rest.NewRequest("GET", "cse://pilotv2server/sayhello/world")
		if err != nil {
			lager.Logger.Errorf("new request failed.", err)
		}
		defer req.Close()

		resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
		if err != nil {
			lager.Logger.Errorf("do request failed.", err)
		}
		defer resp.Close()
		lager.Logger.Infof("REST Server sayhello[GET]: %s", string(resp.ReadBody()))
		time.Sleep(5 * time.Second)
	}
}
