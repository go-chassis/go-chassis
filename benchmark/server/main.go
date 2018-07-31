package main

import (
	"fmt"
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/benchmark/helpers/datacollect"
	"github.com/go-chassis/go-chassis/benchmark/helpers/helloworld"
	"github.com/go-chassis/go-chassis/core/server"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var resultFile *datacollect.ResultFile

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:7777", nil))
	}()

	var err error
	chassis.Init()

	resultTime := time.Now().Format("20060102150405")
	resultFilePath := filepath.Join(os.Getenv("CHASSIS_HOME"), "result_"+resultTime+".txt")
	resultFile := &datacollect.ResultFile{Path: resultFilePath}
	resultFile.NewFile()
	go func() {
		for {
			timeStr := fmt.Sprintf("\n[%s]", time.Now().Format("2006-01-02 15:04:05"))
			sysInfos := datacollect.GetAllSysInfos()
			resultInfo := []string{timeStr}
			resultInfo = append(resultInfo, sysInfos...)
			err := resultFile.Write(resultInfo)
			if err != nil {
				panic(err)
			}
			fmt.Println(strings.Join(resultInfo, "\n"))
			time.Sleep(time.Second * 10)
		}
	}()
	chassis.RegisterSchema("rest", &helloworld.RestHelloServer{}, server.WithSchemaID("HelloServer"))
	chassis.RegisterSchema("highway", &helloworld.HelloServer{}, server.WithSchemaID("HelloServer"))

	if err != nil {
		panic(err)
	}

	fmt.Println("server start success")
	chassis.Run()
}
