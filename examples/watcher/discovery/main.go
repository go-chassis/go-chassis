package main

import (
	"encoding/json"
	"fmt"
	scregistry "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/sc-client"
	"log"
	"os"
	"time"
)

func main() {

	registryClient, err := sc.NewClient(
		sc.Options{
			Endpoints: []string{"127.0.0.1:30100"},
		})
	if err != nil {
		fmt.Printf("err[%v]\n", err)
		os.Exit(1)
	}

	disService := &scregistry.MicroService{
		AppId:       "default",
		ServiceName: "dismyserver1",
		Version:     "0.0.1",
		Environment: "",
	}

	disSid, err := registryClient.RegisterService(disService)
	if err != nil {
		fmt.Printf("err[%v]\n", err)
		os.Exit(1)
	}
	fmt.Printf("sid[%v]\n", disSid)

	//需要住一个的是服务注册方和 发现方的 appid要一致，否则有跨域问题
	err = registryClient.WatchMicroService(disSid, printEvent) //通知sc, watch;
	if err != nil {
		log.Panicf("WatchMicroService, err[%v]", err)
	}
	_, _ = registryClient.FindMicroServiceInstances(disSid, "default", "myserver1", "0.0.1") //告诉sc, 关注的provider信息

	for {
		time.Sleep(time.Second)
	}
}

func printEvent(event *sc.MicroServiceInstanceChangedEvent) {
	content, _ := json.Marshal(event)
	fmt.Printf("event[%v]\n", string(content))
}
