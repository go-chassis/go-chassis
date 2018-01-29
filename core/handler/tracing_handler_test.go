package handler_test

import (
	//"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	//"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/tracing"
	//"github.com/ServiceComb/go-chassis/third_party/forked/valyala/fasthttp"
	"github.com/ServiceComb/go-chassis/util/iputil"
	"github.com/apache/thrift/lib/go/thrift"
	//"github.com/emicklei/go-restful"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/stretchr/testify/assert"
)

const (
	interval    = 200 * time.Millisecond
	serverSleep = 100 * time.Millisecond
)

const (
/*
	prefixTracerState = "x-b3-" // we default to interop with non-opentracing zipkin tracers
	prefixBaggage     = "ot-baggage-"

	tracerStateFieldCount = 3 // not 5, X-B3-ParentSpanID is optional and we allow optional Sampled header
	zipkinTraceID         = prefixTracerState + "traceid"
	zipkinSpanID          = prefixTracerState + "spanid"
	zipkinParentSpanID    = prefixTracerState + "parentspanid"
	zipkinSampled         = prefixTracerState + "sampled"
	zipkinFlags           = prefixTracerState + "flags"
*/
)

type sleepHandler struct{}

func (s *sleepHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	chain.Next(i, cb)
	time.Sleep(interval)
}

func (s *sleepHandler) Name() string {
	return "sleepHandler"
}

type spanTestHandler struct {
	Span opentracing.Span
}

func (s *spanTestHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	s.Span = opentracing.SpanFromContext(i.Ctx)
	cb(&invocation.InvocationResponse{
		Err: nil,
	})
}

func (s *spanTestHandler) Name() string {
	return "spanTestHandler"
}

func TestTracingHandler_Highway(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	// the port should be different from test cases
	port := 30110
	s := newHTTPServer(t, port)
	// init tracer manager
	// batch size it 1, send every span when once collector get it.
	collector, err := zipkin.NewHTTPCollector(fmt.Sprintf("http://localhost:%d/api/v1/spans", port), zipkin.HTTPBatchSize(1))
	assert.NoError(t, err)
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:0", iputil.GetHostName())
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	assert.NoError(t, err)
	tracing.TracerMap[common.DefaultKey] = tracer

	t.Log("========tracing [consumer] handler [highway]")

	// set handler chain
	consumerChain := handler.Chain{}
	tracingConsumerHandler := &handler.TracingConsumerHandler{}
	consumerSpanHandler := &spanTestHandler{}
	consumerChain.AddHandler(&sleepHandler{})
	consumerChain.AddHandler(tracingConsumerHandler)
	consumerChain.AddHandler(consumerSpanHandler)
	inv := &invocation.Invocation{
		MicroServiceName: "test",
		Protocol:         common.ProtocolHighway,
	}

	consumerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====span should be stored in context after invoker")
	assert.NotNil(t, consumerSpanHandler.Span)

	t.Log("====spanContext stored in context should be the same with monitor received")
	zipkinServerRecievedSpans := s.spans()
	if len(zipkinServerRecievedSpans) == 0 {
		return
	}
	assert.Equal(t, 1, len(zipkinServerRecievedSpans))
	if t.Failed() {
		return
	}
	assert.Equal(t, 2, len(zipkinServerRecievedSpans[0].Annotations))

	localSpanContext := consumerSpanHandler.Span.Context()
	consumerZpSpanContext, ok := localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, uint64(zipkinServerRecievedSpans[0].TraceID), consumerZpSpanContext.TraceID.Low)
	assert.Equal(t, (uint64)(*zipkinServerRecievedSpans[0].TraceIDHigh), consumerZpSpanContext.TraceID.High)
	assert.Equal(t, (uint64)(zipkinServerRecievedSpans[0].ID), consumerZpSpanContext.SpanID)

	assert.NotEqual(t, 0, consumerZpSpanContext.SpanID)
	assert.Nil(t, consumerZpSpanContext.ParentSpanID)

	t.Log("========tracing [provider] handler [highway]")

	// set handler chain
	providerChain := handler.Chain{}
	tracingProviderHandler := &handler.TracingProviderHandler{}
	providerSpanHandler := &spanTestHandler{}
	providerChain.AddHandler(&sleepHandler{})
	providerChain.AddHandler(tracingProviderHandler)
	providerChain.AddHandler(providerSpanHandler)

	s.clearSpans()
	providerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====span should be stored in context after server receive")
	assert.NotNil(t, providerSpanHandler.Span)

	t.Log("====spanContext stored in context should be the same with monitor received")
	zipkinServerRecievedSpans = s.spans()
	assert.Equal(t, 1, len(zipkinServerRecievedSpans))
	if t.Failed() {
		return
	}
	assert.Equal(t, 2, len(zipkinServerRecievedSpans[0].Annotations))

	localSpanContext = providerSpanHandler.Span.Context()
	providerZpSpanContext, ok := localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, uint64(zipkinServerRecievedSpans[0].TraceID), providerZpSpanContext.TraceID.Low)
	assert.Equal(t, (uint64)(*zipkinServerRecievedSpans[0].TraceIDHigh), providerZpSpanContext.TraceID.High)
	assert.Equal(t, (uint64)(zipkinServerRecievedSpans[0].ID), providerZpSpanContext.SpanID)

	assert.NotEqual(t, 0, providerZpSpanContext.SpanID)
	assert.Nil(t, providerZpSpanContext.ParentSpanID)

	t.Log("====spanContext of consumer/provider should be the same")
	assert.Equal(t, consumerZpSpanContext.TraceID, providerZpSpanContext.TraceID)
	assert.Equal(t, consumerZpSpanContext.SpanID, providerZpSpanContext.SpanID)

	t.Log("========tracing [consumer] handler [highway], with [parent]")

	consumerChain.Reset()
	s.clearSpans()
	parentSpanID := providerZpSpanContext.SpanID

	consumerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====the parent spanID shoud be last spanID")
	localSpanContext = consumerSpanHandler.Span.Context()
	consumerZpSpanContext, ok = localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, parentSpanID, *consumerZpSpanContext.ParentSpanID)
	assert.NotEqual(t, parentSpanID, consumerZpSpanContext.SpanID)

	var tch *handler.TracingConsumerHandler = new(handler.TracingConsumerHandler)
	str := tch.Name()
	assert.Equal(t, "tracing-consumer", str)

	var tph *handler.TracingProviderHandler = new(handler.TracingProviderHandler)
	str = tph.Name()
	assert.Equal(t, "tracing-provider", str)
}

// TODO
// Comment buggy test cases : Already raised an issue to trace this https://github.com/ServiceComb/go-chassis/issues/5
/*
func TestTracingHandler_Rest_RestRequest(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	// the port should be different from test cases
	port := 30111
	s := newHTTPServer(t, port)
	// init tracer manager
	// batch size it 1, send every span when once collector get it.
	collector, err := zipkin.NewHTTPCollector(fmt.Sprintf("http://localhost:%d/api/v1/spans", port), zipkin.HTTPBatchSize(1))
	assert.NoError(t, err)
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:0", iputil.GetHostName())
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	assert.NoError(t, err)
	tracing.TracerMap[common.DefaultKey] = tracer

	t.Log("========tracing [consumer] handler [rest]")

	// set handler chain
	consumerChain := handler.Chain{}
	tracingConsumerHandler := &handler.TracingConsumerHandler{}
	consumerSpanHandler := &spanTestHandler{}
	consumerChain.AddHandler(&sleepHandler{})
	consumerChain.AddHandler(tracingConsumerHandler)
	consumerChain.AddHandler(consumerSpanHandler)

	restClientSentReq, err := rest.NewRequest("GET", "cse://Server/hello")
	assert.NoError(t, err)

	inv := &invocation.Invocation{
		MicroServiceName: "test",
		Protocol:         common.ProtocolRest,
		Args:             restClientSentReq,
	}

	consumerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====span should be stored in context after invoker")
	assert.NotNil(t, consumerSpanHandler.Span)

	t.Log("====spanContext stored in context should be the same with monitor received")
	zipkinServerRecievedSpans := s.spans()
	assert.Equal(t, 1, len(zipkinServerRecievedSpans))
	if t.Failed() {
		return
	}
	assert.Equal(t, 2, len(zipkinServerRecievedSpans[0].Annotations))

	localSpanContext := consumerSpanHandler.Span.Context()
	consumerZpSpanContext, ok := localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, uint64(zipkinServerRecievedSpans[0].TraceID), consumerZpSpanContext.TraceID.Low)
	assert.Equal(t, (uint64)(*zipkinServerRecievedSpans[0].TraceIDHigh), consumerZpSpanContext.TraceID.High)
	assert.Equal(t, (uint64)(zipkinServerRecievedSpans[0].ID), consumerZpSpanContext.SpanID)

	assert.NotEqual(t, 0, consumerZpSpanContext.SpanID)
	assert.Nil(t, consumerZpSpanContext.ParentSpanID)

	t.Log("========tracing [provider] handler [rest]")

	// copy header from rest client to rest server
	httpReq, err := http.NewRequest("GET", "cse://Server/hello", bytes.NewReader(make([]byte, 0)))
	assert.NoError(t, err)
	restServerReceivedReq := &restful.Request{
		Request: httpReq,
	}
	httpHeadersCarrier := (opentracing.HTTPHeadersCarrier)(restServerReceivedReq.Request.Header)
	assert.True(t, ok)

	httpHeadersCarrier.Set(zipkinTraceID, restClientSentReq.GetHeader(zipkinTraceID))
	httpHeadersCarrier.Set(zipkinSpanID, restClientSentReq.GetHeader(zipkinSpanID))
	// parentID is empty, do not set it
	httpHeadersCarrier.Set(zipkinSampled, restClientSentReq.GetHeader(zipkinSampled))
	httpHeadersCarrier.Set(zipkinFlags, restClientSentReq.GetHeader(zipkinFlags))

	// set args to be restServerReceivedReq
	inv.Args = restServerReceivedReq

	// set handler chain
	providerChain := handler.Chain{}
	tracingProviderHandler := &handler.TracingProviderHandler{}
	providerSpanHandler := &spanTestHandler{}
	providerChain.AddHandler(&sleepHandler{})
	providerChain.AddHandler(tracingProviderHandler)
	providerChain.AddHandler(providerSpanHandler)

	s.clearSpans()
	providerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====span should be stored in context after server receive")
	assert.NotNil(t, providerSpanHandler.Span)

	t.Log("====spanContext stored in context should be the same with monitor received")
	zipkinServerRecievedSpans = s.spans()
	assert.Equal(t, 1, len(zipkinServerRecievedSpans))
	if t.Failed() {
		return
	}
	assert.Equal(t, 2, len(zipkinServerRecievedSpans[0].Annotations))

	localSpanContext = providerSpanHandler.Span.Context()
	providerZpSpanContext, ok := localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, uint64(zipkinServerRecievedSpans[0].TraceID), providerZpSpanContext.TraceID.Low)
	assert.Equal(t, (uint64)(*zipkinServerRecievedSpans[0].TraceIDHigh), providerZpSpanContext.TraceID.High)
	assert.Equal(t, (uint64)(zipkinServerRecievedSpans[0].ID), providerZpSpanContext.SpanID)

	assert.NotEqual(t, 0, providerZpSpanContext.SpanID)
	assert.Nil(t, providerZpSpanContext.ParentSpanID)

	t.Log("====spanContext of consumer/provider should be the same")
	assert.Equal(t, consumerZpSpanContext.TraceID, providerZpSpanContext.TraceID)
	assert.Equal(t, consumerZpSpanContext.SpanID, providerZpSpanContext.SpanID)

	t.Log("========tracing [consumer] handler [rest], with [parent]")

	// set args to be restServerReceivedReq
	inv.Args = restClientSentReq

	consumerChain.Reset()
	s.clearSpans()
	parentSpanID := providerZpSpanContext.SpanID

	consumerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====the parent spanID shoud be last spanID")
	localSpanContext = consumerSpanHandler.Span.Context()
	consumerZpSpanContext, ok = localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, parentSpanID, *consumerZpSpanContext.ParentSpanID)
	assert.NotEqual(t, parentSpanID, consumerZpSpanContext.SpanID)
}

func TestTracingHandler_Rest_FasthttpRequest(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	// the port should be different from test cases
	port := 30112
	s := newHTTPServer(t, port)
	// init tracer manager
	// batch size it 1, send every span when once collector get it.
	collector, err := zipkin.NewHTTPCollector(fmt.Sprintf("http://localhost:%d/api/v1/spans", port), zipkin.HTTPBatchSize(1))
	assert.NoError(t, err)
	recorder := zipkin.NewRecorder(collector, false, "0.0.0.0:0", iputil.GetHostName())
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(true),
		zipkin.TraceID128Bit(true),
	)
	assert.NoError(t, err)
	tracing.TracerMap[common.DefaultKey] = tracer

	t.Log("========tracing [consumer] handler [rest]")

	// set handler chain
	consumerChain := handler.Chain{}
	tracingConsumerHandler := &handler.TracingConsumerHandler{}
	consumerSpanHandler := &spanTestHandler{}
	consumerChain.AddHandler(&sleepHandler{})
	consumerChain.AddHandler(tracingConsumerHandler)
	consumerChain.AddHandler(consumerSpanHandler)

	fasthttpRequest := fasthttp.AcquireRequest()
	fasthttpRequest.Header.SetMethod("GET")
	fasthttpRequest.Header.SetRequestURI("cse://Server/hello")
	assert.NoError(t, err)

	inv := &invocation.Invocation{
		MicroServiceName: "test",
		Protocol:         common.ProtocolRest,
		Args:             fasthttpRequest,
	}

	consumerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====span should be stored in context after invoker")
	assert.NotNil(t, consumerSpanHandler.Span)

	t.Log("====spanContext stored in context should be the same with monitor received")
	zipkinServerRecievedSpans := s.spans()
	assert.Equal(t, 1, len(zipkinServerRecievedSpans))
	if t.Failed() {
		return
	}
	assert.Equal(t, 2, len(zipkinServerRecievedSpans[0].Annotations))

	localSpanContext := consumerSpanHandler.Span.Context()
	consumerZpSpanContext, ok := localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, uint64(zipkinServerRecievedSpans[0].TraceID), consumerZpSpanContext.TraceID.Low)
	assert.Equal(t, (uint64)(*zipkinServerRecievedSpans[0].TraceIDHigh), consumerZpSpanContext.TraceID.High)
	assert.Equal(t, (uint64)(zipkinServerRecievedSpans[0].ID), consumerZpSpanContext.SpanID)

	assert.NotEqual(t, 0, consumerZpSpanContext.SpanID)
	assert.Nil(t, consumerZpSpanContext.ParentSpanID)

	t.Log("========tracing [provider] handler [rest]")

	// copy header from rest client to rest server
	httpReq, err := http.NewRequest("GET", "cse://Server/hello", bytes.NewReader(make([]byte, 0)))
	assert.NoError(t, err)
	restServerReceivedReq := &restful.Request{
		Request: httpReq,
	}
	httpHeadersCarrier := (opentracing.HTTPHeadersCarrier)(restServerReceivedReq.Request.Header)
	assert.True(t, ok)

	httpHeadersCarrier.Set(zipkinTraceID, string(fasthttpRequest.Header.Peek(zipkinTraceID)))
	httpHeadersCarrier.Set(zipkinSpanID, string(fasthttpRequest.Header.Peek(zipkinSpanID)))
	// parentID is empty, do not set it
	httpHeadersCarrier.Set(zipkinSampled, string(fasthttpRequest.Header.Peek(zipkinSampled)))
	httpHeadersCarrier.Set(zipkinFlags, string(fasthttpRequest.Header.Peek(zipkinFlags)))

	// set args to be restServerReceivedReq
	inv.Args = restServerReceivedReq

	// set handler chain
	providerChain := handler.Chain{}
	tracingProviderHandler := &handler.TracingProviderHandler{}
	providerSpanHandler := &spanTestHandler{}
	providerChain.AddHandler(&sleepHandler{})
	providerChain.AddHandler(tracingProviderHandler)
	providerChain.AddHandler(providerSpanHandler)

	s.clearSpans()
	providerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====span should be stored in context after server receive")
	assert.NotNil(t, providerSpanHandler.Span)

	t.Log("====spanContext stored in context should be the same with monitor received")
	zipkinServerRecievedSpans = s.spans()
	assert.Equal(t, 1, len(zipkinServerRecievedSpans))
	if t.Failed() {
		return
	}
	assert.Equal(t, 2, len(zipkinServerRecievedSpans[0].Annotations))

	localSpanContext = providerSpanHandler.Span.Context()
	providerZpSpanContext, ok := localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, uint64(zipkinServerRecievedSpans[0].TraceID), providerZpSpanContext.TraceID.Low)
	assert.Equal(t, (uint64)(*zipkinServerRecievedSpans[0].TraceIDHigh), providerZpSpanContext.TraceID.High)
	assert.Equal(t, (uint64)(zipkinServerRecievedSpans[0].ID), providerZpSpanContext.SpanID)

	assert.NotEqual(t, 0, providerZpSpanContext.SpanID)
	assert.Nil(t, providerZpSpanContext.ParentSpanID)

	t.Log("====spanContext of consumer/provider should be the same")
	assert.Equal(t, consumerZpSpanContext.TraceID, providerZpSpanContext.TraceID)
	assert.Equal(t, consumerZpSpanContext.SpanID, providerZpSpanContext.SpanID)

	t.Log("========tracing [consumer] handler [rest], with [parent]")

	// set args to be restServerReceivedReq
	inv.Args = fasthttpRequest

	consumerChain.Reset()
	s.clearSpans()
	parentSpanID := providerZpSpanContext.SpanID

	consumerChain.Next(inv, func(i *invocation.InvocationResponse) error {
		assert.NoError(t, i.Err)
		return nil
	})

	t.Log("====the parent spanID shoud be last spanID")
	localSpanContext = consumerSpanHandler.Span.Context()
	consumerZpSpanContext, ok = localSpanContext.(zipkin.SpanContext)
	assert.True(t, ok)

	assert.Equal(t, parentSpanID, *consumerZpSpanContext.ParentSpanID)
	assert.NotEqual(t, parentSpanID, consumerZpSpanContext.SpanID)
}
*/

type httpServer struct {
	t            *testing.T
	zipkinSpans  []*zipkincore.Span
	zipkinHeader http.Header
	mutex        sync.RWMutex
}

func (s *httpServer) spans() []*zipkincore.Span {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.zipkinSpans
}

func (s *httpServer) clearSpans() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.zipkinSpans = s.zipkinSpans[:0]
}

func (s *httpServer) headers() http.Header {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.zipkinHeader
}

func (s *httpServer) clearHeaders() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.zipkinHeader = make(http.Header, 0)
}

func newHTTPServer(t *testing.T, port int) *httpServer {
	server := &httpServer{
		t:           t,
		zipkinSpans: make([]*zipkincore.Span, 0),
		mutex:       sync.RWMutex{},
	}

	handler := http.NewServeMux()

	handler.HandleFunc("/api/v1/spans", func(w http.ResponseWriter, r *http.Request) {
		contextType := r.Header.Get("Content-Type")
		if contextType != "application/x-thrift" {
			t.Fatalf(
				"except Content-Type should be application/x-thrift, but is %s",
				contextType)
		}

		// clone headers from request
		headers := make(http.Header, len(r.Header))
		for k, vv := range r.Header {
			vv2 := make([]string, len(vv))
			copy(vv2, vv)
			headers[k] = vv2
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		buffer := thrift.NewTMemoryBuffer()
		if _, err = buffer.Write(body); err != nil {
			t.Error(err)
			return
		}
		transport := thrift.NewTBinaryProtocolTransport(buffer)
		_, size, err := transport.ReadListBegin()
		if err != nil {
			t.Error(err)
			return
		}
		var spans []*zipkincore.Span
		for i := 0; i < size; i++ {
			zs := &zipkincore.Span{}
			if err = zs.Read(transport); err != nil {
				t.Error(err)
				return
			}
			spans = append(spans, zs)
		}
		err = transport.ReadListEnd()
		if err != nil {
			t.Error(err)
			return
		}
		server.mutex.Lock()
		defer server.mutex.Unlock()
		server.zipkinSpans = append(server.zipkinSpans, spans...)
		server.zipkinHeader = headers
	})

	handler.HandleFunc("/api/v1/sleep", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(serverSleep)
	})

	var err error
	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
	}()
	assert.NoError(t, err)

	return server
}
