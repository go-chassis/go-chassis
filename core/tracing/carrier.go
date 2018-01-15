package tracing

import (
	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
)

const (
	prefixTracerState  = "x-b3-" // we default to interop with non-opentracing zipkin tracers
	zipkinTraceID      = prefixTracerState + "traceid"
	zipkinSpanID       = prefixTracerState + "spanid"
	zipkinParentSpanID = prefixTracerState + "parentspanid"
	zipkinSampled      = prefixTracerState + "sampled"
	zipkinFlags        = prefixTracerState + "flags"
)

// RestClientHeaderWriter rest client header writer
type RestClientHeaderWriter rest.Request

// Set to set the header while call API
func (r RestClientHeaderWriter) Set(key, val string) {
	restReq := rest.Request(r)
	restReq.SetHeader(key, val)
}

// FasthttpHeaderCarrier fast http heaer carrier
type FasthttpHeaderCarrier struct {
	*fasthttp.RequestHeader
}

// ForeachKey to check the headers of zipkin
func (f *FasthttpHeaderCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range f.header2TextMap() {
		if err := handler(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (f *FasthttpHeaderCarrier) header2TextMap() map[string]string {
	m := make(map[string]string)
	for _, s := range []string{zipkinTraceID, zipkinSpanID, zipkinParentSpanID, zipkinSampled, zipkinFlags} {
		if v := f.Peek(s); len(v) != 0 {
			m[s] = string(v)
		}
	}
	return m
}
