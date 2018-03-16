package highway

// Forked from github.com/micro/go-micro
// Some parts of this file have been modified to make it functional in this package
import (
	"github.com/ServiceComb/go-chassis/core/client"
	"sync"

	microClient "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/metadata"
	"golang.org/x/net/context"
)

const (
	//Name is a variable of type string
	Name                  = "highway"
	DefaultConnectTimeOut = 60
	DefaultSendTimeOut    = 300
)

//higway client
type highwayClient struct {
	once     sync.Once
	opts     microClient.Options
	reqMutex sync.Mutex // protects following
}

func (c *highwayClient) Init(opts ...microClient.Option) error {
	for _, o := range opts {
		o(&c.opts)
	}

	return nil
}

func (c *highwayClient) NewRequest(service, schemaID, operationID string, arg interface{}, reqOpts ...microClient.RequestOption) *microClient.Request {
	var opts microClient.RequestOptions

	for _, o := range reqOpts {
		o(&opts)
	}
	i := &microClient.Request{
		MicroServiceName: service,
		Struct:           schemaID,
		Method:           operationID,
		Arg:              arg,
	}

	i.ID = int(GenerateMsgID())
	return i
}

//NewHighwayClient is a function
func NewHighwayClient(options ...microClient.Option) microClient.Client {
	opts := microClient.Options{
		PoolTTL: microClient.DefaultPoolTTL,
	}
	for _, o := range options {
		o(&opts)
	}

	if opts.Codecs == nil {
		opts.Codecs = codec.GetCodecMap()
	}
	if len(opts.ContentType) == 0 {
		//TODO take effect of that option
		opts.ContentType = "application/protobuf"
	}

	rc := &highwayClient{
		once: sync.Once{},
		opts: opts,
	}

	c := microClient.Client(rc)

	return c
}

func (c *highwayClient) String() string {
	return "highway_client"
}

func (c *highwayClient) Options() microClient.Options {
	return c.opts
}

func (c *highwayClient) Call(ctx context.Context, addr string, req *microClient.Request, rsp interface{}, opts ...microClient.CallOption) error {
	connParams := &ConnParams{}
	connParams.TlsConfig = c.opts.TLSConfig
	connParams.Addr = addr
	connParams.Timeout = DefaultConnectTimeOut
	baseClient, err := CachedClients.GetClient(connParams)
	if err != nil {
		return err
	}
	tmpRsp := &HighwayRespond{0, Ok, "", 0, rsp, nil}
	highwayReq := &HighwayRequest{}
	highwayReq.MsgID = uint64(req.ID)
	highwayReq.MethodName = req.Method
	highwayReq.Schema = req.Struct
	highwayReq.Arg = req.Arg
	highwayReq.SvcName = req.MicroServiceName
	//Current only twoway
	highwayReq.TwoWay = true
	var ok bool
	highwayReq.Attachments, ok = metadata.FromContext(ctx)
	if !ok {
		highwayReq.Attachments = make(map[string]string)
	}
	err = baseClient.Send(highwayReq, tmpRsp, DefaultSendTimeOut)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	client.InstallPlugin(Name, NewHighwayClient)
}
