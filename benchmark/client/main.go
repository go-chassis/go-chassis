package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chassis/go-chassis"
	tc "github.com/go-chassis/go-chassis/benchmark/helpers/config"
	"github.com/go-chassis/go-chassis/benchmark/helpers/datacollect"
	"github.com/go-chassis/go-chassis/benchmark/helpers/helloworld/protobuf"
	"github.com/go-chassis/go-chassis/client/rest"
	"github.com/go-chassis/go-chassis/core"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var resultFile *datacollect.ResultFile

const restProtocol = "rest"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)

	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:7778", nil))
	}()
	chassis.Init()

	resultTime := time.Now().Format("20060102150405")
	resultFilePath := filepath.Join(os.Getenv("CHASSIS_HOME"), "result_"+resultTime+".csv")
	resultFile = &datacollect.ResultFile{Path: resultFilePath}
	resultFile.NewFile()

	var msgSize int64 = 1024
	threadCount := 50
	printCount := 5000
	registryEnable := true
	endPoint := "127.0.0.1:8080"
	microServiceName := "Server"
	testCfg, _ := tc.GetTestDefinition()
	protocolName := testCfg.Protocol
	if protocolName == "" {
		protocolName = "highway"
	}
	if testCfg != nil {
		msgSize = int64(testCfg.MessageSize)
		threadCount = testCfg.ThreadCount
		printCount = testCfg.PrintCount
		registryEnable = testCfg.RegistryEnable
		endPoint = testCfg.EndPoint
		if testCfg.MicroServiceName != "" {
			microServiceName = testCfg.MicroServiceName
		}
	}

	var invoker *core.RPCInvoker
	var restInvoker *core.RestInvoker

	invoker = core.NewRPCInvoker()
	replyOne := &protobuf.HelloReply{}
	msg := ""
	go func() {
		timeStr := fmt.Sprintf("\n[%s] ThreadCount:%v, Message Size:%v, Print each Count:%v, Registry enable:%v", time.Now().Format("2006-01-02 15:04:05"), threadCount, msgSize, protocolName, printCount, registryEnable)
		sysInfos := datacollect.GetAllSysInfos()
		resultInfo := []string{timeStr}
		resultInfo = append(resultInfo, sysInfos...)
		resultInfo = append(resultInfo, "\nTime, Total count, Sucess count, Cycle TPS, Token time(s), Average latency time(ms), Mem(MB),Mem Free(MB),Mem Used(MB),Mem Usage, CPU Used, Network(bytes)(Recv / Sent)")
		err := resultFile.Write(resultInfo)
		if err != nil {
			panic(err)
		}
	}()

	var err error
	if registryEnable {
		if protocolName == restProtocol {
			restInvoker = core.NewRestInvoker()
			url := "cse://" + microServiceName + "/createmessage"
			request := &protobuf.MessageRequest{Str: "a", Count: msgSize}
			reqBody, err := json.Marshal(request)
			if err != nil {
				panic(err)
			}
			req, err := rest.NewRequest(http.MethodPost, url, reqBody)
			if err != nil {
				panic(err)
			}
			resp, err := restInvoker.ContextDo(context.TODO(), req)
			if err != nil {
				panic(err)
			}
			msg = string(resp.ReadBody())
		} else {
			err = invoker.Invoke(nil, microServiceName, "HelloServer", "CreateMessage", &protobuf.MessageRequest{Str: "a", Count: msgSize}, replyOne, core.WithProtocol(protocolName))
			if err != nil {
				panic(err)
			}
			msg = replyOne.GetMessage()
		}
	} else {
		err = invoker.Invoke(nil, microServiceName, "HelloServer", "CreateMessage", &protobuf.MessageRequest{Str: "a", Count: msgSize}, replyOne, core.WithProtocol(protocolName), core.WithEndpoint(endPoint))
		if err != nil {
			panic(err)
		}
	}
	beginTime := time.Now()
	for i := 0; i < threadCount; i++ {
		go Call(beginTime, invoker, restInvoker, nil, protocolName, printCount, registryEnable, endPoint, microServiceName, msg, msgSize)
	}
	done := make(chan bool)
	<-done
}

var count, sucessCount, totalTime int64
var mutex sync.Mutex

func Call(beginTime time.Time, invoker *core.RPCInvoker, restInvoker *core.RestInvoker, ctx context.Context, protocolName string, n int, registryEnabled bool, endPoint string, microServiceName string, msg string, msgSize int64) {
	for {
		replyOne := &protobuf.HelloReply{}
		result := ""
		var err error
		if registryEnabled {
			if protocolName == restProtocol {
				url := "cse://" + microServiceName + "/getmessage"
				req, _ := rest.NewRequest(http.MethodGet, url, nil)
				resp, err := restInvoker.ContextDo(ctx, req)
				if err != nil {
					panic(err)
				}
				result = string(resp.ReadBody())
			} else {
				err = invoker.Invoke(nil, microServiceName, "HelloServer", "CreateMessage", &protobuf.MessageRequest{Str: "a", Count: msgSize}, replyOne, core.WithProtocol(protocolName))
				if err != nil {
					panic(err)
				}
				msg = replyOne.GetMessage()
				result = replyOne.GetMessage()
			}
		} else {
			err = invoker.Invoke(nil, microServiceName, "HelloServer", "CreateMessage", &protobuf.MessageRequest{Str: "a", Count: msgSize}, replyOne, core.WithProtocol(protocolName), core.WithEndpoint(endPoint))
			if err != nil {
				panic(err)
			}
			result = replyOne.GetMessage()
		}
		endTime := time.Now()
		mutex.Lock()
		count++
		if err != nil {
			fmt.Println("Error have occured: ", err.Error())
		} else if result != msg || len(msg) == 0 {
			fmt.Println("Error have occured: ", fmt.Sprintf("Expected Msg: %v, But: %v", msg, result))
		} else {
			sucessCount++
		}
		totalTime = totalTime + endTime.Sub(beginTime).Nanoseconds()
		mutex.Unlock()
		if count%int64(n) == 0 {
			go func() {
				tokenTime := float64(endTime.Sub(beginTime).Nanoseconds()) / 1e6 / 1000.0
				//log.Println(fmt.Sprintf("%d/%f", count, endTime.Sub(beginTime).Seconds()))
				cycleTPS := float64(count) / (tokenTime)
				//log.Println(fmt.Sprintf("%d/%f", count, float64(tokenTime) / 1000))
				averageTime := float64(tokenTime) * 1000 / float64(count)
				now := time.Now().Format("2006-01-02 15:04:05")
				resultPrint := fmt.Sprintf("\n[%v] \nTotal count:%v \nSuccess count:%v \nCycle TPS:%v \nTaken Time(s):%v \nAverage Latency time(ms):%v", now, count, sucessCount, cycleTPS, tokenTime, averageTime)
				sysInfos := datacollect.GetAllSysInfos()
				resultInfoPrint := []string{resultPrint}
				resultInfoPrint = append(resultInfoPrint, sysInfos...)

				// For info on each, see: https://golang.org/pkg/runtime/#MemStats
				err := resultFile.Write([]string{result})
				if err != nil {
					panic(err)
				}
				fmt.Println(strings.Join(resultInfoPrint, "\n"))
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				fmt.Printf("Alloc = %v MB\n", bToMb(m.Alloc))
				fmt.Printf("Sys = %v MB\n", bToMb(m.Sys))
				fmt.Printf("NumGC = %v\n", m.NumGC)
			}()
		}
	}
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
