package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/ServiceComb/go-chassis"
	_ "github.com/ServiceComb/go-chassis/bootstrap"
	"github.com/ServiceComb/go-chassis/client/rest"
	_ "github.com/ServiceComb/go-chassis/config-center"
	"github.com/ServiceComb/go-chassis/core"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/examples/schemas/employ"
	"github.com/ServiceComb/go-chassis/examples/schemas/helloworld"
)

var wg sync.WaitGroup

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/path/to/conf/folder
func main() {
	//chassis operation
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.", err)
		return
	}

	// use the configured chain
	invoker := core.NewRPCInvoker()
	call(invoker)
	n := 10
	wg.Add(n)
	restInvoker := core.NewRestInvoker()
	for m := 0; m < n; m++ {
		go callRest(restInvoker)
	}
	wg.Wait()
}

func call(invoker *core.RPCInvoker) {
	replyOne := &helloworld.HelloReply{}
	replyTwo := &employ.EmployResponse{}
	// create context with attachments
	ctx := context.WithValue(context.Background(), common.ContextValueKey{}, map[string]string{
		"X-User": "tianxiaoliang",
	})
	err := invoker.Invoke(ctx, "Server", "HelloServer", "SayHello", &helloworld.HelloRequest{Name: "Peter"}, replyOne)
	if err != nil {
		log.Println(err)
	}
	log.Println("SayHello ------------------------------ ", replyOne)
	err = invoker.Invoke(ctx, "Server", "EmployServer", "AddEmploy", &employ.EmployRequest{
		Employ: &employ.EmployStruct{
			Name:  "One",
			Phone: "15989351111",
		},
		EmployList: nil,
	}, replyTwo)
	if err != nil {
		lager.Logger.Errorf(err, "Invoke failed")
	}
	log.Println("AddEmploy ------------------------------", replyTwo)

	err = invoker.Invoke(ctx, "Server", "EmployServer", "AddEmploy", &employ.EmployRequest{
		Employ: &employ.EmployStruct{
			Name:  "Two",
			Phone: "15989352222",
		},
		EmployList: replyTwo.EmployList,
	}, replyTwo)
	if err != nil {
		lager.Logger.Errorf(err, "Invoke failed")
	}
	log.Println("AddEmploy ------------------------------", replyTwo)

	err = invoker.Invoke(ctx, "Server", "EmployServer", "EditEmploy", &employ.EmployRequest{
		Name: "Two",
		Employ: &employ.EmployStruct{
			Name:  "Two",
			Phone: "15989353333",
		},
		EmployList: replyTwo.EmployList,
	}, replyTwo)
	if err != nil {
		lager.Logger.Errorf(err, "Invoke failed")
	}
	log.Println("EditEmploy ------------------------------", replyTwo)

	err = invoker.Invoke(ctx, "Server", "EmployServer", "GetEmploys", &employ.EmployRequest{
		Name:       "One",
		Employ:     nil,
		EmployList: replyTwo.EmployList,
	}, replyTwo)
	if err != nil {
		lager.Logger.Errorf(err, "Invoke failed")
	}
	log.Println("GetEmploys ------------------------------", replyTwo)

	err = invoker.Invoke(ctx, "Server", "EmployServer", "DeleteEmploys", &employ.EmployRequest{
		Name:       "Two",
		Employ:     nil,
		EmployList: replyTwo.EmployList,
	}, replyTwo)
	if err != nil {
		lager.Logger.Errorf(err, "Invoke failed")
	}
	log.Println("DeleteEmploys ------------------------------", replyTwo)

}

func callRest(invoker *core.RestInvoker) {
	defer wg.Done()
	req, _ := rest.NewRequest("GET", "cse://Server/sayhello/myidtest")

	//use the invoker like http client.
	resp1, err := invoker.ContextDo(context.TODO(), req)
	if err != nil {
		lager.Logger.Errorf(err, "call request fail (%s) (%d) ", string(resp1.ReadBody()), resp1.GetStatusCode())
		return
	}
	log.Printf("Rest Server sayhello[Get] %s", string(resp1.ReadBody()))
	log.Printf("Cookie from LB %s", string(resp1.GetCookie(common.LBSessionID)))
	req.SetCookie(common.LBSessionID, string(resp1.GetCookie(common.LBSessionID)))
	req, _ = rest.NewRequest(http.MethodPost, "cse://Server/sayhi", []byte(`{"name": "peter wang and me"}`))
	req.SetHeader("Content-Type", "application/json")
	req.SetCookie(common.LBSessionID, string(resp1.GetCookie(common.LBSessionID)))
	resp1, err = invoker.ContextDo(context.TODO(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Rest Server sayhi[POST] %s", string(resp1.ReadBody()))

	req, _ = rest.NewRequest(http.MethodGet, "cse://Server/sayerror", []byte(""))
	req.SetCookie(common.LBSessionID, string(resp1.GetCookie(common.LBSessionID)))
	resp1, err = invoker.ContextDo(context.TODO(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Rest Server sayerror[GET] %s ", string(resp1.ReadBody()))

	req.Close()
	resp1.Close()
}
