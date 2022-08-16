package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chassis/go-chassis/v2"
	_ "github.com/go-chassis/go-chassis/v2/bootstrap"
	"github.com/go-chassis/go-chassis/v2/client/rest"
	_ "github.com/go-chassis/go-chassis/v2/configserver"
	"github.com/go-chassis/go-chassis/v2/core"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/pkg/util/httputil"
	"github.com/go-chassis/openlog"
)

var wg sync.WaitGroup

// if you use go run main.go instead of binary run, plz export CHASSIS_HOME=/{path}/{to}/discovery/client/
func main() {
	//chassis operation
	if err := chassis.Init(); err != nil {
		openlog.Error("Init failed.")
		return
	}

	n := 10
	wg.Add(n)
	restInvoker := core.NewRestInvoker()
	for m := 0; m < n; m++ {
		go callRest(restInvoker)
	}
	wg.Wait()
}

func callRest(invoker *core.RestInvoker) {
	defer wg.Done()
	req, _ := rest.NewRequest("GET", "http://Server/sayhello/myidtest", nil)

	//use the invoker like http client.
	resp1, err := invoker.ContextDo(context.TODO(), req)
	if err != nil {
		openlog.Error(fmt.Sprintf("call request fail (%s) (%d) ", string(httputil.ReadBody(resp1)), resp1.StatusCode))
		return
	}
	log.Printf("Rest Server sayhello[Get] %s", string(httputil.ReadBody(resp1)))
	log.Printf("Cookie from LB %s", string(httputil.GetRespCookie(resp1, common.LBSessionID)))
	httputil.SetCookie(req, common.LBSessionID, string(httputil.GetRespCookie(resp1, common.LBSessionID)))
	req, _ = rest.NewRequest(http.MethodPost, "http://Server/sayhi", []byte(`{"name": "peter wang and me"}`))
	req.Header.Set("Content-Type", "application/json")
	httputil.SetCookie(req, common.LBSessionID, string(httputil.GetRespCookie(resp1, common.LBSessionID)))
	resp1, err = invoker.ContextDo(context.TODO(), req)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Rest Server sayhi[POST] %s %s", httputil.GetRespCookie(resp1, common.LBSessionID), httputil.ReadBody(resp1))

	req, _ = rest.NewRequest(http.MethodGet, "http://Server/sayerror", []byte(""))
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
