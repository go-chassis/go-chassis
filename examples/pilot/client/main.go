package main

import (
	"context"
	"github.com/go-chassis/openlog"
	"log"
	"time"

	"github.com/go-chassis/go-chassis/v2"
	_ "github.com/go-chassis/go-chassis/v2/bootstrap"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	"github.com/go-chassis/go-chassis/v2/core"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
)

// if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	//chassis operation
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed." + err.Error())
		return
	}
	restInvoker := core.NewRestInvoker()

	// use the configured chain
	for {
		callRest(restInvoker, 10)
		<-time.After(time.Second)
	}
}

func callRest(invoker *core.RestInvoker, i int) {
	url := "http://istioserver/sayhello/b"
	if i < 10 {
		url = "http://istioserver/sayhello/a"
	}
	req, _ := rest.NewRequest("GET", url, nil)
	//use the invoker like http client.
	resp1, err := invoker.ContextDo(context.TODO(), req)
	if err != nil {
		log.Println(err)
	}
	log.Println(i, "REST SayHello ------------------------------ ", resp1.StatusCode, string(httputil.ReadBody(resp1)))

	resp1.Body.Close()
}
