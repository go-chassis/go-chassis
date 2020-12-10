package main

import (
	"fmt"
	scregistry "github.com/go-chassis/cari/discovery"
	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/sc-client"
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

	service := &scregistry.MicroService{
		AppId:       "default",
		ServiceName: "myserver1",
		Version:     "0.0.1",
		Environment: "",
	}

	sid, err := registryClient.RegisterService(service)
	if err != nil {
		fmt.Printf("err[%v]\n", err)
		os.Exit(1)
	}
	fmt.Printf("sid[%v]\n", sid)

	instance := scregistry.MicroServiceInstance{
		ServiceId: sid,
		HostName:  "insdfsdsdfsdff2233",
		Status:    common.DefaultStatus,
		Endpoints: []string{"localhost:808"},
		Properties: map[string]string{
			"Name": "12",
		},
	}

	iid, err := registryClient.RegisterMicroServiceInstance(&instance)
	if err != nil {
		fmt.Printf("RegisterMicroServiceInstance, err[%v]\n", err)
		os.Exit(1)
	}

	count := 0
	for {
		count++
		if count == 10 {
			break
		}
		time.Sleep(time.Second)
	}
	registryClient.UnregisterMicroServiceInstance(sid, iid)
	fmt.Printf("sid[%v], iid[%v]\n", sid, iid)
}
