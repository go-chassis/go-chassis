package main

import (
	"context"
	"fmt"
	"github.com/go-chassis/openlog"

	"github.com/go-chassis/go-chassis/v2"
	_ "github.com/go-chassis/go-chassis/v2/bootstrap"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
	"time"
)

// if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/
func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	// run is only to enable metric exporter
	go chassis.Run()
	ctx := context.WithValue(context.TODO(), common.ContextHeaderKey{}, map[string]string{
		"user": "peter",
	})
	invoker := core.NewRestInvoker()

	for i := 0; i < 500; i++ {
		req, err := rest.NewRequest("GET", "http://ErrServer/lock", nil)
		if err != nil {
			openlog.Error("new request failed.")
			return
		}
		resp, err := invoker.ContextDo(ctx, req)
		if err != nil {
			openlog.Error(fmt.Sprintf("deadlock request failed. %s", err.Error()))
		} else {
			openlog.Info("REST Server [GET]: " + string(httputil.ReadBody(resp)))
		}

	}

	for {
		req, err := rest.NewRequest("GET", "http://ErrServer/sayhimessage", nil)
		if err != nil {
			openlog.Error("new request failed.")
			return
		}
		resp, err := invoker.ContextDo(ctx, req)
		if err != nil {
			openlog.Error(fmt.Sprintf("normal request failed. %s", err.Error()))
		} else {
			openlog.Info("REST Server [GET]: " + string(httputil.ReadBody(resp)))
		}
		time.Sleep(30 * time.Second)
	}
}
