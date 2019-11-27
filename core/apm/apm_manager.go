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
	APMURI        = "URI"
	APMServerType = "servertype"
)

var troption tracing.TracingOptions

//CreateEntrySpan use invocation to make spans for apm
func CreateEntrySpan(i *invocation.Invocation) (interface{}, error) {
	openlogging.GetLogger().Debug("CreateEntrySpan:" + i.MicroServiceName)
	spanCtx := tracing.SpanContext{Ctx: i.Ctx, OperationName: i.MicroServiceName + i.URLPathFormat, ParTraceCtx: i.Headers(), Method: i.Protocol, URL: i.MicroServiceName + i.URLPathFormat}
	span, err := apm.CreateEntrySpan(&spanCtx, troption)
	if err != nil {
		openlogging.GetLogger().Error("CreateEntrySpan err:" + err.Error())
		var span interface{}
		return span, err
	}
	i.Ctx = spanCtx.Ctx
	return span, nil
}

//CreateExitSpan use invocation to make spans for apm
func CreateExitSpan(i *invocation.Invocation) (interface{}, error) {
	openlogging.GetLogger().Debug("CreateExitSpan:" + i.MicroServiceName)
	spanCtx := tracing.SpanContext{Ctx: i.Ctx, OperationName: i.MicroServiceName + i.URLPathFormat, ParTraceCtx: i.Headers(), Method: i.Protocol, URL: i.MicroServiceName + i.URLPathFormat, Peer: i.Endpoint + i.URLPathFormat, TraceCtx: map[string]string{}}
	span, err := apm.CreateExitSpan(&spanCtx, troption)
	if err != nil {
		openlogging.GetLogger().Error("CreateExitSpan err:" + err.Error())
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
	openlogging.GetLogger().Debug("EndSpan " + strconv.Itoa(status))
	apm.EndSpan(span, status, troption)
	return nil
}

//Init apm
func Init() error {
	openlogging.GetLogger().Debug("Apm Init " + config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer)
	if config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer != "" && config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings != nil && config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings[APMURI] != "" {
		troption = tracing.TracingOptions{APMName: config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer, MicServiceName: config.MicroserviceDefinition.ServiceDescription.Name, ServerURI: config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings["URI"]}
		if serverType, ok := config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings[APMServerType]; ok { //
			troption.MicServiceType, _ = strconv.Atoi(serverType)
		}
		apm.Init(troption)
	} else {
		openlogging.GetLogger().Error("Apm Init failed. check apm config " + config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer)
	}
	openlogging.Info("Apm Init:" + config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer + " micname:" + config.MicroserviceDefinition.ServiceDescription.Name)
	return nil
}
