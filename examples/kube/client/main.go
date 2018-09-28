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
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}
	for {
		req, err := rest.NewRequest("GET", "cse://kubeserver/hello", nil)
		if err != nil {
			lager.Logger.Error("new request failed." + err.Error())
		}

		resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req)
		if err != nil {
			lager.Logger.Error("do request failed." + err.Error())
		}
		defer resp.Close()
		lager.Logger.Info("REST Server sayhello[GET]: " + string(httputil.ReadBody(resp)))
		time.Sleep(1 * time.Second)

		req, err = rest.NewRequest("GET", "cse://kubeserver:legacy/legacy", nil)
		if err != nil {
			lager.Logger.Error("new request failed." + err.Error())
		}

		resp, err = core.NewRestInvoker().ContextDo(context.TODO(), req)
		if err != nil {
			lager.Logger.Error("do request failed." + err.Error())
		}
		defer resp.Body.Close()
		lager.Logger.Info("REST Server sayhello[GET]: " + string(httputil.ReadBody(resp)))
		time.Sleep(1 * time.Second)
	}
}
