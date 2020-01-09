package tracing

import "context"

//SpanContext for span info
type SpanContext struct {
	Ctx           context.Context
	OperationName string
	ParTraceCtx   map[string]string
	TraceCtx      map[string]string
	Peer          string
	Method        string
	URL           string
	ComponentID   string
	SpanLayerID   string
	ServiceName   string
}
