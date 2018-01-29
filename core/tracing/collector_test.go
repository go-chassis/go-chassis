package tracing_test

import (
	"log"
	"os"
	"testing"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/tracing"
	"github.com/apache/thrift/lib/go/thrift"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/stretchr/testify/assert"
)

func TestNewCollector(t *testing.T) {
	log.Println("Test NewCollector")

	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)

	t.Log("========new http collector")
	collector, err := tracing.NewCollector(tracing.TracingZipkinCollector, "http://127.0.0.1/zipkin")
	assert.NoError(t, err)
	httpCollector, ok := collector.(*zipkin.HTTPCollector)
	assert.True(t, ok)
	assert.NotNil(t, httpCollector)

	t.Log("========new named pipe collector")
	target := "namedPipeTracing.log"
	_, err = os.Stat(target)
	if err != nil {
		assert.True(t, os.IsNotExist(err))
	} else {
		err = os.Remove(target)
		assert.NoError(t, err)
	}

	t.Log("========new collector of no-supported type")
	collector, err = tracing.NewCollector("no-support", target)
	assert.NotNil(t, err)
	assert.Nil(t, collector)
}

func TestSerialize(t *testing.T) {
	t.Log("========Test thrift seriliaze")
	var timeStamp int64 = 1
	var duration int64 = 2
	var traceIDHigh int64 = 3
	span := &zipkincore.Span{
		Name:              "test",
		ID:                1,
		Timestamp:         &timeStamp,
		Duration:          &duration,
		TraceID:           1,
		TraceIDHigh:       &traceIDHigh,
		Annotations:       make([]*zipkincore.Annotation, 0),
		BinaryAnnotations: make([]*zipkincore.BinaryAnnotation, 0),
	}
	byteBuffer := tracing.Serialize([]*zipkincore.Span{span})
	buffer := thrift.NewTMemoryBuffer()
	if _, err := buffer.Write(byteBuffer.Bytes()); err != nil {
		t.Error(err)
		return
	}
	transport := thrift.NewTBinaryProtocolTransport(buffer)
	_, size, err := transport.ReadListBegin()
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 1, size)
	spanAfterTransport := &zipkincore.Span{}
	err = spanAfterTransport.Read(transport)
	assert.NoError(t, err)
	assert.Equal(t, span, spanAfterTransport)
}
