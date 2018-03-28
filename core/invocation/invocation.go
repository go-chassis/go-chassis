package invocation

import (
	"context"

	"github.com/ServiceComb/go-chassis/core/config"

	"github.com/ServiceComb/go-chassis/core/loadbalancer"
)

// constant values for consumer and provider
const (
	Consumer = iota
	Provider
)

// InvocationResponse invocation response struct
type InvocationResponse struct {
	Status int
	Result interface{}
	Err    error
}

// ResponseCallBack process invocation response
type ResponseCallBack func(*InvocationResponse) error

//Invocation is the basic struct that used in go sdk to make client and transport layer transparent .
//developer should implements a client which is able to  encode from invocation to there own request
type Invocation struct {
	//service's ip and port, it is decided in load balancer
	Endpoint string
	//specify rest,highway
	Protocol string
	//Ctx value will be send as header in transport
	Ctx                context.Context
	SourceServiceID    string
	SourceMicroService string
	MicroServiceName   string
	Version            string
	AppID              string
	//correspond struct
	SchemaID string
	//correspond struct func
	OperationID string
	Args        interface{}
	//an url path has muntil path params such as "/v2/microsvice/:id/instance/:instanceid",http client use this to format correct url
	URLPathFormat string

	Reply interface{}
	//just in local memory
	Metadata map[string]interface{}
	//loadbalancer stratery
	//Strategy loadbalancer.Strategy
	Strategy string
	Filters  []loadbalancer.Filter
}

// CreateInvocation create invocation
func CreateInvocation() *Invocation {
	return &Invocation{
		SourceServiceID: config.SelfServiceID,
	}
}
