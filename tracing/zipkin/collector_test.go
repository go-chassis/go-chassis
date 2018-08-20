package zipkin_test

import (
	"testing"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/go-chassis/go-chassis/tracing/zipkin"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/stretchr/testify/assert"
)

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
	byteBuffer := zipkin.Serialize([]*zipkincore.Span{span})
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
