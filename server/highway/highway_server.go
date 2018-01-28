package tcp

import (
	"context"
	"errors"
	"io"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/ServiceComb/go-chassis/client/highway/pb"
	"github.com/ServiceComb/go-chassis/core/codec"
	"github.com/ServiceComb/go-chassis/core/common"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/handler"
	"github.com/ServiceComb/go-chassis/core/invocation"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/ServiceComb/go-chassis/core/server"
	"github.com/ServiceComb/go-chassis/core/util/metadata"
	microServer "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/server"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
	"log"
)

const (
	//Name is a variable of type string which says about the protocol used
	Name = "highway"
)

// constants for request and login
const (
	Request = 0
	Login   = 1
)

var remoteLogin = true

type highwayServer struct {
	tr   transport.Transport
	opts microServer.Options
	exit chan chan error
	sync.RWMutex
}

func (s *highwayServer) Init(opts ...microServer.Option) error {
	s.Lock()
	for _, o := range opts {
		o(&s.opts)
	}
	lager.Logger.Debugf("server init,transport:%s", s.opts.Transport.String())
	s.Unlock()
	return nil
}
func (s *highwayServer) Options() microServer.Options {
	s.RLock()
	opts := s.opts
	s.RUnlock()
	return opts
}
func (s *highwayServer) Register(schema interface{}, options ...microServer.RegisterOption) (string, error) {
	opts := microServer.RegisterOptions{}
	var pn string
	for _, o := range options {
		o(&opts)
	}
	if opts.MicroServiceName == "" {
		opts.MicroServiceName = config.SelfServiceName
	}
	mc := config.MicroserviceDefinition
	if mc == nil {
		pn = common.DefaultProvider
	}
	if mc == nil || mc.Provider == "" {
		pn = common.DefaultProvider
	} else {
		if mc.Provider == "" {
			pn = common.DefaultProvider
		} else {
			pn = mc.Provider
		}

	}
	provider.RegisterProvider(pn, opts.MicroServiceName)
	if opts.SchemaID != "" {
		err := provider.RegisterSchemaWithName(opts.MicroServiceName, opts.SchemaID, schema)
		return opts.SchemaID, err
	}
	schemaID, err := provider.RegisterSchema(opts.MicroServiceName, schema)
	return schemaID, err

}

func (s *highwayServer) Start() error {
	opts := s.Options()
	//TODO lot of options
	l, err := opts.Transport.Listen(opts.Address)
	if err != nil {
		return err
	}
	lager.Logger.Warnf(nil, "Highway server listening on: %s", l.Addr())
	s.Lock()
	s.opts.Address = l.Addr()
	s.Unlock()

	go l.Accept(s.accept)

	go func() {
		ch := <-s.exit
		ch <- l.Close()
	}()
	return nil
}
func (s *highwayServer) Stop() error {
	ch := make(chan error)
	s.exit <- ch
	return <-ch
}
func (s *highwayServer) serveSocket(sock transport.Socket, Header []byte, Body []byte, metadata map[string]string, ID int) {
	rpcRequest := &highway.RequestHeader{}
	rpcResponse := &highway.ResponseHeader{}

	codecFunc := codec.NewPBCodec()
	defer func() {
		if r := recover(); r != nil {
			err := r.(string)
			writeError(rpcResponse, codecFunc, sock, ID, errors.New(err))
		}
	}()

	//默认为OK状态
	rpcResponse.StatusCode = 200
	rpcResponse.Reason = ""
	rpcResponse.Flags = 0

	//解码 请求头
	err := codecFunc.Unmarshal(Header, rpcRequest)
	if err != nil {
		writeError(rpcResponse, codecFunc, sock, ID, err)
		return
	}

	//TODO 请求头是Login
	switch rpcRequest.MsgType {
	case Login:
		err = s.loginHandler(sock, rpcResponse, Body, ID)
		if err != nil {
			lager.Logger.Errorf(err, "highway server deal with login message failed")
			return
		}

	case Request:
		err = s.messageHandler(sock, Body, rpcRequest, rpcResponse, ID, s.opts.ChainName)
		if err != nil {
			lager.Logger.Errorf(err, "highway server deal with request message failed")
			return
		}
	default:
		lager.Logger.Errorf(err, "highway server receive an unknow  message type")
		//TODO 异常请求不断链是否OK
		return
	}
}

func (s *highwayServer) accept(sock transport.Socket) {
	defer func() {
		// close socket
		sock.Close()
		if r := recover(); r != nil {
			lager.Logger.Warnf(nil, string(debug.Stack()), r)
		}
	}()

	for {
		Header, Body, md, ID, err := sock.Recv()
		if err != nil {
			if err != io.EOF {
				lager.Logger.Errorf(err, "Server Receive Err")
			}
			return
		}
		s.serveSocket(sock, Header, Body, md, ID)
	}
}

func newHighwayServer(opts ...microServer.Option) microServer.Server {
	return &highwayServer{
		opts: newOptions(opts...),
		exit: make(chan chan error),
	}
}
func newOptions(opt ...microServer.Option) microServer.Options {
	opts := microServer.Options{
		Metadata: map[string]string{},
	}
	if opts.Codecs == nil {
		opts.Codecs = codec.GetCodecMap()
	}
	for _, o := range opt {
		o(&opts)
	}

	return opts
}
func (s *highwayServer) String() string {
	return Name
}
func init() {
	server.InstallPlugin(Name, newHighwayServer)
}

func writeError(rpcResponse *highway.ResponseHeader,
	codec codec.Codec,
	sock transport.Socket,
	id int,
	err error) {
	lager.Logger.Errorf(err, "highway server socket  error")
	rpcResponse.Reason = err.Error()
	//TODO 505 定义为服务端异常
	rpcResponse.StatusCode = 505

	respBytes, err := codec.Marshal(rpcResponse)
	if err != nil {
		lager.Logger.Errorf(err, "server marshal err")
		return
	}

	//TODO 只发送响应
	err = sock.Send(respBytes, nil, nil, id)
	if err != nil {
		lager.Logger.Errorf(err, "sock send err")
	}
	return
}

func (s *highwayServer) loginHandler(sock transport.Socket, rpcResponse *highway.ResponseHeader, Body []byte, ID int) error {
	codecFunc := codec.NewPBCodec()
	loginRequest := &highway.LoginRequest{}

	//TODO 解码请求头
	err := codecFunc.Unmarshal(Body, loginRequest)
	if err != nil {
		lager.Logger.Errorf(err, "highway server unmarshal loginRequest failed")
		return err
	}

	//header:ResponseHeader
	//Body  :LoginResponse

	if loginRequest.UseProtobufMapCodec == remoteLogin {
		loginHeaderBytes, err := codecFunc.Marshal(rpcResponse)
		if err != nil {
			lager.Logger.Errorf(err, "login server marshal loginRequest failed")
			return err
		}
		//TODO 设置服务端编码支持新的编码方式
		loginResponse := &highway.LoginResponse{
			Protocol:            "highway",
			ZipName:             "z",
			UseProtobufMapCodec: remoteLogin,
		}

		loginResponseBytes, err := codecFunc.Marshal(loginResponse)
		if err != nil {
			lager.Logger.Errorf(err, "login server marshal loginResponse failed")
			return err
		}

		err = sock.Send(loginHeaderBytes, loginResponseBytes, nil, ID)
		if err != nil {
			lager.Logger.Errorf(err, "sock send err")
		}
		return err
	}
	return nil
}

func (s *highwayServer) messageHandler(sock transport.Socket, Body []byte, rpcRequest *highway.RequestHeader, rpcResponse *highway.ResponseHeader, ID int, chainName string) error {

	codecFunc := codec.NewPBCodec()
	op, err := provider.GetOperation(rpcRequest.DestMicroservice, rpcRequest.SchemaID, rpcRequest.OperationName)
	if err != nil {
		writeError(rpcResponse, codecFunc, sock, ID, err)
		return err
	}

	i := &invocation.Invocation{}
	if op != nil && op.Args() != nil && len(op.Args()) > 0 {
		if op.Args()[1].Kind() != reflect.Ptr {
			err = errors.New("second arg not ptr")
			writeError(rpcResponse, codecFunc, sock, ID, err)
			return err
		}

		argv := reflect.New(op.Args()[1].Elem())
		err = codecFunc.Unmarshal(Body, argv.Interface())
		if err != nil {
			writeError(rpcResponse, codecFunc, sock, ID, err)
			return err
		}
		i.Args = argv.Interface()
	}

	i.MicroServiceName = rpcRequest.DestMicroservice
	i.SchemaID = rpcRequest.SchemaID
	i.OperationID = rpcRequest.OperationName
	if rpcRequest.GetContext() != nil {
		i.SourceMicroService = rpcRequest.GetContext()[common.HeaderSourceName]
	}
	i.Ctx = metadata.NewContext(context.Background(), rpcRequest.Context)
	i.Protocol = common.ProtocolHighway

	c, err := handler.GetChain(common.Provider, chainName)
	if err != nil {
		lager.Logger.Errorf(err, "Handler chain init err")
	}
	c.Next(i, func(ir *invocation.InvocationResponse) error {
		if ir.Err != nil {
			writeError(rpcResponse, codecFunc, sock, ID, ir.Err)
			return ir.Err
		}
		p, err := provider.GetProvider(i.MicroServiceName)
		if err != nil {
			return err
		}
		r, err := p.Invoke(i)
		if err != nil {
			return err
		}
		log.Println(r)
		result, err := codecFunc.Marshal(r)
		if err != nil {
			lager.Logger.Errorf(err, "Marshal result error")
			writeError(rpcResponse, codecFunc, sock, ID, err)
			return err
		}

		//Todo Context 带什么信息
		respBytes, err := codecFunc.Marshal(rpcResponse)
		if err != nil {
			writeError(rpcResponse, codecFunc, sock, ID, err)
			return err
		}

		//存放响应头
		err = sock.Send(respBytes, result, nil, ID)
		if err != nil {
			lager.Logger.Errorf(err, "sock send err")
			return err
		}
		return err
	})
	return nil
}
