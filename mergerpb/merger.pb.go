// Code generated by protoc-gen-go. DO NOT EDIT.
// source: merger.proto

/*
Package merger is a generated protocol buffer package.

It is generated from these files:
	merger.proto

It has these top-level messages:
	EventRequest
	DataWrapper
	Trace
	DataMeta
	Empty
*/
package merger

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EventRequest struct {
	Payload []*DataWrapper `protobuf:"bytes,1,rep,name=payload" json:"payload,omitempty"`
}

func (m *EventRequest) Reset()                    { *m = EventRequest{} }
func (m *EventRequest) String() string            { return proto.CompactTextString(m) }
func (*EventRequest) ProtoMessage()               {}
func (*EventRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EventRequest) GetPayload() []*DataWrapper {
	if m != nil {
		return m.Payload
	}
	return nil
}

type DataWrapper struct {
	Trace *Trace    `protobuf:"bytes,1,opt,name=trace" json:"trace,omitempty"`
	Meta  *DataMeta `protobuf:"bytes,3,opt,name=meta" json:"meta,omitempty"`
	Data  []byte    `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *DataWrapper) Reset()                    { *m = DataWrapper{} }
func (m *DataWrapper) String() string            { return proto.CompactTextString(m) }
func (*DataWrapper) ProtoMessage()               {}
func (*DataWrapper) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DataWrapper) GetTrace() *Trace {
	if m != nil {
		return m.Trace
	}
	return nil
}

func (m *DataWrapper) GetMeta() *DataMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

func (m *DataWrapper) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Trace struct {
	Id     string   `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Traces []*Trace `protobuf:"bytes,2,rep,name=traces" json:"traces,omitempty"`
}

func (m *Trace) Reset()                    { *m = Trace{} }
func (m *Trace) String() string            { return proto.CompactTextString(m) }
func (*Trace) ProtoMessage()               {}
func (*Trace) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *Trace) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Trace) GetTraces() []*Trace {
	if m != nil {
		return m.Traces
	}
	return nil
}

type DataMeta struct {
	Schema  string `protobuf:"bytes,1,opt,name=schema" json:"schema,omitempty"`
	Version string `protobuf:"bytes,2,opt,name=version" json:"version,omitempty"`
}

func (m *DataMeta) Reset()                    { *m = DataMeta{} }
func (m *DataMeta) String() string            { return proto.CompactTextString(m) }
func (*DataMeta) ProtoMessage()               {}
func (*DataMeta) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DataMeta) GetSchema() string {
	if m != nil {
		return m.Schema
	}
	return ""
}

func (m *DataMeta) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type Empty struct {
}

func (m *Empty) Reset()                    { *m = Empty{} }
func (m *Empty) String() string            { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()               {}
func (*Empty) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func init() {
	proto.RegisterType((*EventRequest)(nil), "merger.EventRequest")
	proto.RegisterType((*DataWrapper)(nil), "merger.DataWrapper")
	proto.RegisterType((*Trace)(nil), "merger.Trace")
	proto.RegisterType((*DataMeta)(nil), "merger.DataMeta")
	proto.RegisterType((*Empty)(nil), "merger.Empty")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for MergeService service

type MergeServiceClient interface {
	PartialEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*Empty, error)
	CompleteEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*Empty, error)
}

type mergeServiceClient struct {
	cc *grpc.ClientConn
}

func NewMergeServiceClient(cc *grpc.ClientConn) MergeServiceClient {
	return &mergeServiceClient{cc}
}

func (c *mergeServiceClient) PartialEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/merger.MergeService/PartialEvent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mergeServiceClient) CompleteEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := grpc.Invoke(ctx, "/merger.MergeService/CompleteEvent", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for MergeService service

type MergeServiceServer interface {
	PartialEvent(context.Context, *EventRequest) (*Empty, error)
	CompleteEvent(context.Context, *EventRequest) (*Empty, error)
}

func RegisterMergeServiceServer(s *grpc.Server, srv MergeServiceServer) {
	s.RegisterService(&_MergeService_serviceDesc, srv)
}

func _MergeService_PartialEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MergeServiceServer).PartialEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/merger.MergeService/PartialEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MergeServiceServer).PartialEvent(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MergeService_CompleteEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MergeServiceServer).CompleteEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/merger.MergeService/CompleteEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MergeServiceServer).CompleteEvent(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _MergeService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "merger.MergeService",
	HandlerType: (*MergeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PartialEvent",
			Handler:    _MergeService_PartialEvent_Handler,
		},
		{
			MethodName: "CompleteEvent",
			Handler:    _MergeService_CompleteEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "merger.proto",
}

func init() { proto.RegisterFile("merger.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xc1, 0x4b, 0xf3, 0x40,
	0x10, 0xc5, 0xbf, 0xa4, 0x6d, 0xf2, 0x39, 0x49, 0x45, 0x46, 0x91, 0xc5, 0x53, 0x58, 0x15, 0x72,
	0xb1, 0x87, 0x88, 0x9e, 0xd4, 0x8b, 0xf6, 0x58, 0x90, 0x55, 0xf0, 0x3c, 0x26, 0x83, 0x06, 0x92,
	0x66, 0xdd, 0xac, 0x81, 0x1e, 0xfc, 0xdf, 0xa5, 0xdb, 0x44, 0x82, 0x37, 0x6f, 0x3b, 0xf3, 0xe6,
	0xfd, 0xde, 0x30, 0x0b, 0x71, 0xcd, 0xe6, 0x8d, 0xcd, 0x42, 0x9b, 0xc6, 0x36, 0x18, 0xec, 0x2a,
	0x79, 0x0b, 0xf1, 0xb2, 0xe3, 0xb5, 0x55, 0xfc, 0xf1, 0xc9, 0xad, 0xc5, 0x0b, 0x08, 0x35, 0x6d,
	0xaa, 0x86, 0x0a, 0xe1, 0x25, 0x93, 0x34, 0xca, 0x0e, 0x17, 0xbd, 0xef, 0x81, 0x2c, 0xbd, 0x18,
	0xd2, 0x9a, 0x8d, 0x1a, 0x66, 0x64, 0x05, 0xd1, 0xa8, 0x8f, 0xa7, 0x30, 0xb3, 0x86, 0x72, 0x16,
	0x5e, 0xe2, 0xa5, 0x51, 0x36, 0x1f, 0xbc, 0xcf, 0xdb, 0xa6, 0xda, 0x69, 0x78, 0x06, 0xd3, 0x9a,
	0x2d, 0x89, 0x89, 0x9b, 0x39, 0x18, 0xf3, 0x57, 0x6c, 0x49, 0x39, 0x15, 0x11, 0xa6, 0x05, 0x59,
	0x12, 0xd3, 0xc4, 0x4b, 0x63, 0xe5, 0xde, 0xf2, 0x0e, 0x66, 0x8e, 0x84, 0xfb, 0xe0, 0x97, 0x85,
	0x0b, 0xd9, 0x53, 0x7e, 0x59, 0xe0, 0x39, 0x04, 0x8e, 0xdd, 0x0a, 0xdf, 0x2d, 0xfd, 0x2b, 0xb8,
	0x17, 0xe5, 0x0d, 0xfc, 0x1f, 0x52, 0xf0, 0x18, 0x82, 0x36, 0x7f, 0xe7, 0x9a, 0x7a, 0x4c, 0x5f,
	0xa1, 0x80, 0xb0, 0x63, 0xd3, 0x96, 0xcd, 0x5a, 0xf8, 0x4e, 0x18, 0x4a, 0x19, 0xc2, 0x6c, 0x59,
	0x6b, 0xbb, 0xc9, 0xbe, 0x20, 0x5e, 0x6d, 0xf1, 0x4f, 0x6c, 0xba, 0x32, 0x67, 0xbc, 0x82, 0xf8,
	0x91, 0x8c, 0x2d, 0xa9, 0x72, 0xa7, 0xc4, 0xa3, 0x21, 0x7d, 0x7c, 0xd9, 0x93, 0x9f, 0x9d, 0x1c,
	0x44, 0xfe, 0xc3, 0x6b, 0x98, 0xdf, 0x37, 0xb5, 0xae, 0xd8, 0xf2, 0x5f, 0x7c, 0xaf, 0x81, 0xfb,
	0xc1, 0xcb, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x71, 0xf1, 0x8a, 0x14, 0xd1, 0x01, 0x00, 0x00,
}
