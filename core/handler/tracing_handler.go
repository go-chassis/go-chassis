package handler

import (
	"context"
	"errors"
	"net/url"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/tracing"
	"github.com/ServiceComb/go-chassis/util/iputil"

	"github.com/emicklei/go-restful"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/openzipkin/zipkin-go-opentracing/thrift/gen-go/zipkincore"
	"github.com/valyala/fasthttp"
)

// TracingProviderHandler tracing provider handler
type TracingProviderHandler struct{}

// Handle is to handle the provider tracing related things
func (t *TracingProviderHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	tracer := t.getTracer(i)

	var (
		wireContext   opentracing.SpanContext
		err           error
		interfaceName = "unknown"
		carrier       interface{}
	)

	// extract span context
	switch i.Protocol {
	case common.ProtocolRest:
		// header stored in args, in rest
		// handler chain doesn't resolve rest request, so never return
		switch i.Args.(type) {
		case *restful.Request:
			req := i.Args.(*restful.Request)
			carrier = (opentracing.HTTPHeadersCarrier)(req.Request.Header)
		case *fasthttp.Request:
			req := i.Args.(*fasthttp.Request)
			headerMap := make(map[string]string)
			req.Header.VisitAll(func(key, value []byte) {
				headerMap[string(key)] = string(value)
			})
			carrier = &tracing.HeaderCarrier{Header: headerMap}
		default:
			lager.Logger.Error("rest consumer call arg is neither *restful.Request|*fasthttp.Request type.", nil)
			err = errors.New("Type invalid")
		}
		if err != nil {
			break
		}
		// set url path to span name
		if u, e := url.Parse(i.URLPathFormat); e != nil {
			lager.Logger.Error("parse request url failed.", e)
		} else {
			interfaceName = u.Path
		}
	default:
		interfaceName = i.OperationID

		// header stored in context
		if i.Ctx == nil {
			lager.Logger.Debug("No metadata found in Invocation.Ctx")
			break
		}
		at, ok := i.Ctx.Value(common.ContextValueKey{}).(map[string]string)
		if !ok {
			lager.Logger.Debug("No metadata found in Invocation.Ctx")
			break
		}
		carrier = (opentracing.TextMapCarrier)(at)
	}

	wireContext, err = tracer.Extract(
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
	operationName := genOperaitonName(i.MicroServiceName, interfaceName)
	span := tracer.StartSpan(operationName, ext.RPCServerOption(wireContext))
	// set span kind to be server
	ext.SpanKindRPCServer.Set(span)

	// store span in context
	newCtx := opentracing.ContextWithSpan(i.Ctx, span)
	i.Ctx = newCtx

	// To ensure accuracy, spans should finish immediately once server responds.
	// So the best way is that spans finish in the callback func, not after it.
	// But server may respond in the callback func too, that we have to remove
	// span finishing from callback func's inside to outside.
	chain.Next(i, func(r *invocation.InvocationResponse) (err error) {
		err = cb(r)
		switch i.Protocol {
		case common.ProtocolRest:
			span.SetTag(zipkincore.HTTP_METHOD, i.Metadata[common.RestMethod])
			span.SetTag(zipkincore.HTTP_PATH, interfaceName)
			span.SetTag(zipkincore.HTTP_STATUS_CODE, r.Status)
			span.SetTag(zipkincore.HTTP_HOST, i.Endpoint)
		default:
		}
		span.Finish()
		return
	})
}

// Name returns tracing-provider string
func (t *TracingProviderHandler) Name() string {
	return TracingProvider
}

func (t *TracingProviderHandler) getTracer(i *invocation.Invocation) opentracing.Tracer {
	caller := i.MicroServiceName + ":" + iputil.GetHostName()
	return tracing.GetTracer(caller)
}

func newTracingProviderHandler() Handler {
	return &TracingProviderHandler{}
}

// TracingConsumerHandler tracing consumer handler
type TracingConsumerHandler struct{}

// Handle is handle consumer tracing related things
func (t *TracingConsumerHandler) Handle(chain *Chain, i *invocation.Invocation, cb invocation.ResponseCallBack) {
	// the span context is in invocation.Ctx
	// TODO distinguish rpc msg type
	tracer := t.getTracer(i)

	if i.Ctx == nil {
		i.Ctx = context.Background()
	}

	// start a new span from context
	var span opentracing.Span
	opts := make([]opentracing.StartSpanOption, 0)

	interfaceName := "unknown"
	interfaceName = setInterfaceName(interfaceName, i)

	operationName := genOperaitonName(i.MicroServiceName, interfaceName)
	if parentSpan := opentracing.SpanFromContext(i.Ctx); parentSpan != nil {
		opts = append(opts, opentracing.ChildOf(parentSpan.Context()))
		span = tracer.StartSpan(operationName, opts...)
	} else {
		span = tracer.StartSpan(operationName, opts...)
	}
	// set span kind to be client
	ext.SpanKindRPCClient.Set(span)

	// store span in context
	i.Ctx = opentracing.ContextWithSpan(i.Ctx, span)

	// inject span context into carrier
	switch i.Protocol {
	case common.ProtocolRest:
		var carrier interface{}
		var err error
		// header stored in args, in rest consumer
		// handler chain doesn't resolve rest request, so never return
		switch i.Args.(type) {
		case *rest.Request:
			req := i.Args.(*rest.Request)
			carrier = (*tracing.RestClientHeaderWriter)(req)
		case *fasthttp.Request:
			carrier = &(i.Args.(*fasthttp.Request).Header)
		default:
			lager.Logger.Error("rest consumer call arg is neither *rest.Request|*fasthttp.Request type.", nil)
			err = errors.New("Type invalid")
		}
		if err != nil {
			break
		}

		if err = tracer.Inject(
			span.Context(),
			opentracing.TextMap,
			carrier,
		); err != nil {
			lager.Logger.Errorf(err, "Inject span failed")
		}
	default:
		// header stored in context
		var header map[string]string
		attachments, ok := i.Ctx.Value(common.ContextValueKey{}).(map[string]string)
		if !ok {
			header = make(map[string]string)
			i.Ctx = context.WithValue(i.Ctx, common.ContextValueKey{}, header)
		} else {
			header = attachments
		}

		if err := tracer.Inject(
			span.Context(),
			opentracing.TextMap,
			(opentracing.TextMapCarrier)(header),
		); err != nil {
			lager.Logger.Errorf(err, "Inject span failed")
		} else {
			i.Ctx = context.WithValue(i.Ctx, common.ContextValueKey{}, header)
		}
	}

	// To ensure accuracy, spans should finish immediately once client send req.
	// So the best way is that spans finish in the callback func, not after it.
	// But client may send req in the callback func too, that we have to remove
	// span finishing from callback func's inside to outside.
	chain.Next(i, func(r *invocation.InvocationResponse) (err error) {
		err = cb(r)
		switch i.Protocol {
		case common.ProtocolRest:
			span.SetTag(zipkincore.HTTP_METHOD, i.Metadata[common.RestMethod])
			span.SetTag(zipkincore.HTTP_PATH, interfaceName)
			span.SetTag(zipkincore.HTTP_STATUS_CODE, r.Status)
			span.SetTag(zipkincore.HTTP_HOST, i.Endpoint)
		default:
		}
		span.Finish()
		return
	})
}

func setInterfaceName(interfaceName string, i *invocation.Invocation) string {
	switch i.Protocol {
	case common.ProtocolRest:
		// set url path to span name
		if u, e := url.Parse(i.URLPathFormat); e != nil {
			lager.Logger.Error("parse request url failed.", e)
		} else {
			interfaceName = u.Path
		}
	default:
		interfaceName = i.OperationID
	}

	return interfaceName
}

// Name returns tracing-consumer string
func (t *TracingConsumerHandler) Name() string {
	return TracingConsumer
}

func (t *TracingConsumerHandler) getTracer(i *invocation.Invocation) opentracing.Tracer {
	caller := common.DefaultValue
	if c, ok := i.Metadata[common.CallerKey].(string); ok && c != "" {
		caller = c + ":" + iputil.GetHostName()
	}
	return tracing.GetTracer(caller)
}

func newTracingConsumerHandler() Handler {
	return &TracingConsumerHandler{}
}

func genOperaitonName(microserviceName, interfaceName string) string {
	return "[" + microserviceName + "]:[" + interfaceName + "]"
}
