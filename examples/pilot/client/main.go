package main

import (
	"context"
	"log"
	"sync"

	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/client/rest"
	_ "github.com/ServiceComb/go-chassis/config-center"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
	"time"
)

var wg sync.WaitGroup

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	//chassis operation
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}
	n := 25
	wg.Add(n * 2)

	rpcInvoker := core.NewRPCInvoker()
	restInvoker := core.NewRestInvoker()

	// use the configured chain
	for m := 0; m < n; m++ {
		callRest(restInvoker, m)
	}

	<-time.After(60 * time.Second)
	log.Println("circuit breaker recover")

	for m := 0; m < n; m++ {
		call(rpcInvoker, m)
	}

	wg.Wait()
}

func call(invoker *core.RPCInvoker, i int) {
	defer wg.Done()
	replyOne := &helloworld.HelloReply{}
	//replyTwo := &employ.EmployResponse{}
	// create context with attachments
	ctx := context.WithValue(context.Background(), common.ContextValueKey{}, map[string]string{
		"X-User": "tianxiaoliang",
	})
	h := &helloworld.HelloRequest{Name: "b"}
	if i < 10 {
		h.Name = "a"
	}
	err := invoker.Invoke(ctx, "Server", "HelloServer", "SayHello", h, replyOne)
	if err != nil {
		log.Println(err)
	}
	log.Println(i, "RPC SayHello ------------------------------ ", replyOne)
	//err = invoker.Invoke(ctx, "Server", "EmployServer", "AddEmploy", &employ.EmployRequest{
	//	Employ: &employ.EmployStruct{
	//		Name:  "One",
	//		Phone: "15989351111",
	//	},
	//	EmployList: nil,
	//}, replyTwo)
	//if err != nil {
	//	lager.Logger.Errorf(err, "Invoke failed")
	//}
	//log.Println("AddEmploy ------------------------------", replyTwo)
	//
	//err = invoker.Invoke(ctx, "Server", "EmployServer", "AddEmploy", &employ.EmployRequest{
	//	Employ: &employ.EmployStruct{
	//		Name:  "Two",
	//		Phone: "15989352222",
	//	},
	//	EmployList: replyTwo.EmployList,
	//}, replyTwo)
	//if err != nil {
	//	lager.Logger.Errorf(err, "Invoke failed")
	//}
	//log.Println("AddEmploy ------------------------------", replyTwo)
	//
	//err = invoker.Invoke(ctx, "Server", "EmployServer", "EditEmploy", &employ.EmployRequest{
	//	Name: "Two",
	//	Employ: &employ.EmployStruct{
	//		Name:  "Two",
	//		Phone: "15989353333",
	//	},
	//	EmployList: replyTwo.EmployList,
	//}, replyTwo)
	//if err != nil {
	//	lager.Logger.Errorf(err, "Invoke failed")
	//}
	//log.Println("EditEmploy ------------------------------", replyTwo)
	//
	//err = invoker.Invoke(ctx, "Server", "EmployServer", "GetEmploys", &employ.EmployRequest{
	//	Name:       "One",
	//	Employ:     nil,
	//	EmployList: replyTwo.EmployList,
	//}, replyTwo)
	//if err != nil {
	//	lager.Logger.Errorf(err, "Invoke failed")
	//}
	//log.Println("GetEmploys ------------------------------", replyTwo)
	//
	//err = invoker.Invoke(ctx, "Server", "EmployServer", "DeleteEmploys", &employ.EmployRequest{
	//	Name:       "Two",
	//	Employ:     nil,
	//	EmployList: replyTwo.EmployList,
	//}, replyTwo)
	//if err != nil {
	//	lager.Logger.Errorf(err, "Invoke failed")
	//}
	//log.Println("DeleteEmploys ------------------------------", replyTwo)

}

func callRest(invoker *core.RestInvoker, i int) {
	defer wg.Done()
	url := "cse://Server/sayhello/b"
	if i < 10 {
		url = "cse://Server/sayhello/a"
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
