package codec_test

import (
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/third_party/forked/go-micro/codec"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/examples/helloworld/helloworld"
	"reflect"
	"testing"
)

type Data struct {
	Name string `json:"name"`
	Mny  int    `json:"mny"`
}

func TestJSONMarshal_Marshal(t *testing.T) {
	t.Log("testing json codec")
	var data Data
	data.Name = "中文"
	data.Mny = 123
	codec := codec.NewJSONCodec()
	jsonBytes, err := codec.Marshal(data)
	assert.NoError(t, err)
	if err != nil {
		t.Errorf("Unexpected Marshal err: %v", err)
	}

	var decodeJSON Data
	err1 := codec.Unmarshal(jsonBytes, &decodeJSON)
	assert.NoError(t, err1)
	if err1 != nil {
		t.Errorf("Unexpected Unmarshal err: %v", err1)
	}
	if (decodeJSON.Name != data.Name) || (decodeJSON.Mny != data.Mny) {
		t.Errorf("Unexpected Unmarshal")
	}
}

func TestProtobuffer_Marshal(t *testing.T) {
	t.Log("testing protobuf codec")
	req := &helloworld.HelloRequest{
		Name: "peter",
	}

	cc := codec.NewPBCodec()

	protobyte, err := cc.Marshal(req)
	assert.NoError(t, err)
	if err != nil {
		t.Errorf("Unexpected Marshal err: %v", err)
	}
	typ := reflect.ValueOf(req).Interface()

	err1 := cc.Unmarshal(protobyte, typ.(proto.Message))
	assert.NoError(t, err1)
	if err1 != nil {
		t.Errorf("Unexpected Unmarshal err: %v", err1)
	}
	if req != typ {
		t.Errorf("Unexpected Unmarshal")
	}
}

func TestInstallPlugin_GetCodecMap(t *testing.T) {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	c := codec.NewPBCodec()
	f := func() codec.Codec {
		return c
	}
	codec.InstallPlugin("abc", f)
	cmap := codec.GetCodecMap()
	assert.Equal(t, cmap["abc"], c)
}
