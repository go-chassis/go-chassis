package apm

import (
	"github.com/go-chassis/go-chassis-apm"
	"github.com/go-chassis/go-chassis-apm/tracing"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-mesh/openlogging"
	"reflect"
	"strconv"
)

//monitoring.yaml
const (
	APMURI        = "URI"
	APMServerType = "servertype"
)

var micoption tracing.TracingOptions

//CreateEntrySpan use invocation to make spans for apm
func CreateEntrySpan(i *invocation.Invocation) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateEntrySpan. inv:%v", i)
	spanCtx := tracing.SpanContext{Ctx: i.Ctx, OperationName: i.MicroServiceName + i.URLPathFormat, ParTraceCtx: i.Headers(), Method: i.Protocol, URL: i.MicroServiceName + i.URLPathFormat}
	span, err := apm.CreateEntrySpan(&spanCtx, micoption)
	if err != nil {
		openlogging.GetLogger().Errorf("CreateEntrySpan err:%s", err.Error())
		var span interface{}
		return span, err
	}
	i.Ctx = spanCtx.Ctx
	return span, nil
}

//CreateExitSpan use invocation to make spans for apm
func CreateExitSpan(i *invocation.Invocation) (interface{}, error) {
	openlogging.GetLogger().Debugf("CreateExitSpan. inv:%v", i)
	spanCtx := tracing.SpanContext{Ctx: i.Ctx, OperationName: i.MicroServiceName + i.URLPathFormat, ParTraceCtx: i.Headers(), Method: i.Protocol, URL: i.MicroServiceName + i.URLPathFormat, Peer: i.Endpoint + i.URLPathFormat, TraceCtx: map[string]string{}}
	span, err := apm.CreateExitSpan(&spanCtx, micoption)
	if err != nil {
		openlogging.GetLogger().Errorf("CreateExitSpan err:%s", err.Error())
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
	openlogging.GetLogger().Debugf("EndSpan. %v %v %v", span, reflect.TypeOf(span), strconv.Itoa(status))
	apm.EndSpan(span, status, micoption)
	return nil
}

//Init apm
func Init() error {
	openlogging.GetLogger().Debugf("Apm Init %v", *config.MonitorCfgDef)
	if config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer != "" && config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings != nil && config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings[APMURI] != "" {
		micoption = tracing.TracingOptions{APMName: config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer, MicServiceName: config.MicroserviceDefinition.ServiceDescription.Name, ServerUri: config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings["URI"]}
		if serverType, ok := config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings[APMServerType]; ok { //
			micoption.MicServiceType, _ = strconv.Atoi(serverType)
		}
		apm.Init(micoption)
	} else {
		openlogging.GetLogger().Errorf("Apm Init failed. check apm config %v %v %v", config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer, config.MonitorCfgDef.ServiceComb.APM.Tracing.Settings, config.MicroserviceDefinition.ServiceDescription.Name)
	}
	openlogging.Info("Apm Init:" + config.MonitorCfgDef.ServiceComb.APM.Tracing.Tracer + " micname:" + config.MicroserviceDefinition.ServiceDescription.Name)
	return nil
}
