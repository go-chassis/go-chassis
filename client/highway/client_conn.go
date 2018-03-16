package highway

import (
	"bufio"
	"fmt"
	"github.com/ServiceComb/go-chassis/client/highway/pb"
	"github.com/ServiceComb/go-chassis/core/lager"
	"net"
	"sync"
)

// constant for buffer size
const (
	DefaultReadBufferSize  = 0
	DefaultWriteBufferSize = 1024
)

//Highway client connection
type HighwayClientConnection struct {
	remoteAddr string
	baseConn   net.Conn
	client     *HighwayBaseClient
	mtx        sync.Mutex
	closed     bool
}

//creat Highway client connection
func NewHighwayClientConnection(conn net.Conn, client *HighwayBaseClient) *HighwayClientConnection {
	tmp := new(HighwayClientConnection)
	//conn.SetKeepAlive(true)
	tmp.baseConn = conn
	tmp.client = client
	tmp.closed = false
	return tmp
}

//Init Highway client connection
func (this *HighwayClientConnection) Open() error {
	err := this.Hello()
	if err != nil {
		this.Close()
		return err
	}
	go this.msgRecvLoop()
	return nil
}

//Highway handshake
func (hwClientConn *HighwayClientConnection) Hello() error {
	wBuf := bufio.NewWriterSize(hwClientConn.baseConn, DefaultWriteBufferSize)
	protoObj := &HighWayProtocalObject{}
	protoObj.SerializeHelloReq(wBuf)
	err := wBuf.Flush()
	if err != nil {
		return err
	}
	rdBuf := bufio.NewReaderSize(hwClientConn.baseConn, DefaultReadBufferSize)
	rsp := &HighwayRespond{}
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
func (hwClientConn *HighwayClientConnection) Close() {
	hwClientConn.mtx.Lock()
	defer hwClientConn.mtx.Unlock()
	if hwClientConn.closed {
		return
	}
	hwClientConn.closed = true
	hwClientConn.baseConn.Close()
}

func (hwClientConn *HighwayClientConnection) msgRecvLoop() {
	rdBuf := bufio.NewReaderSize(hwClientConn.baseConn, DefaultReadBufferSize)
	for {
		protoObj := &HighWayProtocalObject{}
		err := protoObj.DeSerializeFrame(rdBuf)
		if err != nil {
			break
		}
		hwClientConn.processMsg(protoObj)
	}
	hwClientConn.Close()
}

func (hwClientConn *HighwayClientConnection) processMsg(protoObj *HighWayProtocalObject) {
	ctx := hwClientConn.client.GetWaitMsg(protoObj.FrHead.MsgID)
	if ctx != nil {
		protoObj.DeSerializeRsp(ctx.Rsp)
		ctx.Done()
	} else {
		lager.Logger.Info(fmt.Sprintf("Cann't find the msg, perhaps it's timeout:%d", protoObj.FrHead.MsgID))
	}
}

//Highway send message
func (hwClientConn *HighwayClientConnection) AsyncSendMsg(ctx *InvocationContext) error {
	wBuf := bufio.NewWriterSize(hwClientConn.baseConn, DefaultWriteBufferSize)
	protoObj := &HighWayProtocalObject{}
	protoObj.SerializeReq(ctx.Req, wBuf)
	err := wBuf.Flush()
	if err != nil {
		ctx.Rsp.Err = err.Error()
		ctx.Done()
	}
	return err
}

//Highway post message,	 Respond  is  needless
func (hwClientConn *HighwayClientConnection) PostMsg(req *HighwayRequest) error {

	wBuf := bufio.NewWriterSize(hwClientConn.baseConn, DefaultWriteBufferSize)
	protoObj := &HighWayProtocalObject{}
	protoObj.SerializeReq(req, wBuf)
	return wBuf.Flush()
}

/*
func (this *HighwayClientConnection) SyncSendMsg(req *client.Request, rsp *client.Response) error {
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

//Highway connection status
func (hwClientConn *HighwayClientConnection) Closed() bool {
	return hwClientConn.closed
}
