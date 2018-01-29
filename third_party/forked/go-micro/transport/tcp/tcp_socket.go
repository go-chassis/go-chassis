// Package tcp provides a TCP transport
package tcp

import (
	"bufio"
	"encoding/gob"
	"net"
	"time"
)

type tcpTransportSocket struct {
	baseConn
	conn    net.Conn
	enc     *gob.Encoder
	dec     *gob.Decoder
	encBuf  *bufio.Writer
	timeout time.Duration
}

func (t *tcpTransportSocket) Recv() ([]byte, []byte, map[string]string, int, error) {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	reqID, requestHeader, requestBody, err := t.baseConn.recvHeader()
	if err != nil {
		return nil, nil, nil, int(reqID), err
	}

	return requestHeader, requestBody, nil, int(reqID), nil
}

func (t *tcpTransportSocket) Send(Header []byte, Body []byte, Metadata map[string]string, ID int) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	headerLen := int32(len(Header))
	totalLen := int32(len(Body)) + headerLen

	//Assemble HighWay CommHeader
	commonHeader, err := t.initCommHeader(ID, totalLen, headerLen)
	if err != nil {
		return err
	}

	err = t.sendFrame(commonHeader, Header, Body)
	if err != nil {
		return err
	}
	return t.baseConn.w.Flush()
}

func (t *tcpTransportSocket) Close() error {
	return t.conn.Close()
}
