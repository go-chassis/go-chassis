package main

import (
	"context"
	"log"

	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/client/rest"
	_ "github.com/ServiceComb/go-chassis/config-center"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/lager"
	"time"
)

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	//chassis operation
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
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
	url := "cse://istioserver/sayhello/b"
	if i < 10 {
		url = "cse://istioserver/sayhello/a"
	}
	req, _ := rest.NewRequest("GET", url)
	//use the invoker like http client.
	resp1, err := invoker.ContextDo(context.TODO(), req)
	if err != nil {
		log.Println(err)
		//lager.Logger.Errorf(err, "call request fail (%s) (%d) ", string(resp1.ReadBody()), resp1.GetStatusCode())
		//return
	}
	log.Println(i, "REST SayHello ------------------------------ ", resp1.GetStatusCode(), string(resp1.ReadBody()))

	//req, _ = rest.NewRequest(http.MethodPost, "cse://Server/sayhi", []byte(`{"name": "peter wang and me"}`))
	//req.SetHeader("Content-Type", "application/json")
	//resp1, err = invoker.ContextDo(context.TODO(), req)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Printf("Rest Server sayhi[POST] %s", string(resp1.ReadBody()))
	//
	//req, _ = rest.NewRequest(http.MethodGet, "cse://Server/sayerror", []byte(""))
	//resp1, err = invoker.ContextDo(context.TODO(), req)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//log.Printf("Rest Server sayerror[GET] %s ", string(resp1.ReadBody()))

	req.Close()
	resp1.Close()
}
