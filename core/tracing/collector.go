package tracing

import (
	"bytes"
	"errors"
	"os"

	"github.com/apache/thrift/lib/go/thrift"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
)

// constant for tracing zipkin and pipe collectors
const (
	TracingZipkinCollector    = "zipkin"
	TracingNamedPipeCollector = "namedPipe"
)

//FileCollector collects span to file
type FileCollector struct {
	Fd *os.File
}

// Collect serialize the zipkin spans and write into the file collector
func (f *FileCollector) Collect(s *zipkincore.Span) error {
	buf := Serialize([]*zipkincore.Span{s})
	_, err := f.Fd.Write(buf.Bytes())
	return err
}

// Close close file collector
func (f *FileCollector) Close() error {
	return f.Fd.Close()
}

// Serialize serialize the zipkin spans
func Serialize(spans []*zipkincore.Span) *bytes.Buffer {
	t := thrift.NewTMemoryBuffer()
	p := thrift.NewTBinaryProtocolTransport(t)
	if err := p.WriteListBegin(thrift.STRUCT, len(spans)); err != nil {
		panic(err)
	}
	for _, s := range spans {
		if err := s.Write(p); err != nil {
			panic(err)
		}
	}
	if err := p.WriteListEnd(); err != nil {
		panic(err)
	}
	return t.Buffer
}

// collectorNewer new collector
type collectorNewer func(string) (zipkin.Collector, error)

var newNamedPipeCollector collectorNewer = func(string) (zipkin.Collector, error) {
	return nil, errors.New("OS does not support named pipe")
}

// NewCollector returns the collector object based on collector type
func NewCollector(collectorType, target string) (zipkin.Collector, error) {
	if collectorType == "" {
		collectorType = TracingZipkinCollector
	}
	switch collectorType {
	case TracingZipkinCollector:
		return zipkin.NewHTTPCollector(target, zipkin.HTTPBatchSize(10))
	case TracingNamedPipeCollector:
		return newNamedPipeCollector(target)
	}
	return nil, errors.New("Not support collector type: " + collectorType)
}
