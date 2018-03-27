package highway

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/ServiceComb/go-chassis/client/highway/pb"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/core/provider"
	"github.com/ServiceComb/go-chassis/core/util/string"
	"github.com/golang/protobuf/proto"
)

//number const
const (
	FrameHeadLen = 23
	MagicLen     = 7
)

//status code
const (
	Ok          = 200
	ServerError = 505
)

var localSupportLogin = true
var gCurMSGID uint64
var msgIDMtx sync.Mutex

//GenerateMsgID generate message ID
func GenerateMsgID() uint64 {
	msgIDMtx.Lock()
	defer msgIDMtx.Unlock()
	gCurMSGID++
	return gCurMSGID
}

//HighwayRequest Highway request
type HighwayRequest struct {
	MsgID       uint64
	MsgType     int
	TwoWay      bool
	Arg         interface{}
	MethodName  string
	SvcName     string
	Schema      string
	Attachments map[string]string
}

//HighwayRespond Highway respond
type HighwayRespond struct {
	MsgID       uint64
	Status      int
	Err         string
	MsgType     int
	Result      interface{}
	Attachments map[string]string
}

var magID = "CSE.TCP"

var magicID = [MagicLen]byte{0x43, 0x53, 0x45, 0x2E, 0x54, 0x43, 0x50}

type highwayFrameHead struct {
	Magic     [MagicLen]byte
	MsgID     uint64
	TotalLen  uint32
	HeaderLen uint32
}

func (frHead *highwayFrameHead) serialize() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, frHead.Magic)
	binary.Write(buf, binary.BigEndian, frHead.MsgID)
	binary.Write(buf, binary.BigEndian, frHead.TotalLen)
	binary.Write(buf, binary.BigEndian, frHead.HeaderLen)
	return buf.Bytes()
}

func (frHead *highwayFrameHead) deserialize(buf []byte) error {
	if len(buf) < FrameHeadLen {
		return errors.New("Too few bytes")
	}
	rdBuf := bytes.NewBuffer(buf)
	binary.Read(rdBuf, binary.BigEndian, &frHead.Magic)
	if !strings.EqualFold(magID, stringutil.Bytes2str(frHead.Magic[0:])) {
		return errors.New("Invalid magicID")
	}
	binary.Read(rdBuf, binary.BigEndian, &frHead.MsgID)
	binary.Read(rdBuf, binary.BigEndian, &frHead.TotalLen)
	binary.Read(rdBuf, binary.BigEndian, &frHead.HeaderLen)
	return nil
}

func newHeadFrame(msgID uint64) *highwayFrameHead {
	return &highwayFrameHead{magicID, msgID, 0, 0}
}

//HighWayProtocalObject highway Protocal
type HighWayProtocalObject struct {
	FrHead  highwayFrameHead
	payLoad []byte
}

//ProtocalName protocal name
func (msgObj *HighWayProtocalObject) ProtocalName() string {
	return "Highway"
}

//SerializeReq Serialize request
func (msgObj *HighWayProtocalObject) SerializeReq(req *HighwayRequest, wBuf *bufio.Writer) {
	frHead := newHeadFrame(uint64(req.MsgID))
	//flags:Indicates whether compression , temporarily not to use
	reqHeader := highway.RequestHeader{
		MsgType:          highway.MsgTypeRequest,
		Flags:            int32(0),
		DestMicroservice: req.SvcName,
		OperationName:    req.MethodName,
		SchemaID:         req.Schema,
		Context:          req.Attachments,
	}

	header, err := proto.Marshal(&reqHeader)
	if err != nil {
		lager.Logger.Errorf(err, "client marshal highway request header failed.")
		return
	}
	frHead.HeaderLen = uint32(len(header))
	body, err := proto.Marshal(req.Arg.(proto.Message))
	if err != nil {
		return
	}
	frHead.TotalLen = frHead.HeaderLen + uint32(len(body))
	wBuf.Write(frHead.serialize())
	wBuf.Write(header)
	wBuf.Write(body)
}

//SerializeRsp Serialize frame
func (msgObj *HighWayProtocalObject) SerializeRsp(rsp *HighwayRespond, wBuf *bufio.Writer) {
	frHead := newHeadFrame(uint64(rsp.MsgID))
	//todo parse meta
	//flags:Indicates whether compression , temporarily not to use
	respHeader := &highway.ResponseHeader{}
	respHeader.StatusCode = int32(rsp.Status)
	respHeader.Reason = rsp.Err
	rsp.Attachments = respHeader.Context
	respHeader.Flags = 0
	header, err := proto.Marshal(respHeader)
	if err != nil {
		lager.Logger.Errorf(err, "client marshal highway request header failed.")
		return
	}
	var body []byte
	frHead.HeaderLen = uint32(len(header))
	if rsp.Result != nil {
		body, err = proto.Marshal(rsp.Result.(proto.Message))
		if err != nil {
			return
		}
		frHead.TotalLen = frHead.HeaderLen + uint32(len(body))
	} else {
		frHead.TotalLen = frHead.HeaderLen
	}

	wBuf.Write(frHead.serialize())
	wBuf.Write(header)
	if body != nil {
		wBuf.Write(body)
	}
}

//DeSerializeFrame Deserialize frame
func (msgObj *HighWayProtocalObject) DeSerializeFrame(rdBuf *bufio.Reader) error {
	var err error
	var count int
	//Parse frame head
	buf := make([]byte, FrameHeadLen)
	count = 0
	for count < FrameHeadLen {
		tmpsize, rdErr := rdBuf.Read(buf[count:])
		if rdErr != nil {
			lager.Logger.Errorf(rdErr, "Recv Frame head  failed.")
			return rdErr
		}
		count += tmpsize
	}

	msgObj.FrHead = highwayFrameHead{}
	err = msgObj.FrHead.deserialize(buf)
	if err != nil {
		lager.Logger.Errorf(err, "Frame head error.")
		return err
	}
	msgObj.payLoad = make([]byte, msgObj.FrHead.TotalLen)

	count = 0
	for count < int(msgObj.FrHead.TotalLen) {
		tmpsize, rdErr := rdBuf.Read(msgObj.payLoad[count:])
		if rdErr != nil {
			lager.Logger.Errorf(rdErr, "Read frame body  failed")
			return rdErr
		}
		count += tmpsize
	}

	return nil
}

//DeSerializeRsp Deserialize rsp
func (msgObj *HighWayProtocalObject) DeSerializeRsp(rsp *HighwayRespond) error {
	var err error
	rsp.MsgID = msgObj.FrHead.MsgID
	respHeader := &highway.ResponseHeader{}
	//Head
	err = proto.Unmarshal(msgObj.payLoad[0:msgObj.FrHead.HeaderLen], respHeader)
	if err != nil {
		lager.Logger.Errorf(err, "Unmarshal response header failed")
		return err
	}
	rsp.Status = int(respHeader.GetStatusCode())
	rsp.Err = respHeader.GetReason()
	rsp.Attachments = respHeader.Context

	//Body
	if msgObj.FrHead.HeaderLen != msgObj.FrHead.TotalLen {
		err = proto.Unmarshal(msgObj.payLoad[msgObj.FrHead.HeaderLen:], (rsp.Result).(proto.Message))
		if err != nil {
			lager.Logger.Errorf(err, "Unmarshal response body  failed")
			rsp.Err = err.Error()
			return err
		}
	}
	return nil
}

//DeSerializeReq Deserialize req
func (msgObj *HighWayProtocalObject) DeSerializeReq(req *HighwayRequest) error {
	var err error
	req.MsgID = msgObj.FrHead.MsgID
	reqHeader := &highway.RequestHeader{}

	err = proto.Unmarshal(msgObj.payLoad[0:msgObj.FrHead.HeaderLen], reqHeader)
	if err != nil {
		lager.Logger.Errorf(err, "Unmarshal request header failed")
		return err
	}
	if req.Arg == nil {
		req.MethodName = reqHeader.GetOperationName()
		req.SvcName = reqHeader.GetDestMicroservice()
		req.Schema = reqHeader.GetSchemaID()
		req.Attachments = reqHeader.Context
		req.MsgType = int(reqHeader.MsgType)
		//Here we need to parse Attachments, indicating whether it is TwoWay,Current only twoway
		req.TwoWay = true
		var op provider.Operation
		op, err = provider.GetOperation(req.SvcName, req.Schema, req.MethodName)
		if err != nil {
			return err
		}
		if op != nil && op.Args() != nil && len(op.Args()) > 0 {
			if op.Args()[1].Kind() != reflect.Ptr {
				err = errors.New("second arg not ptr")
				return err
			}
			argv := reflect.New(op.Args()[1].Elem())
			req.Arg = argv.Interface()
			//Body
			err = proto.Unmarshal(msgObj.payLoad[msgObj.FrHead.HeaderLen:], (req.Arg).(proto.Message))
			if err != nil {
				lager.Logger.Errorf(err, "Unmarshal request body  failed")
				return err
			}
		}
	} else {
		err = proto.Unmarshal(msgObj.payLoad[msgObj.FrHead.HeaderLen:], (req.Arg).(proto.Message))
		if err != nil {
			lager.Logger.Errorf(err, "Unmarshal hello request body  failed")
			return err
		}
	}
	return nil
}

//SerializeHelloReq Serialize hello req
func (msgObj *HighWayProtocalObject) SerializeHelloReq(wBuf *bufio.Writer) error {
	frHead := newHeadFrame(GenerateMsgID())
	reqHeader := highway.RequestHeader{
		MsgType:          highway.MsgTypeLogin,
		Flags:            int32(0),
		DestMicroservice: "",
		OperationName:    "",
		SchemaID:         "",
		Context:          nil,
	}
	header, err := proto.Marshal(&reqHeader)
	if err != nil {
		lager.Logger.Errorf(err, "Marshal highway login header failed")
		return err
	}
	frHead.HeaderLen = uint32(len(header))

	loginBody := highway.LoginRequest{
		Protocol:            "highway",
		ZipName:             "z",
		UseProtobufMapCodec: localSupportLogin,
	}
	body, err := proto.Marshal(&loginBody)
	if err != nil {
		lager.Logger.Errorf(err, "Marshal highway login body failed")
		return err
	}
	frHead.TotalLen = uint32(len(body)) + frHead.HeaderLen
	wBuf.Write(frHead.serialize())
	wBuf.Write(header)
	wBuf.Write(body)

	return nil
}

//SerializelLoginRsp Serialize hello req
func (msgObj *HighWayProtocalObject) SerializelLoginRsp(msgID uint64, wBuf *bufio.Writer) error {
	frHead := newHeadFrame(msgID)
	reqHeader := &highway.ResponseHeader{
		Flags:      int32(0),
		StatusCode: Ok,
		Reason:     "",
		Context:    nil,
	}
	header, err := proto.Marshal(reqHeader)
	if err != nil {
		lager.Logger.Errorf(err, "Marshal highway login header failed")
		return err
	}

	frHead.HeaderLen = uint32(len(header))

	loginRspBody := &highway.LoginResponse{
		Protocol:            "highway",
		ZipName:             "z",
		UseProtobufMapCodec: true,
	}

	body, err := proto.Marshal(loginRspBody)
	if err != nil {
		lager.Logger.Errorf(err, "Marshal highway login body failed")
		return err
	}
	frHead.TotalLen = uint32(len(body)) + frHead.HeaderLen
	wBuf.Write(frHead.serialize())
	wBuf.Write(header)
	wBuf.Write(body)
	return nil
}
