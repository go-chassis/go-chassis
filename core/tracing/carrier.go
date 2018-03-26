package tracing

import (
	"github.com/ServiceComb/go-chassis/client/rest"
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

// HeaderCarrier http heaer carrier
type HeaderCarrier struct {
	Header map[string]string
}

// ForeachKey to check the headers of zipkin
func (f *HeaderCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, v := range f.header2TextMap() {
		if err := handler(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (f *HeaderCarrier) header2TextMap() map[string]string {
	m := make(map[string]string)
	for _, s := range []string{zipkinTraceID, zipkinSpanID, zipkinParentSpanID, zipkinSampled, zipkinFlags} {
		if _, ok := f.Header[s]; ok {
			m[s] = f.Header[s]
		}
	}
	return m
}
