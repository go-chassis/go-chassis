package handler_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/tracing"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/emicklei/go-restful"
	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/stretchr/testify/assert"
)

const (
	serverSleep = 100 * time.Millisecond
)

type spanTestHandler struct {
	Span opentracing.Span
}

func (s *spanTestHandler) Handle(chain *handler.Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	s.Span = opentracing.SpanFromContext(i.Ctx)
	cb(&invocation.Response{
		Err: nil,
	})
}

func (s *spanTestHandler) Name() string {
	return "spanTestHandler"
}

func TestTracingHandler_Highway_Consumer(t *testing.T) {
	lager.Initialize("", "DEBUG", "", "size", true, 1, 10, 7)
	// the port should be different from test cases
	port := 30110
	s := newHTTPServer(t, port)
	tracing.Init(&tracing.Option{
		ServiceName:     "test",
		CollectorType:   tracing.TracingZipkinCollector,
		CollectorTarget: fmt.Sprintf("http://localhost:%d/api/v1/spans", port),
		ProtocolEndpointMap: map[string]string{
			common.RestMethod:      "localhost:8000",
			common.ProtocolHighway: "localhost:8001",
		},
	})

	testTracer(t, s, common.ProtocolHighway)
	testTracer(t, s, common.ProtocolRest)
}

func testTracer(t *testing.T, s *httpServer, protocol string) {
	baseSpanNum := len(s.spans())
	// set handler chain
	consumerChain := handler.Chain{}
	spanHandler := &spanTestHandler{}
	consumerChain.AddHandler(&handler.TracingConsumerHandler{})
	consumerChain.AddHandler(spanHandler)
	consumerParentOp, consumerInv := generateOpAndInv(t, protocol, common.Consumer, "parent", nil)

	t.Logf("[%s] consumer create parent span", protocol)
	consumerChain.Next(consumerInv, func(i *invocation.Response) error {
		assert.NoError(t, i.Err)
		return nil
	})
	time.Sleep(2 * time.Second)
	assert.True(t, len(s.spans()) == baseSpanNum+1)
	consumerParentSpanFromZipkin, err := filterSpanByOperationName(consumerParentOp, s.spans()...)
	assert.NoError(t, err)
	consumerParentSpanContext, ok := spanHandler.Span.Context().(zipkin.SpanContext)
	assert.True(t, ok)
	assert.Equal(t, uint64(consumerParentSpanFromZipkin.ID), consumerParentSpanContext.SpanID)

	t.Log("====span should stored in ctx")
	assert.NotNil(t, opentracing.SpanFromContext(consumerInv.Ctx))

	t.Logf("[%s] provider create span from carrier", protocol)
	providerChain := handler.Chain{}
	providerSpanHandler := &spanTestHandler{}
	providerChain.AddHandler(&handler.TracingProviderHandler{})
	providerChain.AddHandler(providerSpanHandler)
	providerFromCarrierOp, providerInv := generateOpAndInv(t, protocol, common.Provider, "fromcarrier", consumerInv)
	providerChain.Next(providerInv, func(i *invocation.Response) error {
		assert.NoError(t, i.Err)
		return nil
	})
	time.Sleep(2 * time.Second)
	assert.True(t, len(s.spans()) == baseSpanNum+2)
	providerFromCarrierSpanFromZipkin, err := filterSpanByOperationName(providerFromCarrierOp, s.spans()...)
	assert.NoError(t, err)
	providerFromCarrierSpanContext, ok := providerSpanHandler.Span.Context().(zipkin.SpanContext)
	assert.True(t, ok)
	assert.Equal(t, uint64(providerFromCarrierSpanFromZipkin.ID), providerFromCarrierSpanContext.SpanID)
	assert.Equal(t, providerFromCarrierSpanContext.SpanID, consumerParentSpanContext.SpanID)

	t.Logf("[%s] consumer create child span", protocol)
	childOp, consumerInv := generateOpAndInv(t, protocol, common.Consumer, "child", providerInv)
	consumerChain.Reset()
	consumerChain.Next(consumerInv, func(i *invocation.Response) error {
		assert.NoError(t, i.Err)
		return nil
	})
	time.Sleep(2 * time.Second)
	assert.True(t, len(s.spans()) == baseSpanNum+3)
	consumerChildSpanFromZipkin, err := filterSpanByOperationName(childOp, s.spans()...)
	assert.NoError(t, err)
	assert.Equal(t, consumerParentSpanFromZipkin.ID, *consumerChildSpanFromZipkin.ParentID)

	t.Logf("[%s] provider create new span", protocol)
	providerNewOp, providerInv := generateOpAndInv(t, protocol, common.Provider, "new", nil)
	providerChain.Reset()
	providerChain.Next(providerInv, func(i *invocation.Response) error {
		assert.NoError(t, i.Err)
		return nil
	})
	time.Sleep(2 * time.Second)
	assert.True(t, len(s.spans()) == baseSpanNum+4)
	providerNewSpanFromZipkin, err := filterSpanByOperationName(providerNewOp, s.spans()...)
	assert.NoError(t, err)
	providerNewSpanContext, ok := providerSpanHandler.Span.Context().(zipkin.SpanContext)
	assert.True(t, ok)
	assert.Equal(t, uint64(providerNewSpanFromZipkin.ID), providerNewSpanContext.SpanID)
	assert.NotEqual(t, providerNewSpanContext.TraceID, consumerParentSpanContext.TraceID)
}

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

	h := http.NewServeMux()

	h.HandleFunc("/api/v1/spans", func(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Get %d new spans", len(spans))
	})

	h.HandleFunc("/api/v1/sleep", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(serverSleep)
	})

	var err error
	go func() {
		err = http.ListenAndServe(fmt.Sprintf(":%d", port), h)
	}()
	assert.NoError(t, err)

	return server
}

func filterSpanByOperationName(o string, span ...*zipkincore.Span) (*zipkincore.Span, error) {
	for _, s := range span {
		if s.Name == o {
			return s, nil
		}
	}
	return nil, errors.New("no matched span")
}

func generateOpAndInv(t *testing.T, proto, serviceType, op string, old *invocation.Invocation) (name string, inv *invocation.Invocation) {
	switch proto {
	case common.ProtocolRest:
		baseUri := "http://localhost:8000"
		name = strings.Join([]string{"", proto, serviceType, op}, "/")
		inv = &invocation.Invocation{
			MicroServiceName: "test",
			Protocol:         proto,
		}
		switch serviceType {
		case common.Consumer:
			req, err := rest.NewRequest("Get", baseUri+name)
			assert.NoError(t, err)
			inv.Args = req
			inv.URLPathFormat = req.Req.URL.Path
			// this means consumer needs inherit old's ctx
			if old != nil {
				inv.Ctx = old.Ctx
			} else {
				inv.Ctx = common.NewContext(nil)
			}
		default:
			tmpReq, err := rest.NewRequest("GET", baseUri+name)
			assert.NoError(t, err)
			if old != nil {
				req, ok := old.Args.(*rest.Request)
				assert.True(t, ok)
				if !ok {
					t.FailNow()
				}
				tmpReq.Req.Header = req.Req.Header
			}
			providerFromCarrierReq := restful.NewRequest(tmpReq.Req)
			inv.Args = providerFromCarrierReq
			inv.URLPathFormat = providerFromCarrierReq.Request.URL.Path
			inv.Ctx = common.NewContext(nil)
		}
	default:
		strings.Join([]string{proto, serviceType, op}, ".")
		name = strings.Join([]string{proto, serviceType, op}, ".")
		inv = &invocation.Invocation{
			MicroServiceName: "test",
			Protocol:         proto,
			OperationID:      name,
		}
		if old != nil {
			inv.Ctx = old.Ctx
		}
	}
	if inv.Ctx == nil {
		inv.Ctx = common.NewContext(nil)
	}
	return
}
