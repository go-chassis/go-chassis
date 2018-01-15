package codec

import (
	jsonEnc "encoding/json"
)

//NewJSONCodec 创建JSON编解码对象实例
func NewJSONCodec() Codec {
	return new(jsonCodec)
}

//jsonCodec JSON编解码器
type jsonCodec struct {
}

//Marshal 编码函数.
func (json *jsonCodec) Marshal(v interface{}) (res []byte, err error) {
	return jsonEnc.Marshal(v)
}

//Unmarshal 解码函数.
func (json *jsonCodec) Unmarshal(data []byte, v interface{}) (err error) {
	return jsonEnc.Unmarshal(data, v)
}
