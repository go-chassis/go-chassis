// Package codec is an interface for encoding messages
package codec

import (
	"io"

	"github.com/ServiceComb/go-chassis/core/lager"
)

// MessageType is gives the info about message type
type MessageType int

// NewClientCodec takes in a connection/buffer and returns a new client codec
type NewClientCodec func(io.ReadWriteCloser) ClientCodec

// NewServerCodec takes in a connection/buffer and returns a new server codec
type NewServerCodec func(io.ReadWriteCloser) ServerCodec

// ClientCodec writes RPC requests and reads RPC responses
// in the client side of an RPC session.
// ReadResponseHeader and ReadResponseBody are called in pairs
// to read requests.
// WriteRequest writes a request to the connection
// ReadResponseBody could be called with a nil param to
// force the body of the response to be read and then discarded.
type ClientCodec interface {
	ReadResponseHeader(*Response) error
	ReadResponseBody(interface{}) error
	// WriteRequest must be safe for concurrent use by multiple goroutines.
	// don't return bytes
	WriteRequest(*Request) error
	Close() error
}

// ServerCodec reads RPC requests and writes RPC responses
// in the server side of an RPC session.
// ReadRequestHeader and ReadRequestBody are called in pairs
// to read requests from the connection.
// WriteResponse writes a response back.
// ReadRequestBody could be called with a nil param to
// force the body of the request to be read and discarded.
type ServerCodec interface {
	ReadRequestHeader(*Request) error
	ReadRequestBody(interface{}) error
	// WriteResponse must be safe for concurrent use by multiple goroutines.
	WriteResponse(r *Response) error
	Close() error
}

// Request Message represents detailed information about the communication, likely followed by the body.
// In the case of an error, body may be nil
type Request struct {
	ID          uint64
	SchemaID    string
	OperationID string
	Error       string
	Header      map[string]string
	Arg         interface{}
}

// Response Message represents detailed information about the communication, likely followed by the body.
// In the case of an error, body may be nil
type Response struct {
	ID          uint64
	SchemaID    string
	OperationID string
	Error       string
	Header      map[string]string
	reply       interface{}
}

//Codecs 编解码对象表 key为编解码名称 value为编解码对象接口的实现
var codecs = map[string]func() Codec{}

//Codec 编解码接口
type Codec interface {
	// 编码函数.
	Marshal(v interface{}) ([]byte, error)
	// 解码函数.
	Unmarshal(data []byte, v interface{}) error
}

// InstallPlugin to install the codec plugins
func InstallPlugin(t string, f func() Codec) {
	codecs[t] = f
	lager.Logger.Debugf("Install Codec Plugin,codec_name:%s", t)
}

// GetCodecMap is to get the codec map
func GetCodecMap() map[string]Codec {
	cm := make(map[string]Codec)
	for k, v := range codecs {
		cm[k] = v()
	}
	return cm
}

// init is to initialize the codec functions
func init() {
	codecs["application/json"] = NewJSONCodec
	codecs["application/protobuf"] = NewPBCodec
}
