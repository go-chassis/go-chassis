package highway

import (
	"crypto/tls"
	"errors"
	"github.com/ServiceComb/go-chassis/core/lager"
	"net"
	"sync"
	"time"
)

type ConnParams struct {
	Addr      string
	TlsConfig *tls.Config
	Timeout   time.Duration
	ConnNum   int
}

type HighwayBaseClient struct {
	addr          string
	mtx           sync.Mutex
	mapMutex      sync.Mutex
	msgWaitRspMap map[uint64]*InvocationContext
	highwayConns  []*HighwayClientConnection
	closed        bool
	connParams    *ConnParams
}

var CachedClients *ClientMgr

func init() {
	CachedClients = NewClientMgr()
}

type InvocationContext struct {
	Req  *HighwayRequest
	Rsp  *HighwayRespond
	Wait *chan int
}

func (this *InvocationContext) Done() {
	*this.Wait <- 1
}

type ClientMgr struct {
	mapMutex sync.Mutex
	clients  map[string]*HighwayBaseClient
}

func NewClientMgr() *ClientMgr {
	tmp := new(ClientMgr)
	tmp.clients = make(map[string]*HighwayBaseClient)
	return tmp
}

func (this *ClientMgr) GetClient(connParmas *ConnParams) (*HighwayBaseClient, error) {
	this.mapMutex.Lock()
	defer this.mapMutex.Unlock()
	if tmp, ok := this.clients[connParmas.Addr]; ok {
		if !tmp.Closed() {
			lager.Logger.Info("GetClient from cached addr:" + connParmas.Addr)
			return tmp, nil
		} else {
			delete(this.clients, connParmas.Addr)
		}
	}

	lager.Logger.Info("GetClient from new open addr:" + connParmas.Addr)
	tmp := NewHighwayBaseClient(connParmas)
	err := tmp.Open()
	if err != nil {
		return nil, err
	} else {
		this.clients[connParmas.Addr] = tmp
		return tmp, nil
	}
}

func NewHighwayBaseClient(connParmas *ConnParams) *HighwayBaseClient {
	tmp := &HighwayBaseClient{}
	tmp.addr = connParmas.Addr
	tmp.closed = true
	tmp.connParams = connParmas
	tmp.msgWaitRspMap = make(map[uint64]*InvocationContext)
	return tmp
}

func (this *HighwayBaseClient) GetAddr() string {
	return this.addr
}

func (this *HighwayBaseClient) makeConnection() (*HighwayClientConnection, error) {
	var baseConn net.Conn
	var errDial error

	if this.connParams.TlsConfig != nil {
		dialer := &net.Dialer{Timeout: this.connParams.Timeout * time.Second}
		baseConn, errDial = tls.DialWithDialer(dialer, "tcp", this.addr, this.connParams.TlsConfig)
	} else {
		baseConn, errDial = net.DialTimeout("tcp", this.addr, this.connParams.Timeout*time.Second)
	}
	if errDial != nil {
		lager.Logger.Error("the addr: "+this.addr, errDial)
		return nil, errDial
	}
	higwayConn := NewHighwayClientConnection(baseConn, this)
	err := higwayConn.Open()
	if err != nil {
		lager.Logger.Error("higwayConn open: "+this.addr, errDial)
		return nil, err
	}

	return higwayConn, nil
}

func (this *HighwayBaseClient) initConns() error {
	if this.connParams.ConnNum == 0 {
		this.connParams.ConnNum = 4
	}

	this.highwayConns = make([]*HighwayClientConnection, this.connParams.ConnNum)
	for i := 0; i < this.connParams.ConnNum; i++ {
		higwayConn, err := this.makeConnection()
		if err != nil {
			return err
		}
		this.highwayConns[i] = higwayConn
	}
	return nil
}

func (this *HighwayBaseClient) Open() error {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	err := this.initConns()
	if err != nil {
		this.clearConns()
		return err
	}
	this.closed = false
	return nil
}

func (this *HighwayBaseClient) Close() {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.close()
}

func (this *HighwayBaseClient) close() {
	if this.closed {
		return
	}
	this.mapMutex.Lock()
	for _, v := range this.msgWaitRspMap {
		v.Done()
	}
	this.msgWaitRspMap = make(map[uint64]*InvocationContext)
	this.mapMutex.Unlock()
	this.clearConns()
	this.closed = true
}

func (this *HighwayBaseClient) clearConns() {
	for i := 0; i < this.connParams.ConnNum; i++ {
		conn := this.highwayConns[i]
		if conn != nil {
			conn.Close()
			this.highwayConns[i] = nil
		}
	}
}

func (this *HighwayBaseClient) AddWaitMsg(msgID uint64, result *InvocationContext) {
	this.mapMutex.Lock()
	if this.msgWaitRspMap != nil {
		this.msgWaitRspMap[msgID] = result
	}
	this.mapMutex.Unlock()
}

func (this *HighwayBaseClient) RemoveWaitMsg(msgID uint64) {
	this.mapMutex.Lock()
	if this.msgWaitRspMap != nil {
		delete(this.msgWaitRspMap, msgID)
	}
	this.mapMutex.Unlock()
}

func (this *HighwayBaseClient) Send(req *HighwayRequest, rsp *HighwayRespond, timeout time.Duration) error {
	if this.closed {
		this.mtx.Lock()
		if this.closed {
			this.mtx.Unlock()
			return errors.New("Client is closed.")
		}
		this.mtx.Unlock()
	}

	msgID := req.MsgID
	idx := msgID % uint64(this.connParams.ConnNum)
	highwayConn := this.highwayConns[idx]
	if highwayConn == nil || highwayConn.Closed() {
		this.mtx.Lock()
		highwayConn = this.highwayConns[idx]
		if highwayConn == nil || highwayConn.Closed() {
			highwayConnTmp, err := this.makeConnection()
			if err != nil {
				this.mtx.Unlock()
				return err
			}
			highwayConn = highwayConnTmp
			this.highwayConns[idx] = highwayConn
		}
		this.mtx.Unlock()
	}
	if req.TwoWay {
		wait := make(chan int)
		ctx := &InvocationContext{req, rsp, &wait}
		this.AddWaitMsg(msgID, ctx)

		err := highwayConn.AsyncSendMsg(ctx)
		if err != nil {
			rsp.Err = err.Error()
			lager.Logger.Error("AsyncSendMsg err:", err)
			return err
		}

		var bTimeout bool = false
		select {
		case <-wait:
			bTimeout = false
		case <-time.After(timeout * time.Second):
			bTimeout = true
		}

		this.RemoveWaitMsg(msgID)
		close(wait)
		if bTimeout {
			rsp.Err = "Client send timeout"
			return errors.New("Client send timeout")
		} else {
			return nil
		}
	} else {
		// Respond of postMsg  is  needless
		err := highwayConn.PostMsg(req)
		if err != nil {
			lager.Logger.Error("PostMsg err:", err)
			return err
		}
	}
	return nil
}

func (this *HighwayBaseClient) GetWaitMsg(msgID uint64) *InvocationContext {
	this.mapMutex.Lock()
	defer this.mapMutex.Unlock()

	if _, ok := this.msgWaitRspMap[msgID]; ok {
		return this.msgWaitRspMap[msgID]
	}
	return nil
}

func (this *HighwayBaseClient) Closed() bool {
	return this.closed
}
