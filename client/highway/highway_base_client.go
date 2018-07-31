package highway

import (
	"crypto/tls"
	"errors"
	"github.com/go-chassis/go-chassis/core/lager"
	"net"
	"sync"
	"time"
)

//ConnParams highway connect parmas
type ConnParams struct {
	Addr      string
	TLSConfig *tls.Config
	Timeout   time.Duration
	ConnNum   int
}

//BaseClient highway base client
type BaseClient struct {
	addr          string
	mtx           sync.Mutex
	mapMutex      sync.Mutex
	msgWaitRspMap map[uint64]*InvocationContext
	highwayConns  []*ClientConnection
	closed        bool
	connParams    *ConnParams
}

//CachedClients client cache
var CachedClients *ClientMgr

func init() {
	CachedClients = newClientMgr()
}

//InvocationContext Highway context
type InvocationContext struct {
	Req  *Request
	Rsp  *Response
	Wait *chan int
	once sync.Once
}

//Done Notify done.
func (ctx *InvocationContext) Done() {
	ctx.once.Do(func() {
		close(*ctx.Wait)
	})
}

//ClientMgr client manage
type ClientMgr struct {
	mapMutex sync.Mutex
	clients  map[string]*BaseClient
}

func newClientMgr() *ClientMgr {
	tmp := new(ClientMgr)
	tmp.clients = make(map[string]*BaseClient)
	return tmp
}

//GetClient Obtain  client
func (mgr *ClientMgr) GetClient(connParmas *ConnParams) (*BaseClient, error) {
	mgr.mapMutex.Lock()
	defer mgr.mapMutex.Unlock()
	if tmp, ok := mgr.clients[connParmas.Addr]; ok {
		if !tmp.Closed() {
			//GetClient from cached addr
			return tmp, nil
		}
		delete(mgr.clients, connParmas.Addr)
	}

	lager.Logger.Info("GetClient from new open addr:" + connParmas.Addr)
	tmp := newHighwayBaseClient(connParmas)
	err := tmp.Open()
	if err != nil {
		return nil, err
	}
	mgr.clients[connParmas.Addr] = tmp
	return tmp, nil

}

func newHighwayBaseClient(connParmas *ConnParams) *BaseClient {
	tmp := &BaseClient{}
	tmp.addr = connParmas.Addr
	tmp.closed = true
	tmp.connParams = connParmas
	tmp.msgWaitRspMap = make(map[uint64]*InvocationContext)
	return tmp
}

//GetAddr Obtain the address
func (baseClient *BaseClient) GetAddr() string {
	return baseClient.addr
}

func (baseClient *BaseClient) makeConnection() (*ClientConnection, error) {
	var baseConn net.Conn
	var errDial error

	if baseClient.connParams.TLSConfig != nil {
		dialer := &net.Dialer{Timeout: baseClient.connParams.Timeout * time.Second}
		baseConn, errDial = tls.DialWithDialer(dialer, "tcp", baseClient.addr, baseClient.connParams.TLSConfig)
	} else {
		baseConn, errDial = net.DialTimeout("tcp", baseClient.addr, baseClient.connParams.Timeout*time.Second)
	}
	if errDial != nil {
		lager.Logger.Error("the addr: "+baseClient.addr, errDial)
		return nil, errDial
	}
	highwayConn := NewHighwayClientConnection(baseConn, baseClient)
	err := highwayConn.Open()
	if err != nil {
		lager.Logger.Error("highwayConn open: "+baseClient.addr, errDial)
		return nil, err
	}

	return highwayConn, nil
}

func (baseClient *BaseClient) initConns() error {
	if baseClient.connParams.ConnNum == 0 {
		baseClient.connParams.ConnNum = 4
	}

	baseClient.highwayConns = make([]*ClientConnection, baseClient.connParams.ConnNum)
	for i := 0; i < baseClient.connParams.ConnNum; i++ {
		higwayConn, err := baseClient.makeConnection()
		if err != nil {
			return err
		}
		baseClient.highwayConns[i] = higwayConn
	}
	return nil
}

//Open  client
func (baseClient *BaseClient) Open() error {
	baseClient.mtx.Lock()
	defer baseClient.mtx.Unlock()
	err := baseClient.initConns()
	if err != nil {
		baseClient.clearConns()
		return err
	}
	baseClient.closed = false
	return nil
}

//Close client
func (baseClient *BaseClient) Close() {
	baseClient.mtx.Lock()
	defer baseClient.mtx.Unlock()
	baseClient.close()
}

//close client, no mutex
func (baseClient *BaseClient) close() {
	if baseClient.closed {
		return
	}
	baseClient.mapMutex.Lock()
	for _, v := range baseClient.msgWaitRspMap {
		v.Done()
	}
	baseClient.msgWaitRspMap = make(map[uint64]*InvocationContext)
	baseClient.mapMutex.Unlock()
	baseClient.clearConns()
	baseClient.closed = true
}

func (baseClient *BaseClient) clearConns() {
	for i := 0; i < baseClient.connParams.ConnNum; i++ {
		conn := baseClient.highwayConns[i]
		if conn != nil {
			conn.Close()
			baseClient.highwayConns[i] = nil
		}
	}
}

//AddWaitMsg add wait msg
func (baseClient *BaseClient) AddWaitMsg(msgID uint64, result *InvocationContext) {
	baseClient.mapMutex.Lock()
	if baseClient.msgWaitRspMap != nil {
		baseClient.msgWaitRspMap[msgID] = result
	}
	baseClient.mapMutex.Unlock()
}

//RemoveWaitMsg remove wait msg
func (baseClient *BaseClient) RemoveWaitMsg(msgID uint64) {
	baseClient.mapMutex.Lock()
	if baseClient.msgWaitRspMap != nil {
		delete(baseClient.msgWaitRspMap, msgID)
	}
	baseClient.mapMutex.Unlock()
}

//Send send msg
func (baseClient *BaseClient) Send(req *Request, rsp *Response, timeout time.Duration) error {
	if baseClient.closed {
		baseClient.mtx.Lock()
		if baseClient.closed {
			baseClient.mtx.Unlock()
			return errors.New("client is closed")
		}
		baseClient.mtx.Unlock()
	}

	msgID := req.MsgID
	idx := msgID % uint64(baseClient.connParams.ConnNum)
	highwayConn := baseClient.highwayConns[idx]
	if highwayConn == nil || highwayConn.Closed() {
		baseClient.mtx.Lock()
		highwayConn = baseClient.highwayConns[idx]
		if highwayConn == nil || highwayConn.Closed() {
			highwayConnTmp, err := baseClient.makeConnection()
			if err != nil {
				baseClient.mtx.Unlock()
				return err
			}
			highwayConn = highwayConnTmp
			baseClient.highwayConns[idx] = highwayConn
		}
		baseClient.mtx.Unlock()
	}
	if req.TwoWay {
		wait := make(chan int)
		ctx := &InvocationContext{Req: req, Rsp: rsp, Wait: &wait}
		baseClient.AddWaitMsg(msgID, ctx)

		err := highwayConn.AsyncSendMsg(ctx)
		if err != nil {
			baseClient.RemoveWaitMsg(msgID)
			rsp.Err = err.Error()
			lager.Logger.Error("AsyncSendMsg err:", err)
			return err
		}

		var bTimeout bool
		select {
		case <-wait:
			bTimeout = false
		case <-time.After(timeout * time.Second):
			bTimeout = true
		}

		baseClient.RemoveWaitMsg(msgID)

		if bTimeout {
			ctx.Done()
		}

		if bTimeout {
			rsp.Err = "Client send timeout"
			return errors.New("Client send timeout")
		}
		if ctx.Rsp.Status != Ok {
			return errors.New(ctx.Rsp.Err)
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

//GetWaitMsg get wait message
func (baseClient *BaseClient) GetWaitMsg(msgID uint64) *InvocationContext {
	baseClient.mapMutex.Lock()
	defer baseClient.mapMutex.Unlock()

	if _, ok := baseClient.msgWaitRspMap[msgID]; ok {
		return baseClient.msgWaitRspMap[msgID]
	}
	return nil
}

//Closed  client status
func (baseClient *BaseClient) Closed() bool {
	return baseClient.closed
}
