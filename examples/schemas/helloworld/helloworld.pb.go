// Code generated by protoc-gen-go. DO NOT EDIT.
// source: helloworld.proto

/*
Package helloworld is a generated protocol buffer package.

It is generated from these files:
	helloworld.proto

It has these top-level messages:
	HelloRequest
	HelloReply
*/
package helloworld

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
	Name string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
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
	Message string `protobuf:"bytes,1,opt,name=message" json:"message,omitempty"`
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

func init() {
	proto.RegisterType((*HelloRequest)(nil), "helloworld.HelloRequest")
	proto.RegisterType((*HelloReply)(nil), "helloworld.HelloReply")
}

func init() { proto.RegisterFile("helloworld.proto", fileDescriptor0) }

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
