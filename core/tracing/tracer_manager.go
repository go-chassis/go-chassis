package tracing

import (
	"errors"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/pkg/runtime"
	"github.com/ServiceComb/go-chassis/pkg/util/iputil"

	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

// Option is used to init tracing module
type Option struct {
	ServiceName         string
	ProtocolEndpointMap map[string]string
	CollectorType       string
	CollectorTarget     string
}

var defaultTracer opentracing.Tracer
var protocolTracerMap map[string]opentracing.Tracer

// GetTracer get tracer
func DefaultTracer() opentracing.Tracer {
	return defaultTracer
}

// ConsumerTracer get tracer for consumer
func ConsumerTracer() opentracing.Tracer {
	return defaultTracer
}

// ProviderTracer get tracer for provider
func ProviderTracer(protocol string) (opentracing.Tracer, error) {
	if t, ok := protocolTracerMap[protocol]; ok {
		return t, nil
	}
	return nil, errors.New("no tracer for protocol: " + protocol)
}

// Init initialize the tracer
func Init(opt *Option) error {
	if opt == nil {
		return errors.New("tracing init option is nil")
	}
	if opt.CollectorType == "" {
		lager.Logger.Info("Collector type empty, use noop tracer")
		return nil
	}

	collector, err := NewCollector(opt.CollectorType, opt.CollectorTarget)
	if err != nil {
		return err
	}

	svcName := opt.ServiceName
	if svcName == "" {
		svcName = runtime.HostName
	}

	defaultRecorder := zipkin.NewRecorder(collector, false, iputil.GetLocalIP(), svcName)
	if t, err := newZipkinTracer(defaultRecorder); err != nil {
		return err
	} else {
		defaultTracer = t
	}
	if len(opt.ProtocolEndpointMap) == 0 {
		lager.Logger.Debug("No protocol endpoint provided")
		return nil
	}
	for proto, ep := range opt.ProtocolEndpointMap {
		recorder := zipkin.NewRecorder(collector, false, ep, svcName)
		t, err := newZipkinTracer(recorder)
		if err != nil {
			return err
		}
		protocolTracerMap[proto] = t
	}
	lager.Logger.Info("Tracing init success", nil)
	return nil
}

func newZipkinTracer(r zipkin.SpanRecorder) (opentracing.Tracer, error) {
	return zipkin.NewTracer(
		r,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
}

func init() {
	defaultTracer = opentracing.NoopTracer{}
	protocolTracerMap = make(map[string]opentracing.Tracer)
}
