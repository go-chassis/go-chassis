package highway

import (
	"context"
	"sync"

	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/common"
)

//const timeout
const (
	//Name is a variable of type string
	Name                  = "highway"
	DefaultConnectTimeOut = 60
	DefaultSendTimeOut    = 300
)

//higway client
type highwayClient struct {
	once     sync.Once
	opts     client.Options
	reqMutex sync.Mutex // protects following
}

//NewHighwayClient is a function
func NewHighwayClient(options client.Options) client.ProtocolClient {

	rc := &highwayClient{
		once: sync.Once{},
		opts: options,
	}

	c := client.ProtocolClient(rc)

	return c
}

func (c *highwayClient) String() string {
	return "highway_client"
}

func (c *highwayClient) Call(ctx context.Context, addr string, req *client.Request, rsp interface{}) error {
	req.ID = int(GenerateMsgID())
	connParams := &ConnParams{}
	connParams.TLSConfig = c.opts.TLSConfig
	connParams.Addr = addr
	connParams.Timeout = DefaultConnectTimeOut
	baseClient, err := CachedClients.GetClient(connParams)
	if err != nil {
		return err
	}
	tmpRsp := &HighwayRespond{0, Ok, "", 0, rsp, nil}
	highwayReq := &HighwayRequest{}
	highwayReq.MsgID = uint64(req.ID)
	highwayReq.MethodName = req.Operation
	highwayReq.Schema = req.Schema
	highwayReq.Arg = req.Arg
	highwayReq.SvcName = req.MicroServiceName
	//Current only twoway
	highwayReq.TwoWay = true
	highwayReq.Attachments = common.FromContext(ctx)

	err = baseClient.Send(highwayReq, tmpRsp, DefaultSendTimeOut)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	client.InstallPlugin(Name, NewHighwayClient)
}
