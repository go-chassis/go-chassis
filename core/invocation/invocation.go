package invocation

import (
	"context"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
)

// constant values for consumer and provider
const (
	Consumer = iota
	Provider
)

// Response is invocation response struct
type Response struct {
	Status int
	Result interface{}
	Err    error
}

// ResponseCallBack process invocation response
type ResponseCallBack func(*Response) error

//Invocation is the basic struct that used in go sdk to make client and transport layer transparent .
//developer should implements a client which is able to  encode from invocation to there own request
type Invocation struct {
	Endpoint           string //service's ip and port, it is decided in load balancing
	Protocol           string
	SourceServiceID    string
	SourceMicroService string
	MicroServiceName   string //Target micro service name
	Version            string
	AppID              string
	SchemaID           string //correspond struct name
	OperationID        string //correspond struct func name
	Args               interface{}
	URLPathFormat      string
	Reply              interface{}
	Ctx                context.Context //ctx can save protocol header
	Metadata           map[string]interface{}
	RouteTags          map[string]string //route tags is decided in router handler
	Strategy           string            //load balancing strategy
	Filters            []string
}

// CreateConsumerInvocation create invocation
func CreateConsumerInvocation() *Invocation {
	return &Invocation{
		SourceServiceID: config.SelfServiceID,
		Metadata:        make(map[string]interface{}),
	}
}

//GetSessionID return session id
func (inv *Invocation) GetSessionID() string {
	return inv.Metadata[common.LBSessionID].(string)

}

//SetSessionID set session id to invocation
func (inv *Invocation) SetSessionID(value string) {
	inv.Metadata[common.LBSessionID] = value
}
