package apm

import (
	"github.com/go-chassis/go-chassis-apm"
	"github.com/go-chassis/go-chassis-apm/tracing"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"strconv"
)

//monitoring.yaml
const (
	URI        = "URI"
	ServerType = "serverType"
)

var troption tracing.TracingOptions

//CreateEntrySpan use invocation to make spans for apm
func CreateEntrySpan(i *invocation.Invocation) (interface{}, error) {
	openlogging.Debug("CreateEntrySpan:" + i.MicroServiceName)
	spanCtx := tracing.SpanContext{Ctx: i.Ctx, OperationName: i.MicroServiceName + i.URLPathFormat, ParTraceCtx: i.Headers(), Method: i.Protocol, URL: i.MicroServiceName + i.URLPathFormat}
	span, err := apm.CreateEntrySpan(&spanCtx, troption)
	if err != nil {
		openlogging.Error("CreateEntrySpan err:" + err.Error())
		var span interface{}
		return span, err
	}
	i.Ctx = spanCtx.Ctx
	return span, nil
}

//CreateExitSpan use invocation to make spans for apm
func CreateExitSpan(i *invocation.Invocation) (interface{}, error) {
	openlogging.Debug("CreateExitSpan:" + i.MicroServiceName)
	spanCtx := tracing.SpanContext{Ctx: i.Ctx, OperationName: i.MicroServiceName + i.URLPathFormat, ParTraceCtx: i.Headers(), Method: i.Protocol, URL: i.MicroServiceName + i.URLPathFormat, Peer: i.Endpoint + i.URLPathFormat, TraceCtx: map[string]string{}}
	span, err := apm.CreateExitSpan(&spanCtx, troption)
	if err != nil {
		openlogging.Error("CreateExitSpan err:" + err.Error())
		var span interface{}
		return span, err
	}
	for k, v := range spanCtx.TraceCtx { //ctx need transfer by header
		i.SetHeader(k, v)
	}
	return span, nil
}

//EndSpan use invocation to make spans of apm end
func EndSpan(span interface{}, status int) error {
	openlogging.Debug("EndSpan " + strconv.Itoa(status))
	apm.EndSpan(span, status, troption)
	return nil
}

//Init apm
func Init() error {
	openlogging.Debug("apm Init " + config.GetAPM().Tracing.Tracer)
	if config.GetAPM().Tracing.Tracer != "" && config.GetAPM().Tracing.Settings != nil && config.GetAPM().Tracing.Settings[URI] != "" {
		troption = tracing.TracingOptions{APMName: config.GetAPM().Tracing.Tracer, MicServiceName: config.MicroserviceDefinition.ServiceDescription.Name, ServerURI: config.GetAPM().Tracing.Settings["URI"]}
		if serverType, ok := config.GetAPM().Tracing.Settings[ServerType]; ok { //
			troption.MicServiceType, _ = strconv.Atoi(serverType)
		}
		apm.Init(troption)
		openlogging.Info("apm Init:" + config.GetAPM().Tracing.Tracer + " service:" + config.MicroserviceDefinition.ServiceDescription.Name)
	} else {
		openlogging.Warn("apm Init failed. check apm config " + config.GetAPM().Tracing.Tracer)
	}
	return nil
}
