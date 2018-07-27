package invocation

import (
	"context"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/pkg/runtime"
	"github.com/ServiceComb/go-chassis/pkg/util/tags"
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
	Port               string
	SourceServiceID    string
	SourceMicroService string
	MicroServiceName   string //Target micro service name
	SchemaID           string //correspond struct name
	OperationID        string //correspond struct func name
	Args               interface{}
	URLPathFormat      string
	Reply              interface{}
	Ctx                context.Context        //ctx can save protocol headers
	Metadata           map[string]interface{} //local scope data
	RouteTags          utiltags.Tags          //route tags is decided in router handler
	Strategy           string                 //load balancing strategy
	Filters            []string
}

// New create invocation
func New(ctx context.Context) *Invocation {
	inv := &Invocation{
		SourceServiceID: runtime.ServiceID,
		Ctx:             ctx,
	}
	return inv
}

//GetSessionID return session id
func (inv *Invocation) GetSessionID() string {
	return inv.Metadata[common.LBSessionID].(string)

}

//SetSessionID set session id to invocation
func (inv *Invocation) SetSessionID(value string) {
	headers := inv.Ctx.Value(common.ContextHeaderKey{}).(map[string]string)
	headers[common.LBSessionID] = value
}

//SetMetadata local scope params
func (inv *Invocation) SetMetadata(key string, value interface{}) {
	if inv.Metadata == nil {
		inv.Metadata = make(map[string]interface{})
	}
	inv.Metadata[key] = value
}

//SetHeader set headers, the client and server plugins should use them in protocol headers
func (inv *Invocation) SetHeader(k, v string) {
	if inv.Ctx.Value(common.ContextHeaderKey{}) == nil {
		inv.Ctx = context.WithValue(inv.Ctx, common.ContextHeaderKey{}, map[string]string{})
	}
	m := inv.Ctx.Value(common.ContextHeaderKey{}).(map[string]string)

	m[k] = v
}

//Headers return a map that protocol plugin should deliver in transport
func (inv *Invocation) Headers() map[string]string {
	return inv.Ctx.Value(common.ContextHeaderKey{}).(map[string]string)
}
