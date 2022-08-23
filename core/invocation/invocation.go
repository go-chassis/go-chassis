package invocation

import (
	"context"
	"net/http"

	"github.com/go-chassis/go-chassis/v2/core/common"
	"github.com/go-chassis/go-chassis/v2/pkg/runtime"
	utiltags "github.com/go-chassis/go-chassis/v2/pkg/util/tags"
)

// constant values for consumer and provider
const (
	Consumer = iota
	Provider
)
const (
	MDMark = "mark"
)

// Response is invocation response struct
type Response struct {
	Status int
	Result interface{}
	Err    error
}

// ResponseCallBack process invocation response
type ResponseCallBack func(*Response)

// Invocation is the basic struct which makes transport layer transparent to middleware "handler chain".
// developer should implements a client which is able to transfer invocation to there native protocol request,
// a protocol server should transfer request to invocation and then back to request
type Invocation struct {
	//Invocation is a stateful struct,
	//this index indicates the current index number of handler of a chain
	HandlerIndex       int
	SSLEnable          bool   // indicates whether provider service using TLS communication to serve or not
	Endpoint           string // service's ip and port, it is decided in load balancing or specified by invoker
	Protocol           string // indicates consumer what to use which protocol to communicate with provider
	PortName           string // indicates the name of a service port number
	SourceServiceID    string
	SourceMicroService string
	MicroServiceName   string // provider micro service name that consumer want to request

	// route tags is decided in router handler, it indicates metadata of a microservice,
	// like service version, env, etc.
	RouteTags utiltags.Tags

	SchemaID    string // correspond struct name
	OperationID string // correspond func name of struct

	// relative API path of http request
	URLPath string

	// it holds native request of protocol, use http protocol for example,
	// it is *http.request
	Args interface{}

	// it holds native response of protocol, use http protocol for example,
	// in consumer it is *http.response.
	// in provider it is *http.ResponseWriter
	Reply interface{}

	Ctx      context.Context        // ctx can save protocol headers
	Metadata map[string]interface{} // can save local data, will not send in header on network
	Strategy string                 // load balancing strategy
	Filters  []string
	// http.CheckRedirect
	CheckRedirect func(req *http.Request, via []*http.Request) error
}

// GetMark return match rule name that request matches
func (inv *Invocation) GetMark() string {
	m, ok := inv.Metadata[MDMark].(string)
	if ok {
		return m
	}
	return "none"
}

// Mark marks a invocation, it means the invocation matches a match rule
// so that governance rule can be applied to invocation with specific mark
func (inv *Invocation) Mark(matchRuleName string) {
	inv.Metadata[MDMark] = matchRuleName
}

// New create invocation, context can not be nil
// if you don't set ContextHeaderKey, then New will init it
func New(ctx context.Context) *Invocation {
	inv := &Invocation{
		SourceServiceID: runtime.ServiceID,
		Ctx:             ctx,
	}
	if inv.Ctx == nil {
		inv.Ctx = context.TODO()
	}
	if inv.Ctx.Value(common.ContextHeaderKey{}) == nil {
		inv.Ctx = context.WithValue(inv.Ctx, common.ContextHeaderKey{}, map[string]string{})
	}
	inv.Metadata = make(map[string]interface{}, 1)
	inv.Metadata[MDMark] = "none"
	return inv
}

// SetMetadata local scope data
func (inv *Invocation) SetMetadata(key string, value interface{}) {
	if inv.Metadata == nil {
		inv.Metadata = make(map[string]interface{})
	}
	inv.Metadata[key] = value
}

// SetHeader set headers of protocol request, the client and server plugins should use them in protocol headers
// it is convenience but has lower performance than you use Headers[k]=v, when you have a batch of kv to set
func (inv *Invocation) SetHeader(k, v string) {
	m := inv.Ctx.Value(common.ContextHeaderKey{}).(map[string]string)
	m[k] = v
}

// Headers return a map that protocol plugin should deliver in transport
func (inv *Invocation) Headers() map[string]string {
	return inv.Ctx.Value(common.ContextHeaderKey{}).(map[string]string)
}

// Header return header value
func (inv *Invocation) Header(name string) string {
	m := inv.Ctx.Value(common.ContextHeaderKey{}).(map[string]string)
	return m[name]
}
