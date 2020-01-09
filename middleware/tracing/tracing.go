package tracing

import "github.com/SkyAPM/go2sky"

// TracingClient for apm interface
type TracingClient interface {
	CreateEntrySpan(sc *SpanContext) (go2sky.Span, error)
	CreateExitSpan(sc *SpanContext) (go2sky.Span, error)
	EndSpan(sp go2sky.Span, statusCode int) error
}

var tc TracingClient

//CreateEntrySpan create entry span
func CreateEntrySpan(s *SpanContext) (go2sky.Span, error) {
	return tc.CreateEntrySpan(s)
}

//CreateExitSpan create exit span
func CreateExitSpan(s *SpanContext) (go2sky.Span, error) {
	return tc.CreateExitSpan(s)
}

//EndSpan end span
func EndSpan(span go2sky.Span, status int) error {
	return tc.EndSpan(span, status)
}
