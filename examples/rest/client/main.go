package main

import (
	"context"

	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/lager"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}

	req, err := rest.NewRequest("GET", "cse://RESTServer/sayhello/world")
	if err != nil {
		lager.Logger.Error("new request failed.", err)
		return
	}
	defer req.Close()

	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
	if err != nil {
		lager.Logger.Error("do request failed.", err)
		return
	}
	defer resp.Close()
	lager.Logger.Info("REST Server sayhello[GET]: " + string(resp.ReadBody()))
}
