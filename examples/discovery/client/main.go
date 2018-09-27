package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/go-chassis/go-chassis"
	_ "github.com/go-chassis/go-chassis/bootstrap"
	"github.com/go-chassis/go-chassis/client/rest"
	_ "github.com/go-chassis/go-chassis/config-center"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/examples/schemas/employ"
	"github.com/go-chassis/go-chassis/examples/schemas/helloworld"
	"github.com/go-chassis/go-chassis/pkg/util/httputil"
)

var wg sync.WaitGroup

//if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/discovery/client/
func main() {
	//chassis operation
	if err := chassis.Init(); err != nil {
		lager.Logger.Error("Init failed.")
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
	ctx := context.WithValue(context.Background(), common.ContextHeaderKey{}, map[string]string{
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
		lager.Logger.Errorf("Invoke failed")
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
		lager.Logger.Errorf("Invoke failed")
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
		lager.Logger.Errorf("Invoke failed")
	}
	log.Println("EditEmploy ------------------------------", replyTwo)

	err = invoker.Invoke(ctx, "Server", "EmployServer", "GetEmploys", &employ.EmployRequest{
		Name:       "One",
		Employ:     nil,
		EmployList: replyTwo.EmployList,
	}, replyTwo)
	if err != nil {
		lager.Logger.Errorf("Invoke failed")
	}
	log.Println("GetEmploys ------------------------------", replyTwo)

	err = invoker.Invoke(ctx, "Server", "EmployServer", "DeleteEmploys", &employ.EmployRequest{
		Name:       "Two",
		Employ:     nil,
		EmployList: replyTwo.EmployList,
	}, replyTwo)
	if err != nil {
		lager.Logger.Errorf("Invoke failed")
	}
	log.Println("DeleteEmploys ------------------------------", replyTwo)

}

func callRest(invoker *core.RestInvoker) {
	defer wg.Done()
	req, _ := rest.NewRequest("GET", "cse://Server/sayhello/myidtest", nil)

	//use the invoker like http client.
	resp1, err := invoker.ContextDo(context.TODO(), req)
	if err != nil {
		lager.Logger.Errorf("call request fail (%s) (%d) ", string(httputil.ReadBody(resp1)), resp1.StatusCode)
		return
	}
	log.Printf("Rest Server sayhello[Get] %s", string(httputil.ReadBody(resp1)))
	log.Printf("Cookie from LB %s", string(httputil.GetRespCookie(resp1, common.LBSessionID)))
	httputil.SetCookie(req, common.LBSessionID, string(httputil.GetRespCookie(resp1, common.LBSessionID)))
	req, _ = rest.NewRequest(http.MethodPost, "cse://Server/sayhi", []byte(`{"name": "peter wang and me"}`))
	req.Header.Set("Content-Type", "application/json")
	httputil.SetCookie(req, common.LBSessionID, string(httputil.GetRespCookie(resp1, common.LBSessionID)))
	resp1, err = invoker.ContextDo(context.TODO(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Rest Server sayhi[POST] %s %s", httputil.GetRespCookie(resp1, common.LBSessionID), httputil.ReadBody(resp1))

	req, _ = rest.NewRequest(http.MethodGet, "cse://Server/sayerror", []byte(""))
	httputil.SetCookie(req, common.LBSessionID, string(httputil.GetRespCookie(resp1, common.LBSessionID)))
	resp1, err = invoker.ContextDo(context.TODO(), req)
	if err != nil {
		log.Printf("%s", string(httputil.GetRespCookie(resp1, common.LBSessionID)))
		log.Println(err)
		return
	}
	log.Printf("Rest Server sayerror[GET] %s ", string(httputil.GetRespCookie(resp1, common.LBSessionID)))

	defer resp1.Body.Close()
}
