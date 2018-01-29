package tcp

import (
	"errors"
	"fmt"
	"net/url"
	"runtime"
	"sync"

	"github.com/ServiceComb/go-chassis/core/client"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	microClient "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/client"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	"golang.org/x/net/context"
)

const (
	//Name is a variable of type string
	Name = "highway"
)

var requestID int

//初始状态时未连接的

//Go SDK 没有老版本问题 默认本地和对端都是支持新的编码方式的
var localSupportLogin = true

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

	//计算请求Id
	i.ID = requestID
	requestID++
	if requestID >= ((1 << 31) - 2) {
		requestID = 0
	}
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
	address := "highway://" + addr
	u, err := url.Parse(address)
	if err != nil {
		lager.Logger.Errorf(err, "highway get host failed")
		return err
	}

	//check worker number in configuration
	workerNum := config.GlobalDefinition.Cse.Protocols["highway"].WorkerNumber
	if workerNum == 0 {
		workerNum = runtime.NumCPU() * 4
	}

	//check for the existence of workers if not exist create workers
	err = jobSchdlr.createWorkerSchedulers(addr, workerNum, c, u.Host)
	if err != nil {
		return err
	}

	errCh := make(chan error)
	j := &job{
		req:  req,
		resp: rsp,
		err:  errCh,
		ctx:  ctx,
	}

	//schedule the job to workers
	scheduleErr := jobSchdlr.scheduleJob(addr, j)
	if scheduleErr != nil {
		return scheduleErr
	}
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		err = ctx.Err()
		e := fmt.Sprintf("request timeout: %v", ctx.Err())
		return errors.New(e)
	}

	//return nil
}

//TODO send(requestHeader)

func init() {
	client.InstallPlugin(Name, NewHighwayClient)
}
