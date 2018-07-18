package handler

import (
	"errors"
	"net/url"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/tracing"

	"github.com/emicklei/go-restful"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
)

var ErrRestConsumerArgInvalid = errors.New("rest consumer call arg is not *rest.Request type")
var ErrRestProviderArgInvalid = errors.New("rest provider call arg is not *restful.Request type")

// TracingProviderHandler tracing provider handler
type TracingProviderHandler struct{}

// Handle is to handle the provider tracing related things
func (t *TracingProviderHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	// get carrier
	var carrier interface{}
	switch i.Protocol {
	case common.ProtocolRest:
		switch i.Args.(type) {
		case *restful.Request:
			// carrier is http header in rest
			req := i.Args.(*restful.Request)
			carrier = (opentracing.HTTPHeadersCarrier)(req.Request.Header)
		default:
			lager.Logger.Error(ErrRestProviderArgInvalid.Error(), nil)
		}
	default:
		// carrier is header stored in ctx
		if i.Ctx == nil {
			lager.Logger.Debug("No metadata found in Invocation.Ctx")
			break
		}
		at := common.FromContext(i.Ctx)
		carrier = (opentracing.TextMapCarrier)(at)
	}

	// extract span context from carrier
	tracer := t.getTracer(i)
	operationName := getOperationName(i)
	wireContext, err := tracer.Extract(
		opentracing.TextMap,
		carrier,
	)
	switch err {
	case nil:
	case opentracing.ErrSpanContextNotFound:
		lager.Logger.Debug(err.Error())
	default:
		lager.Logger.Errorf(err, "Extract span failed")
	}

	// start new span
	span := tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
	// set span kind to be server
	ext.SpanKindRPCServer.Set(span)
	// use invocation.Ctx to transit span across process
	i.Ctx = opentracing.ContextWithSpan(i.Ctx, span)

	// To ensure accuracy, spans should finish immediately once server responds.
	// So the best way is that spans finish in the callback func, not after it.
	// But server may respond in the callback func too, that we have to remove
	// span finishing from callback func's inside to outside.
	chain.Next(i, func(r *invocation.Response) (err error) {
		err = cb(r)
		span.Finish()
		return
	})
}

// Name returns tracing-provider string
func (t *TracingProviderHandler) Name() string {
	return TracingProvider
}

func (t *TracingProviderHandler) getTracer(i *invocation.Invocation) opentracing.Tracer {
	if t, err := tracing.ProviderTracer(i.Protocol); err != nil {
		return tracing.DefaultTracer()
	} else {
		return t
	}
}

func newTracingProviderHandler() Handler {
	return &TracingProviderHandler{}
}

// TracingConsumerHandler tracing consumer handler
type TracingConsumerHandler struct{}

// Handle is handle consumer tracing related things
func (t *TracingConsumerHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	// start a new span
	// TODO distinguish rpc msg type
	tracer := t.getTracer(i)
	operationName := getOperationName(i)

	var span opentracing.Span
	var carrier interface{}
	opts := make([]opentracing.StartSpanOption, 0)
	// use invocation.Ctx to transit span across process
	// in provider side, start a new span from context
	if parentSpan := opentracing.SpanFromContext(i.Ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
		span = tracer.StartSpan(operationName, opts...)
	} else {
		span = tracer.StartSpan(operationName, opts...)
	}
	// set span kind to be client
	ext.SpanKindRPCClient.Set(span)
	// use invocation.Ctx to transit span across process
	i.Ctx = opentracing.ContextWithSpan(i.Ctx, span)

	// inject span context into carrier
	switch i.Protocol {
	case common.ProtocolRest:
		switch i.Args.(type) {
		case *rest.Request:
			req := i.Args.(*rest.Request)
			carrier = (opentracing.HTTPHeadersCarrier)(req.Req.Header)
		default:
			lager.Logger.Error(ErrRestConsumerArgInvalid.Error(), nil)
		}
	default:
		// header stored in context
		at := common.FromContext(i.Ctx)
		carrier = (opentracing.TextMapCarrier)(at)
	}
	if err := tracer.Inject(
		span.Context(),
		opentracing.TextMap,
		carrier,
	); err != nil {
		lager.Logger.Errorf(err, "Inject span failed")
	}

	// To ensure accuracy, spans should finish immediately once client send req.
	// So the best way is that spans finish in the callback func, not after it.
	// But client may send req in the callback func too, that we have to remove
	// span finishing from callback func's inside to outside.
	chain.Next(i, func(r *invocation.Response) (err error) {
		switch i.Protocol {
		case common.ProtocolRest:
			span.SetTag(zipkincore.HTTP_METHOD, i.Metadata[common.RestMethod])
			span.SetTag(zipkincore.HTTP_PATH, operationName)
			span.SetTag(zipkincore.HTTP_STATUS_CODE, r.Status)
			span.SetTag(zipkincore.HTTP_HOST, i.Endpoint)
		default:
		}
		span.Finish()
		return cb(r)
	})
}

func getOperationName(i *invocation.Invocation) string {
	name := "unknown"
	switch i.Protocol {
	case common.ProtocolRest:
		if u, e := url.Parse(i.URLPathFormat); e != nil {
			lager.Logger.Error("parse request url failed.", e)
		} else {
			// span name is uri path in rest
			name = u.Path
		}
	default:
		// span name is operation id
		name = i.OperationID
	}

	return name
}

// Name returns tracing-consumer string
func (t *TracingConsumerHandler) Name() string {
	return TracingConsumer
}

func (t *TracingConsumerHandler) getTracer(i *invocation.Invocation) opentracing.Tracer {
	return tracing.ConsumerTracer()
}

func newTracingConsumerHandler() Handler {
	return &TracingConsumerHandler{}
}
