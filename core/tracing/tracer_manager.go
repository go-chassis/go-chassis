package tracing

import (
	"fmt"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/schema"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/util/iputil"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

// TracerMap tracer map
// key: caller name
// val: tracer
var TracerMap map[string]opentracing.Tracer

// GetTracer get tracer
func GetTracer(caller string) opentracing.Tracer {
	if tracer, ok := TracerMap[caller]; ok {
		return tracer
	}
	return TracerMap[common.DefaultKey]
}

func init() {
	TracerMap = make(map[string]opentracing.Tracer)
}

// Init initialize the tracer
func Init() error {
	lager.Logger.Warn("Tracing enabled. Start to init tracer map.", nil)
	collector, err := NewCollector(config.GlobalDefinition.Tracing.CollectorType, config.GlobalDefinition.Tracing.CollectorTarget)
	if err != nil {
		lager.Logger.Error(err.Error(), nil)
		return fmt.Errorf("unable to create tracing collector: %+v", err)
	}

	microserviceNames := schema.GetMicroserviceNames()
	// key: caller name, val: recorder
	recorderMap := make(map[string]zipkin.SpanRecorder, len(microserviceNames)+1)

	// set default recorder
	defaultCaller := common.DefaultKey
	defaultRecorder := zipkin.NewRecorder(collector, false, "0.0.0.0:0", iputil.GetHostName())
	recorderMap[defaultCaller] = defaultRecorder

	// set recorder map
	for _, msName := range microserviceNames {
		caller := msName + ":" + iputil.GetHostName()
		r := zipkin.NewRecorder(collector, false, "0.0.0.0:0", caller)
		recorderMap[caller] = r
	}

	// set tracer map
	for caller, recorder := range recorderMap {
		// TODO more tracer configuration
		tracer, err := zipkin.NewTracer(
			recorder,
			zipkin.ClientServerSameSpan(true),
			zipkin.TraceID128Bit(true),
		)
		if err != nil {
			lager.Logger.Error(err.Error(), nil)
			return fmt.Errorf("unable to create global tracer: %+v", err)
		}
		TracerMap[caller] = tracer
	}

	return nil
}
