package highway

import (
	"bufio"
	"net"
	"sync"

	highwayclient "github.com/go-chassis/go-chassis/client/highway"
	"github.com/go-chassis/go-chassis/client/highway/pb"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/handler"
	"github.com/go-chassis/go-chassis/core/invocation"
	"github.com/go-chassis/go-chassis/core/lager"
	"github.com/go-chassis/go-chassis/core/provider"
	"io"
)

//ConnectionMgr conn manage
type ConnectionMgr struct {
	conns map[string]*HighwayConnection
	count int
	sync.RWMutex
}

func newConnectMgr() *ConnectionMgr {
	tmp := new(ConnectionMgr)
	tmp.count = 0
	tmp.conns = make(map[string]*HighwayConnection)
	return tmp
}

func (connMgr *ConnectionMgr) createConn(baseConn net.Conn, handlerChain string) *HighwayConnection {
	connMgr.Lock()
	defer connMgr.Unlock()
	conn := newHighwayConnection(baseConn, handlerChain, connMgr)
	connMgr.conns[conn.GetRemoteAddr()] = conn
	connMgr.count++
	return conn
}

func (connMgr *ConnectionMgr) deleteConn(addr string) {
	connMgr.Lock()
	defer connMgr.Unlock()
	delete(connMgr.conns, addr)
	connMgr.count--
}

//DeactiveAllConn close all conn
func (connMgr *ConnectionMgr) DeactiveAllConn() {
	for _, conn := range connMgr.conns {
		conn.Close()
	}
}

//HighwayConnection Highway connection
type HighwayConnection struct {
	remoteAddr   string
	handlerChain string
	baseConn     net.Conn
	mtx          *sync.Mutex
	closed       bool
	connMgr      *ConnectionMgr
}

//newHighwayConnection Create service connection
func newHighwayConnection(conn net.Conn, handlerChain string, connMgr *ConnectionMgr) *HighwayConnection {
	return &HighwayConnection{(conn.(*net.TCPConn)).RemoteAddr().String(), handlerChain, conn, &sync.Mutex{}, false, connMgr}
}

//Open open service connection
func (svrConn *HighwayConnection) Open() {
	go svrConn.msgRecvLoop()
}

//GetRemoteAddr Get remote addr
func (svrConn *HighwayConnection) GetRemoteAddr() string {
	return svrConn.remoteAddr
}

//Close connection
func (svrConn *HighwayConnection) Close() {
	svrConn.mtx.Lock()
	defer svrConn.mtx.Unlock()
	if svrConn.closed {
		return
	}
	svrConn.connMgr.deleteConn(svrConn.remoteAddr)
	svrConn.closed = true
	svrConn.baseConn.Close()
}

//Hello handshake
func (svrConn *HighwayConnection) Hello() error {
	var err error
	rdBuf := bufio.NewReaderSize(svrConn.baseConn, highwayclient.DefaultReadBufferSize)
	protoObj := &highwayclient.ProtocolObject{}
	protoObj.DeSerializeFrame(rdBuf)
	if err != nil {
		return err
	}

	req := &highwayclient.Request{}
	req.Arg = &highway.LoginRequest{}
	err = protoObj.DeSerializeReq(req)
	if err != nil {
		return err
	}

	if loginRequest, ok := req.Arg.(*highway.LoginRequest); ok {
		if loginRequest.UseProtobufMapCodec == true {
			wBuf := bufio.NewWriterSize(svrConn.baseConn, highwayclient.DefaultWriteBufferSize)
			protoObj.SerializeLoginRap(req.MsgID, wBuf)
			err := wBuf.Flush()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (svrConn *HighwayConnection) msgRecvLoop() {
	if svrConn.Hello() != nil {
		//Handshake failed , Close the conn
		svrConn.Close()
		return
	}
	rdBuf := bufio.NewReaderSize(svrConn.baseConn, highwayclient.DefaultReadBufferSize)
	for {
		protoObj := &highwayclient.ProtocolObject{}
		err := protoObj.DeSerializeFrame(rdBuf)
		if err != nil {
			if err != io.EOF {
				lager.Logger.Errorf(err, "DeSerializeFrame failed.")
			}

			break
		}
		go svrConn.handleFrame(protoObj)
	}
	svrConn.Close()
}

//send error msg
func (svrConn *HighwayConnection) writeError(req *highwayclient.Request, err error) {
	if req.TwoWay {
		protoObj := &highwayclient.ProtocolObject{}
		wBuf := bufio.NewWriterSize(svrConn.baseConn, highwayclient.DefaultWriteBufferSize)
		rsp := &highwayclient.Response{}
		rsp.Result = nil
		rsp.MsgID = req.MsgID
		rsp.Err = err.Error()
		rsp.Status = highwayclient.ServerError
		protoObj.SerializeRsp(rsp, wBuf)
		errSnd := wBuf.Flush()
		if errSnd != nil {
			svrConn.Close()
			lager.Logger.Errorf(errSnd, "writeError failed.")
		}
	}
}

func (svrConn *HighwayConnection) handleFrame(protoObj *highwayclient.ProtocolObject) error {
	var err error
	req := &highwayclient.Request{}
	err = protoObj.DeSerializeReq(req)
	if err != nil {
		lager.Logger.Errorf(err, "DeSerializeReq failed")
		svrConn.writeError(req, err)
		return err
	}

	i := invocation.New(common.NewContext(req.Attachments))
	i.Args = req.Arg
	i.MicroServiceName = req.SvcName
	i.SchemaID = req.Schema
	i.OperationID = req.MethodName
	i.SourceMicroService = common.FromContext(i.Ctx)[common.HeaderSourceName]
	i.Protocol = common.ProtocolHighway
	c, err := handler.GetChain(common.Provider, svrConn.handlerChain)
	if err != nil {
		lager.Logger.Errorf(err, "Handler chain init err")
		svrConn.writeError(req, err)
	}

	c.Next(i, func(ir *invocation.Response) error {
		if ir.Err != nil {
			svrConn.writeError(req, ir.Err)
			return ir.Err
		}
		p, err := provider.GetProvider(i.MicroServiceName)
		if err != nil {
			svrConn.writeError(req, err)
			return err
		}
		r, err := p.Invoke(i)
		if err != nil {
			svrConn.writeError(req, err)
			return err
		}
		if req.TwoWay {
			wBuf := bufio.NewWriterSize(svrConn.baseConn, highwayclient.DefaultWriteBufferSize)
			rsp := &highwayclient.Response{}
			rsp.Result = r
			rsp.Status = highwayclient.Ok
			rsp.MsgID = req.MsgID
			protoObj.SerializeRsp(rsp, wBuf)
			err = wBuf.Flush()
			if err != nil {
				lager.Logger.Errorf(err, "Send Respond failed.")
				svrConn.Close()
				return err
			}
		}
		return err

	})

	return nil
}
