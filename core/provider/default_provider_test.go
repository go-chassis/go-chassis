// Forked from github.com/golang/go
// Some parts of this file have been modified to make it functional in this package
package provider_test

import (
	"context"
	"github.com/go-chassis/go-archaius"
	"github.com/go-chassis/go-chassis/v2/core/config/model"
	"testing"

	"github.com/go-chassis/go-chassis/v2/core/config"
	"github.com/go-chassis/go-chassis/v2/core/invocation"
	"github.com/go-chassis/go-chassis/v2/core/lager"
	"github.com/go-chassis/go-chassis/v2/core/provider"

	"github.com/stretchr/testify/assert"
)

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the user's name.
type HelloRequest struct {
	Name string
}

func (m *HelloRequest) Reset()                    { *m = HelloRequest{} }
func (m *HelloRequest) String() string            { return proto.CompactTextString(m) }
func (*HelloRequest) ProtoMessage()               {}
func (*HelloRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *HelloRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

// The response message containing the greetings
type HelloReply struct {
	Message string
}

func (m *HelloReply) Reset()                    { *m = HelloReply{} }
func (m *HelloReply) String() string            { return proto.CompactTextString(m) }
func (*HelloReply) ProtoMessage()               {}
func (*HelloReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *HelloReply) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

var fileDescriptor0 = []byte{
	// 142 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xca, 0x48, 0xcd, 0xc9, 0xc9, 0x2f, 0xcf, 0x2f,
	0xca, 0x49, 0x51, 0x52, 0xe2, 0xe2, 0xf1, 0x00, 0xf1, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b,
	0x84, 0x84, 0xb8, 0x58, 0xf2, 0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0xc0,
	0x6c, 0x25, 0x35, 0x2e, 0x2e, 0xa8, 0x9a, 0x82, 0x9c, 0x4a, 0x21, 0x09, 0x2e, 0xf6, 0xdc, 0xd4,
	0xe2, 0xe2, 0xc4, 0x74, 0x98, 0x22, 0x18, 0xd7, 0xc9, 0x80, 0x4b, 0x3a, 0x33, 0x5f, 0x2f, 0xbd,
	0xa8, 0x20, 0x59, 0x2f, 0xb5, 0x22, 0x31, 0xb7, 0x20, 0x27, 0xb5, 0x58, 0x0f, 0x61, 0x95, 0x13,
	0x3f, 0xd8, 0x90, 0x70, 0x10, 0x3b, 0x00, 0xe4, 0x8e, 0x00, 0xc6, 0x24, 0x36, 0xb0, 0x83, 0x8c,
	0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x27, 0xd2, 0x51, 0x53, 0x9e, 0x00, 0x00, 0x00,
}

type HelloServerfk1 struct{}
type HelloServerfk2 struct{}

func (s *HelloServerfk1) sayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Go Hello  " + in.Name}, nil
}

func TestRegister(t *testing.T) {
	p := &provider.DefaultProvider{}
	err := p.RegisterName("schema1", &HelloServer{})
	assert.NoError(t, err)
}

func TestRegister2(t *testing.T) {

	provider := &provider.DefaultProvider{}
	sname, err := provider.Register(&HelloServer{})
	assert.NoError(t, err)
	assert.NotEqual(t, "", sname)
	assert.NotEqual(t, nil, sname)
	op, _ := provider.GetOperation(sname, "SayHello")
	assert.NotNil(t, op)
	method := op.Method()
	assert.NotEqual(t, nil, method)
	args := op.Args()
	assert.Equal(t, 2, len(args))
	replay := op.Reply()
	assert.Equal(t, 2, len(replay))
	schema1 := "fkschema"
	err = provider.RegisterName(schema1, &HelloServerfk1{})
	assert.Error(t, err)
	schema2 := "fkschema"
	err = provider.RegisterName(schema2, &HelloServerfk2{})
	assert.Error(t, err)

}

func TestProvider_Invoke(t *testing.T) {
	p := &provider.DefaultProvider{}
	schema := "schema1"
	err := p.RegisterName(schema, &HelloServer{})
	assert.NoError(t, err)
	inv := &invocation.Invocation{
		SchemaID:    schema,
		OperationID: "SayHello",
		Args: &HelloRequest{
			Name: "1",
		},
	}
	_, err = p.Invoke(inv)
	assert.NoError(t, err)

}

//HelloServer is a struct
type HelloServer struct {
}

//SayHello is a method used to reply message
func (s *HelloServer) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{Message: "Go Hello  " + in.Name}, nil
}

func TestDefaultProvider_GetOperation(t *testing.T) {
	p := &provider.DefaultProvider{}
	schema := "schema1"
	err := p.RegisterName(schema, &HelloServer{})
	assert.NoError(t, err)
	op, _ := p.GetOperation(schema, "SayHello")
	assert.NotNil(t, op)
	op, _ = p.GetOperation(schema, "SayHelloasd")
	assert.Nil(t, op)
}

func TestDefaultProvider_Exist(t *testing.T) {
	p := &provider.DefaultProvider{}
	schema := "schema1"
	err := p.RegisterName(schema, &HelloServer{})
	assert.NoError(t, err)
	e := p.Exist(schema, "SayHello")
	assert.Equal(t, true, e)
	e = p.Exist(schema, "SayHelloaasda")
	assert.Equal(t, false, e)
	e = p.Exist(schema+"123", "SayHelloaasda")
	assert.Equal(t, false, e)
}

func BenchmarkDefaultProvider_GetOperation(b *testing.B) {
	p := &provider.DefaultProvider{}
	schema := "schema1"
	_ = p.RegisterName(schema, &HelloServer{})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.GetOperation(schema, "SayHello")
	}
}
func init() {
	archaius.Init(archaius.WithMemorySource())
	archaius.Set("servicecomb.noRefreshSchema", true)
	archaius.Set("servicecomb.service.name", "Client")
	archaius.Set("servicecomb.service.hostname", "localhost")
	config.MicroserviceDefinition = &model.ServiceSpec{}
	archaius.UnmarshalConfig(config.MicroserviceDefinition)
	config.ReadGlobalConfigFromArchaius()
	lager.Init(&lager.Options{
		LoggerLevel: "INFO",
	})
}
