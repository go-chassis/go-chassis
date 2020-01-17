package main

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/pkg/scclient"
	"github.com/go-chassis/go-chassis/pkg/scclient/proto"
	"os"
	"time"
)

func main()  {
	registryClient := &client.RegistryClient{}

	err := registryClient.Initialize(
		client.Options{
			Addrs: []string{"127.0.0.1:30100"},
		})
	if err != nil {
		fmt.Printf("err[%v]\n", err)
		os.Exit(1)
	}

	/*
	"myapp1.1", "myserver1", "0.0.1", "")
	*/
	service := &proto.MicroService{
		AppId: "default",
		ServiceName:"myserver1",
		Version: "0.0.1",
		Environment: "",
	}

	sid ,err := registryClient.RegisterService(service)
	if err != nil {
		fmt.Printf("err[%v]\n", err)
		os.Exit(1)
	}
	fmt.Printf("sid[%v]\n", sid)


	instance := proto.MicroServiceInstance{
		ServiceId: sid,
		HostName: "insdfsdsdfsdff2233",
		Status: common.DefaultStatus,
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
	for ; ;  {
		count++
		if count == 10 {
			break
		}
		time.Sleep(time.Second)
	}
	registryClient.UnregisterMicroServiceInstance(sid, iid)
	fmt.Printf("sid[%v], iid[%v]\n", sid, iid)
}


