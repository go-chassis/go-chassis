package tcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"strings"
	//"strconv"
	"github.com/ServiceComb/go-chassis/core/util/string"
)

//中软协议公共头
var magID = "CSE.TCP"

//中软协议关键字
var magicID = [7]byte{0x43, 0x53, 0x45, 0x2E, 0x54, 0x43, 0x50}

type baseConn struct {
	w *bufio.Writer
	r *bufio.Reader
	c io.Closer
	//frameBuf [binary.MaxVarintLen64]byte
	frameBuf      [23]byte //保存公共消息头
	sendBuf       bytes.Buffer
	commHeaderBuf bytes.Buffer
}

// Close closes the underlying connection.
func (c *baseConn) Close() error {
	return c.c.Close()
}

//组装HighWay协议|commHeader|requestHeader|Body|
func (c *baseConn) initCommHeader(id int, totalLen int32, headLen int32) ([]byte, error) {
	//TODO 公共里没有携带编码信息
	//新的消息格式：|magic number(7)|request ID(8)|total length(4)|header length(4)|
	//Assemble CommHeader
	ID := uint64(id)
	c.commHeaderBuf.Reset()
	binary.Write(&c.commHeaderBuf, binary.BigEndian, magicID)
	binary.Write(&c.commHeaderBuf, binary.BigEndian, ID)              //uint64
	binary.Write(&c.commHeaderBuf, binary.BigEndian, int32(totalLen)) //int32
	binary.Write(&c.commHeaderBuf, binary.BigEndian, int32(headLen))  //int32

	//Assemble requestHeader
	return c.commHeaderBuf.Bytes(), nil
}

//TODO 修改基本功能为再组装码流
//1、Make clientBuffer as common buffer of client
func (c *baseConn) sendFrame(comm []byte, Header []byte, Body []byte) error {
	//1、写requestHeader
	c.sendBuf.Reset()
	err := binary.Write(&c.sendBuf, binary.BigEndian, comm)
	if err != nil {
		return err
	}

	//2、写requestHeader
	err = binary.Write(&c.sendBuf, binary.BigEndian, Header)
	if err != nil {
		return err
	}
	//3、写requestBody
	err = binary.Write(&c.sendBuf, binary.BigEndian, Body)
	if err != nil {
		return err
	}

	return c.write(c.w, c.sendBuf.Bytes())
}

func (c *baseConn) write(w io.Writer, data []byte) error {
	for index := 0; index < len(data); {
		n, err := w.Write(data[index:])
		if err != nil {
			if nerr, ok := err.(net.Error); !ok || !nerr.Temporary() {
				return err
			}
		}
		index += n
	}

	//logger.Logger.Debug("send data buffer :", zap.Binary("data", data))
	return nil
}

func (c *baseConn) recv(r io.Reader, data []byte) error {
	_, err := io.ReadFull(r, data)
	if err != nil {
		return err
	}
	return nil
}

//recvDiscard discard msg body
func (c *baseConn) recvDiscard(bodyLen int32) error {

	//获取消息头长度
	if bodyLen == 0 {
		return nil
	}
	data := make([]byte, bodyLen)
	err := c.recv(c.r, data)
	if err != nil {
		return err
	}
	return nil
}

//解析公用头和请求头
//这里不带解码功能 只负责从IO中拦截码流
func (c *baseConn) recvHeader() (uint64, []byte, []byte, error) {
	//解析公共头
	if err := c.recv(c.r, c.frameBuf[:]); err != nil {
		//TODO performance low
		return 0, nil, nil, err
	}
	//判断魔鬼数字
	if !strings.EqualFold(magID, stringutil.Bytes2str(c.frameBuf[0:7])) {
		return 0, nil, nil, errors.New("MagicID Err")
	}

	reqID := binary.BigEndian.Uint64(c.frameBuf[7:15])
	totalLength := binary.BigEndian.Uint32(c.frameBuf[15:19])
	headLength := binary.BigEndian.Uint32(c.frameBuf[19:23])
	bodyLength := totalLength - headLength

	if bodyLength < 0 {
		return reqID, nil, nil, errors.New("Parse Message Length Err")
	}

	requestHeader := make([]byte, headLength) //TODO performance low
	if err := c.recv(c.r, requestHeader); err != nil {
		return reqID, requestHeader, nil, err
	}

	requestBody := make([]byte, bodyLength)
	if err := c.recv(c.r, requestBody); err != nil {
		return reqID, requestHeader, requestBody, err
	}

	return reqID, requestHeader, requestBody, nil
}

func (c *baseConn) recvBody(m interface{}, bodyLen int32) error {
	if bodyLen == 0 {
		return nil
	}
	requestBody := make([]byte, bodyLen)
	if err := c.recv(c.r, requestBody); err != nil {
		return err
	}
	return protoUnmarshal(requestBody, m)
}

func protoUnmarshal(src []byte, msg interface{}) error {
	return proto.Unmarshal(src, msg.(proto.Message))
}
