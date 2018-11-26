package bench

import (
	"context"
	"github.com/go-chassis/go-chassis"
	"github.com/go-chassis/go-chassis/core"
	"github.com/go-mesh/openlogging"
	"log"
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"os"
	"strconv"
	"time"
	"fmt"
)

func ReadBody() []byte {
	if Configs.BodyPath == "" {
		return nil
	}
	f, err := os.Open(Configs.BodyPath)
	if err != nil {
		panic(err)
	}
	b := make([]byte, 0)
	n, err := f.Read(b)
	if err != nil {
		panic(err)
	}
	openlogging.GetLogger().Infof("message size: %d", n)
	return b
}
func Benchmark() {
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:7778", nil))
	}()
	var err error
	concurrency, err := strconv.Atoi(Configs.Concurrency)
	if err != nil {
		panic(err)
	}
	if concurrency == 0 {
		panic("-c must be bigger than 0")
	}
	duration, err := time.ParseDuration(Configs.Duration)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(Configs.Target)
	if err != nil {
		panic(err)
	}
	chassis.Init()
	go chassis.Run()
	prepareMetrics()
	cancels := make([]context.CancelFunc, 0)
	switch u.Scheme {
	case "http":
		var invoker = core.NewRestInvoker()
		ctx, cancel := context.WithCancel(context.Background())
		b := ReadBody()
		for i := 0; i < concurrency; i++ {
			openlogging.GetLogger().Info("launched one http benchmark thread")
			cancels = append(cancels, cancel)
			go callHTTP(ctx, invoker, Configs.Method, u.String(), b)
		}
	default:
		panic("not supported " + u.Scheme)

	}
	//wait duration
	t := time.Tick(duration)
	<-t
	//stop all go routines
	for _, c := range cancels {
		c()
	}
	Report(duration)

}
func Report(d time.Duration) {
	fmt.Println("TPS: ", totalRequest/d.Seconds())
	fmt.Println("Total request: ", totalRequest)
	fmt.Println("Err request: ", errRequest)
	t := latency.Snapshot()
	ps := t.Percentiles([]float64{0.05, 0.25, 0.5, 0.75, 0.90, 0.99})
	meanTime := t.Mean() / float64(time.Millisecond)
	fmt.Println("latency mean: ", meanTime)
	fmt.Println("latency p05: ", ps[0]/float64(time.Millisecond))
	fmt.Println("latency p25: ", ps[1]/float64(time.Millisecond))
	fmt.Println("latency p50: ", ps[2]/float64(time.Millisecond))
	fmt.Println("latency p75: ", ps[3]/float64(time.Millisecond))
	fmt.Println("latency p90: ", ps[4]/float64(time.Millisecond))
	fmt.Println("latency p99: ", ps[5]/float64(time.Millisecond))
}
