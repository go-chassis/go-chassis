package main

import (
	"context"
	"encoding/json"
	"time"

	_ "github.com/go-chassis/go-chassis-config/servicecomb"

	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/rest/client/

// Implement grayscale publishing of the application,version  A is you old service ,version B is you
// new service.you want to small request to access you new service to test new service of new function

func main() {
	//Init framework
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed." + err.Error())
		return
	}

	for i := 0; i < 10; i++ {
		req, err := rest.NewRequest("POST", "http://paymentService/equal", nil)
		if err != nil {
			lager.Logger.Error("new request failed.")
			return
		}
		parm := struct {
			Num  int
			Nums []int
		}{
			Num:  10,
			Nums: []int{2, 5, 3},
		}

		parmByte, _ := json.Marshal(parm)
		httputil.SetBody(req, parmByte)

		resp, err := core.NewRestInvoker().ContextDo(context.Background(), req)
		if err != nil {
			lager.Logger.Error("do request failed.")
			return
		}
		defer resp.Body.Close()
		lager.Logger.Info("paymentService equal num [POST]: " + string(httputil.ReadBody(resp)))

		time.Sleep(1 * time.Second)
	}
}
