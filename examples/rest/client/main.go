package main

import (
	"context"
	"github.com/go-chassis/openlog"

	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}

	req, err := rest.NewRequest("GET", "http://RESTServer/sayhello/world", nil)
	if err != nil {
		openlog.Error("new request failed.")
		return
	}

	ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
		"user": "peter",
	})
	resp, err := core.NewRestInvoker().ContextDo(ctx, req)
	if err != nil {
		openlog.Error("do request failed.")
		return
	}
	defer resp.Body.Close()
	openlog.Info("REST Server sayhello[GET]: " + string(httputil.ReadBody(resp)))
}
