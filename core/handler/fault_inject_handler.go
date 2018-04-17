package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ServiceComb/go-chassis/client/rest"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/fault"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"

	"github.com/valyala/fasthttp"
)

// constant for fault handler name
const (
	FaultHandlerName = "fault-inject"
)

// FaultHandler handler
type FaultHandler struct{}

// FaultHandle fault handle gives the object of FaultHandler
func FaultHandle() Handler {
	return &FaultHandler{}
}

// init is for to register the fault handler
func init() {
	RegisterHandler(FaultHandlerName, FaultHandle)
}

// Name function returns fault-inject string
func (rl *FaultHandler) Name() string {
	return "fault-inject"
}

// Handle is to handle the API
func (rl *FaultHandler) Handle(chain *Chain, inv *invocation.Invocation, cb invocation.ResponseCallBack) {

	faultStruct := GetFaultConfig(inv.Protocol, inv.MicroServiceName, inv.SchemaID, inv.OperationID)
	faultConfig := model.FaultProtocolStruct{}
	faultConfig.Fault = make(map[string]model.Fault)
	faultConfig.Fault[inv.Protocol] = faultStruct

	faultInject, ok := fault.FaultInjectors[inv.Protocol]
	r := &invocation.InvocationResponse{}
	if !ok {
		lager.Logger.Warnf("fault injection doesn't support for protocol ", errors.New(inv.Protocol))
		r.Err = nil
		cb(r)
		return
	}

	faultValue := faultConfig.Fault[inv.Protocol]
	err := faultInject(faultValue, inv)
	if err != nil {
		if strings.Contains(err.Error(), "injecting abort") {
			switch inv.Reply.(type) {
			case *rest.Response:
				resp := inv.Reply.(*rest.Response)
				resp.SetStatusCode(faultConfig.Fault[inv.Protocol].Abort.HTTPStatus)
			case *fasthttp.Response:
				resp := inv.Reply.(*fasthttp.Response)
				resp.SetStatusCode(faultConfig.Fault[inv.Protocol].Abort.HTTPStatus)
			}
			r.Status = faultConfig.Fault[inv.Protocol].Abort.HTTPStatus
		} else {
			switch inv.Reply.(type) {
			case *rest.Response:
				resp := inv.Reply.(*rest.Response)
				resp.SetStatusCode(http.StatusBadRequest)
			case *fasthttp.Response:
				resp := inv.Reply.(*fasthttp.Response)
				resp.SetStatusCode(http.StatusBadRequest)
			}
			r.Status = http.StatusBadRequest
		}

		r.Err = fault.FaultError{Message: err.Error()}
		cb(r)
		return
	}

	chain.Next(inv, func(r *invocation.InvocationResponse) error {
		return cb(r)
	})
}

// GetFaultConfig get faultconfig
func GetFaultConfig(protocol, microServiceName, schemaID, operationID string) model.Fault {

	faultStruct := model.Fault{}
	faultStruct.Abort.Percent = config.GetAbortPercent(protocol, microServiceName, schemaID, operationID)
	faultStruct.Abort.HTTPStatus = config.GetAbortStatus(protocol, microServiceName, schemaID, operationID)
	faultStruct.Delay.Percent = config.GetDelayPercent(protocol, microServiceName, schemaID, operationID)
	faultStruct.Delay.FixedDelay = config.GetFixedDelay(protocol, microServiceName, schemaID, operationID)

	return faultStruct
}
