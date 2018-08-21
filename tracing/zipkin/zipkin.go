package zipkin

import (
	"fmt"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/tracing"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/opentracing/opentracing-go"
	"github.com/openzipkin/zipkin-go-opentracing"
	"strconv"
	"time"
)

//const for default values
const (
	DefaultURI           = "http://127.0.0.1:9411/api/v1/spans"
	DefaultBatchSize     = 10000
	DefaultBatchInterval = time.Second * 10
	DefaultCollector     = "http"
)

// NewTracer returns zipkin tracer
func NewTracer(options map[string]string) (opentracing.Tracer, error) {
	uri := options["URI"]
	if uri == "" {
		uri = DefaultURI
	}
	var batchSize = DefaultBatchSize
	bs := options["batchSize"]
	if bs != "" {
		var err error
		batchSize, err = strconv.Atoi(bs)
		if err != nil {
			return nil, fmt.Errorf("can not convert [%s] to batch size", bs)
		}
	}
	var batchInterval = DefaultBatchInterval
	bi := options["batchInterval"]
	if bi != "" {
		var err error
		batchInterval, err = time.ParseDuration(bi)
		if err != nil {
			return nil, fmt.Errorf("can not convert [%s] to batch interval", bi)
		}
	}
	var collectorOption string
	var collector zipkintracer.Collector
	if options["collector"] == "" {
		collectorOption = DefaultCollector
	}
	lager.Logger.Infof("New Zipkin tracer with options %s,%s,%s", uri, batchSize, batchInterval)
	if collectorOption == DefaultCollector {
		var err error
		collector, err = zipkintracer.NewHTTPCollector(uri, zipkintracer.HTTPBatchSize(batchSize), zipkintracer.HTTPBatchInterval(batchInterval))
		if err != nil {
			lager.Logger.Error(err.Error(), nil)
			return nil, fmt.Errorf("unable to create zipkin collector: %+v", err)
		}
	} else if collectorOption == "namedPipe" {
		var err error
		collector, err = newNamedPipeCollector(uri)
		if err != nil {
			lager.Logger.Error(err.Error(), nil)
			return nil, fmt.Errorf("unable to create zipkin collector: %+v", err)
		}
	} else {
		return nil, fmt.Errorf("unable to create zipkin collector: %s", collectorOption)
	}

	// set default recorder
	defaultRecorder := zipkintracer.NewRecorder(collector, false, "0.0.0.0:0", runtime.HostName)

	// set tracer map
	tracer, err := zipkintracer.NewTracer(
		defaultRecorder,
		zipkintracer.ClientServerSameSpan(true),
		zipkintracer.TraceID128Bit(true),
	)
	if err != nil {
		lager.Logger.Error(err.Error(), nil)
		return nil, fmt.Errorf("unable to create zipkin tracer: %+v", err)
	}
	return tracer, nil
}

func init() {
	tracing.InstallTracer("zipkin", NewTracer)
}
