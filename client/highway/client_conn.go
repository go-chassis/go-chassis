package highway

import (
	"bufio"
	"fmt"
	"github.com/ServiceComb/go-chassis/client/highway/pb"
	"github.com/ServiceComb/go-chassis/core/lager"
	"net"
	"sync"
)

const (
	DefaultReadBufferSize  = 0
	DefaultWriteBufferSize = 1024
)

type HighwayClientConnection struct {
	remoteAddr string
	baseConn   net.Conn
	client     *HighwayBaseClient
	mtx        sync.Mutex
	closed     bool
}

func NewHighwayClientConnection(conn net.Conn, client *HighwayBaseClient) *HighwayClientConnection {
	tmp := new(HighwayClientConnection)
	//conn.SetKeepAlive(true)
	tmp.baseConn = conn
	tmp.client = client
	tmp.closed = false
	return tmp
}

func (this *HighwayClientConnection) Open() error {
	err := this.Hello()
	if err != nil {
		this.Close()
		return err
	}
	go this.MsgRecvLoop()
	return nil
}

func (this *HighwayClientConnection) Hello() error {
	wBuf := bufio.NewWriterSize(this.baseConn, DefaultWriteBufferSize)
	protoObj := &HighWayProtocalObject{}
	protoObj.GenerateHelloReq(wBuf)
	err := wBuf.Flush()
	if err != nil {
		return err
	}
	rdBuf := bufio.NewReaderSize(this.baseConn, DefaultReadBufferSize)
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

func (this *HighwayClientConnection) Close() {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	if this.closed {
		return
	}
	this.closed = true
	this.baseConn.Close()
}

func (this *HighwayClientConnection) MsgRecvLoop() {
	rdBuf := bufio.NewReaderSize(this.baseConn, DefaultReadBufferSize)
	for {
		protoObj := &HighWayProtocalObject{}
		err := protoObj.DeSerializeFrame(rdBuf)
		if err != nil {
			break
		}
		this.ProcessMsg(protoObj)
	}
	this.Close()
}

func (this *HighwayClientConnection) ProcessMsg(protoObj *HighWayProtocalObject) {
	ctx := this.client.GetWaitMsg(protoObj.FrHead.MsgID)
	if ctx != nil {
		protoObj.DeSerializeRsp(ctx.Rsp)
		ctx.Done()
	} else {
		lager.Logger.Info(fmt.Sprintf("Cann't find the msg, perhaps it's timeout:%d", protoObj.FrHead.MsgID))
	}
}

func (this *HighwayClientConnection) AsyncSendMsg(ctx *InvocationContext) error {
	wBuf := bufio.NewWriterSize(this.baseConn, DefaultWriteBufferSize)
	protoObj := &HighWayProtocalObject{}
	protoObj.SerializeReq(ctx.Req, wBuf)
	err := wBuf.Flush()
	if err != nil {
		ctx.Rsp.Err = err.Error()
		ctx.Done()
	}
	return err
}

func (this *HighwayClientConnection) PostMsg(req *HighwayRequest) error {
	// Respond of postMsg  is  needless
	wBuf := bufio.NewWriterSize(this.baseConn, DefaultWriteBufferSize)
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
func (this *HighwayClientConnection) Closed() bool {
	return this.closed
}
