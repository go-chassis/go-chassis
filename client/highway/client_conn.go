package highway

import (
	"bufio"
	"fmt"
	"github.com/go-chassis/go-chassis/client/highway/pb"
	"github.com/go-chassis/go-chassis/core/lager"
	"net"
	"sync"
)

// constant for buffer size
const (
	DefaultReadBufferSize  = 0
	DefaultWriteBufferSize = 1024
)

//ClientConnection Highway client connection
type ClientConnection struct {
	remoteAddr string
	baseConn   net.Conn
	client     *BaseClient
	mtx        sync.Mutex
	closed     bool
}

//NewHighwayClientConnection creat Highway client connection
func NewHighwayClientConnection(conn net.Conn, client *BaseClient) *ClientConnection {
	tmp := new(ClientConnection)
	//conn.SetKeepAlive(true)
	tmp.baseConn = conn
	tmp.client = client
	tmp.closed = false
	return tmp
}

//Open Init Highway client connection
func (hwClientConn *ClientConnection) Open() error {
	err := hwClientConn.Hello()
	if err != nil {
		hwClientConn.Close()
		return err
	}
	go hwClientConn.msgRecvLoop()
	return nil
}

//Hello Highway handshake
func (hwClientConn *ClientConnection) Hello() error {
	wBuf := bufio.NewWriterSize(hwClientConn.baseConn, DefaultWriteBufferSize)
	protoObj := &ProtocolObject{}
	protoObj.SerializeHelloReq(wBuf)
	err := wBuf.Flush()
	if err != nil {
		return err
	}
	rdBuf := bufio.NewReaderSize(hwClientConn.baseConn, DefaultReadBufferSize)
	rsp := &Response{}
	rsp.Result = &highway.LoginResponse{}
	protoObj.DeSerializeFrame(rdBuf)
	if err != nil {
		return err
	}

	err = protoObj.DeSerializeRsp(rsp)
	if err != nil && rsp.Status != Ok {
		return err
	}
	return nil
}

//Close the connection
func (hwClientConn *ClientConnection) Close() {
	hwClientConn.mtx.Lock()
	defer hwClientConn.mtx.Unlock()
	if hwClientConn.closed {
		return
	}
	hwClientConn.closed = true
	hwClientConn.baseConn.Close()
}

func (hwClientConn *ClientConnection) msgRecvLoop() {
	rdBuf := bufio.NewReaderSize(hwClientConn.baseConn, DefaultReadBufferSize)
	for {
		protoObj := &ProtocolObject{}
		err := protoObj.DeSerializeFrame(rdBuf)
		if err != nil {
			break
		}
		hwClientConn.processMsg(protoObj)
	}
	hwClientConn.Close()
}

func (hwClientConn *ClientConnection) processMsg(protoObj *ProtocolObject) {
	ctx := hwClientConn.client.GetWaitMsg(protoObj.FrHead.MsgID)
	if ctx != nil {
		protoObj.DeSerializeRsp(ctx.Rsp)
		ctx.Done()
	} else {
		lager.Logger.Info(fmt.Sprintf("Cann't find the msg, perhaps it's timeout:%d", protoObj.FrHead.MsgID))
	}
}

//AsyncSendMsg Highway send message
func (hwClientConn *ClientConnection) AsyncSendMsg(ctx *InvocationContext) error {
	wBuf := bufio.NewWriterSize(hwClientConn.baseConn, DefaultWriteBufferSize)
	protoObj := &ProtocolObject{}
	protoObj.SerializeReq(ctx.Req, wBuf)
	err := wBuf.Flush()
	if err != nil {
		ctx.Rsp.Err = err.Error()
		ctx.Done()
	}
	return err
}

//PostMsg Highway post message,	 Respond  is  needless
func (hwClientConn *ClientConnection) PostMsg(req *Request) error {

	wBuf := bufio.NewWriterSize(hwClientConn.baseConn, DefaultWriteBufferSize)
	protoObj := &ProtocolObject{}
	protoObj.SerializeReq(req, wBuf)
	return wBuf.Flush()
}

/*
func (this *ClientConnection) SyncSendMsg(req *client.Request, rsp *client.Response) error {
	wBuf := bufio.NewWriterSize(this.baseConn, 1024)
	protoObj := &HighWayProtocalObject_10{}
	protoObj.SerializeReq(req, wBuf)
	err := wBuf.Flush()
	if err != nil {
		//this.client.RspCallBack(nil)
		return err
	}
	rdBuf := bufio.NewReaderSize(this.baseConn, 1024)

	err = protoObj.DeSerializeFrame(rdBuf)
	if err != nil {
		//this.client.RspCallBack(nil)
		return err
	}

	protoObj.DeSerializeRsp(rsp)
	if err != nil {
		//this.client.RspCallBack(nil)
		return err
	}

	return  err
}
*/

//Closed Highway connection status
func (hwClientConn *ClientConnection) Closed() bool {
	return hwClientConn.closed
}
