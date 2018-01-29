package codec

import (
	"github.com/golang/protobuf/proto"
)

//NewPBCodec 创建PB编解码对象实例
func NewPBCodec() Codec {
	return new(pbCodec)
}

//pbCodec PB编解码器
type pbCodec struct {
}

//Marshal 编码函数.
func (pb *pbCodec) Marshal(v interface{}) (res []byte, err error) {
	result, _ := v.(proto.Message)
	return proto.Marshal(result)
}

//Unmarshal 解码函数.
func (pb *pbCodec) Unmarshal(data []byte, v interface{}) (err error) {
	result, _ := v.(proto.Message)
	return proto.Unmarshal(data, result)
}
