package main

import (
	"context"
	"github.com/go-chassis/openlog"

	"github.com/go-chassis/go-chassis/v2"
	_ "github.com/go-chassis/go-chassis/v2/bootstrap"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}

	req, err := rest.NewRequest("GET", "http://127.0.0.1:5001/hello", nil)
	if err != nil {
		openlog.Error("new request failed.")
		return
	}
	resp, err := core.NewRestInvoker().ContextDo(context.TODO(), req, core.WithoutSD())
	if err != nil {
		openlog.Error("do request failed.")
		return
	}
	defer resp.Body.Close()
	openlog.Info("REST Server sayhello[GET]: " + string(httputil.ReadBody(resp)))
}
