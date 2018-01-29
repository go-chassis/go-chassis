package tcp

import (
	"bufio"
	"encoding/gob"
	"net"
	"time"

	transportOption "github.com/ServiceComb/go-chassis/third_party/forked/go-micro/transport"
)

type tcpTransportClient struct {
	baseConn
	dialOpts transportOption.DialOptions
	conn     net.Conn
	enc      *gob.Encoder
	dec      *gob.Decoder
	encBuf   *bufio.Writer
	timeout  time.Duration
}

func (t *tcpTransportClient) Send(Header []byte, Body []byte, metadata map[string]string, ID int) error {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}

	headerLen := int32(len(Header))
	totalLen := int32(len(Body)) + headerLen

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

func (t *tcpTransportClient) Recv() ([]byte, []byte, map[string]string, int, error) {
	// set timeout if its greater than 0
	if t.timeout > time.Duration(0) {
		t.conn.SetDeadline(time.Now().Add(t.timeout))
	}
	reqID, responseHeader, responseBody, err := t.baseConn.recvHeader()
	if err != nil {
		return nil, nil, nil, int(reqID), err
	}

	return responseHeader, responseBody, nil, int(reqID), nil
}

func (t *tcpTransportClient) Close() error {
	return t.conn.Close()
}
