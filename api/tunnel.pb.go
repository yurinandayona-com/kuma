// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tunnel.proto

package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/ptypes/empty"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type RequestHeader struct {
	Method  string    `protobuf:"bytes,1,opt,name=method" json:"method,omitempty"`
	Path    string    `protobuf:"bytes,2,opt,name=path" json:"path,omitempty"`
	Headers []*Header `protobuf:"bytes,3,rep,name=headers" json:"headers,omitempty"`
}

func (m *RequestHeader) Reset()                    { *m = RequestHeader{} }
func (m *RequestHeader) String() string            { return proto.CompactTextString(m) }
func (*RequestHeader) ProtoMessage()               {}
func (*RequestHeader) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *RequestHeader) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *RequestHeader) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *RequestHeader) GetHeaders() []*Header {
	if m != nil {
		return m.Headers
	}
	return nil
}

type RequestBody struct {
	Body []byte `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *RequestBody) Reset()                    { *m = RequestBody{} }
func (m *RequestBody) String() string            { return proto.CompactTextString(m) }
func (*RequestBody) ProtoMessage()               {}
func (*RequestBody) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *RequestBody) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type ResponseHeader struct {
	Status  int32     `protobuf:"varint,1,opt,name=status" json:"status,omitempty"`
	Headers []*Header `protobuf:"bytes,2,rep,name=headers" json:"headers,omitempty"`
}

func (m *ResponseHeader) Reset()                    { *m = ResponseHeader{} }
func (m *ResponseHeader) String() string            { return proto.CompactTextString(m) }
func (*ResponseHeader) ProtoMessage()               {}
func (*ResponseHeader) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *ResponseHeader) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *ResponseHeader) GetHeaders() []*Header {
	if m != nil {
		return m.Headers
	}
	return nil
}

type ResponseBody struct {
	Body []byte `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (m *ResponseBody) Reset()                    { *m = ResponseBody{} }
func (m *ResponseBody) String() string            { return proto.CompactTextString(m) }
func (*ResponseBody) ProtoMessage()               {}
func (*ResponseBody) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{3} }

func (m *ResponseBody) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type ResponseError struct {
	Error string `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
}

func (m *ResponseError) Reset()                    { *m = ResponseError{} }
func (m *ResponseError) String() string            { return proto.CompactTextString(m) }
func (*ResponseError) ProtoMessage()               {}
func (*ResponseError) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{4} }

func (m *ResponseError) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type Header struct {
	Name   string   `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Values []string `protobuf:"bytes,2,rep,name=values" json:"values,omitempty"`
}

func (m *Header) Reset()                    { *m = Header{} }
func (m *Header) String() string            { return proto.CompactTextString(m) }
func (*Header) ProtoMessage()               {}
func (*Header) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{5} }

func (m *Header) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Header) GetValues() []string {
	if m != nil {
		return m.Values
	}
	return nil
}

func init() {
	proto.RegisterType((*RequestHeader)(nil), "kuma.tunnel.RequestHeader")
	proto.RegisterType((*RequestBody)(nil), "kuma.tunnel.RequestBody")
	proto.RegisterType((*ResponseHeader)(nil), "kuma.tunnel.ResponseHeader")
	proto.RegisterType((*ResponseBody)(nil), "kuma.tunnel.ResponseBody")
	proto.RegisterType((*ResponseError)(nil), "kuma.tunnel.ResponseError")
	proto.RegisterType((*Header)(nil), "kuma.tunnel.Header")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Tunnel service

type TunnelClient interface {
	ReceiveHeader(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*RequestHeader, error)
	ReceiveBody(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (Tunnel_ReceiveBodyClient, error)
	SendHeader(ctx context.Context, in *ResponseHeader, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
	SendBody(ctx context.Context, opts ...grpc.CallOption) (Tunnel_SendBodyClient, error)
	SendError(ctx context.Context, in *ResponseError, opts ...grpc.CallOption) (*google_protobuf.Empty, error)
}

type tunnelClient struct {
	cc *grpc.ClientConn
}

func NewTunnelClient(cc *grpc.ClientConn) TunnelClient {
	return &tunnelClient{cc}
}

func (c *tunnelClient) ReceiveHeader(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (*RequestHeader, error) {
	out := new(RequestHeader)
	err := grpc.Invoke(ctx, "/kuma.tunnel.Tunnel/ReceiveHeader", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tunnelClient) ReceiveBody(ctx context.Context, in *google_protobuf.Empty, opts ...grpc.CallOption) (Tunnel_ReceiveBodyClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Tunnel_serviceDesc.Streams[0], c.cc, "/kuma.tunnel.Tunnel/ReceiveBody", opts...)
	if err != nil {
		return nil, err
	}
	x := &tunnelReceiveBodyClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Tunnel_ReceiveBodyClient interface {
	Recv() (*RequestBody, error)
	grpc.ClientStream
}

type tunnelReceiveBodyClient struct {
	grpc.ClientStream
}

func (x *tunnelReceiveBodyClient) Recv() (*RequestBody, error) {
	m := new(RequestBody)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tunnelClient) SendHeader(ctx context.Context, in *ResponseHeader, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/kuma.tunnel.Tunnel/SendHeader", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *tunnelClient) SendBody(ctx context.Context, opts ...grpc.CallOption) (Tunnel_SendBodyClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Tunnel_serviceDesc.Streams[1], c.cc, "/kuma.tunnel.Tunnel/SendBody", opts...)
	if err != nil {
		return nil, err
	}
	x := &tunnelSendBodyClient{stream}
	return x, nil
}

type Tunnel_SendBodyClient interface {
	Send(*ResponseBody) error
	CloseAndRecv() (*google_protobuf.Empty, error)
	grpc.ClientStream
}

type tunnelSendBodyClient struct {
	grpc.ClientStream
}

func (x *tunnelSendBodyClient) Send(m *ResponseBody) error {
	return x.ClientStream.SendMsg(m)
}

func (x *tunnelSendBodyClient) CloseAndRecv() (*google_protobuf.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(google_protobuf.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tunnelClient) SendError(ctx context.Context, in *ResponseError, opts ...grpc.CallOption) (*google_protobuf.Empty, error) {
	out := new(google_protobuf.Empty)
	err := grpc.Invoke(ctx, "/kuma.tunnel.Tunnel/SendError", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Tunnel service

type TunnelServer interface {
	ReceiveHeader(context.Context, *google_protobuf.Empty) (*RequestHeader, error)
	ReceiveBody(*google_protobuf.Empty, Tunnel_ReceiveBodyServer) error
	SendHeader(context.Context, *ResponseHeader) (*google_protobuf.Empty, error)
	SendBody(Tunnel_SendBodyServer) error
	SendError(context.Context, *ResponseError) (*google_protobuf.Empty, error)
}

func RegisterTunnelServer(s *grpc.Server, srv TunnelServer) {
	s.RegisterService(&_Tunnel_serviceDesc, srv)
}

func _Tunnel_ReceiveHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(google_protobuf.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TunnelServer).ReceiveHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kuma.tunnel.Tunnel/ReceiveHeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TunnelServer).ReceiveHeader(ctx, req.(*google_protobuf.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tunnel_ReceiveBody_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(google_protobuf.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TunnelServer).ReceiveBody(m, &tunnelReceiveBodyServer{stream})
}

type Tunnel_ReceiveBodyServer interface {
	Send(*RequestBody) error
	grpc.ServerStream
}

type tunnelReceiveBodyServer struct {
	grpc.ServerStream
}

func (x *tunnelReceiveBodyServer) Send(m *RequestBody) error {
	return x.ServerStream.SendMsg(m)
}

func _Tunnel_SendHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResponseHeader)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TunnelServer).SendHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kuma.tunnel.Tunnel/SendHeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TunnelServer).SendHeader(ctx, req.(*ResponseHeader))
	}
	return interceptor(ctx, in, info, handler)
}

func _Tunnel_SendBody_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TunnelServer).SendBody(&tunnelSendBodyServer{stream})
}

type Tunnel_SendBodyServer interface {
	SendAndClose(*google_protobuf.Empty) error
	Recv() (*ResponseBody, error)
	grpc.ServerStream
}

type tunnelSendBodyServer struct {
	grpc.ServerStream
}

func (x *tunnelSendBodyServer) SendAndClose(m *google_protobuf.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *tunnelSendBodyServer) Recv() (*ResponseBody, error) {
	m := new(ResponseBody)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Tunnel_SendError_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResponseError)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TunnelServer).SendError(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kuma.tunnel.Tunnel/SendError",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TunnelServer).SendError(ctx, req.(*ResponseError))
	}
	return interceptor(ctx, in, info, handler)
}

var _Tunnel_serviceDesc = grpc.ServiceDesc{
	ServiceName: "kuma.tunnel.Tunnel",
	HandlerType: (*TunnelServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReceiveHeader",
			Handler:    _Tunnel_ReceiveHeader_Handler,
		},
		{
			MethodName: "SendHeader",
			Handler:    _Tunnel_SendHeader_Handler,
		},
		{
			MethodName: "SendError",
			Handler:    _Tunnel_SendError_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ReceiveBody",
			Handler:       _Tunnel_ReceiveBody_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SendBody",
			Handler:       _Tunnel_SendBody_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "tunnel.proto",
}

func init() { proto.RegisterFile("tunnel.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 355 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0x4f, 0x4f, 0xf2, 0x40,
	0x10, 0xc6, 0x53, 0xfe, 0xf4, 0x7d, 0x19, 0xd0, 0xc3, 0x6a, 0x4c, 0x2d, 0x17, 0x6c, 0x62, 0xc2,
	0xc5, 0xc5, 0xa0, 0x77, 0x02, 0x86, 0xc4, 0xf3, 0x6a, 0x62, 0xe2, 0x6d, 0xb1, 0x23, 0xa0, 0xb4,
	0x5b, 0xdb, 0x2d, 0x09, 0x1f, 0xd4, 0xef, 0x63, 0x76, 0x76, 0x51, 0x30, 0x34, 0xf1, 0x36, 0xb3,
	0x33, 0xfb, 0x7b, 0x9e, 0x99, 0x81, 0x8e, 0x2e, 0xd3, 0x14, 0x57, 0x3c, 0xcb, 0x95, 0x56, 0xac,
	0xfd, 0x5e, 0x26, 0x92, 0xdb, 0xa7, 0xb0, 0x3b, 0x57, 0x6a, 0xbe, 0xc2, 0x01, 0x95, 0x66, 0xe5,
	0xeb, 0x00, 0x93, 0x4c, 0x6f, 0x6c, 0x67, 0xf4, 0x06, 0x47, 0x02, 0x3f, 0x4a, 0x2c, 0xf4, 0x3d,
	0xca, 0x18, 0x73, 0x76, 0x06, 0x7e, 0x82, 0x7a, 0xa1, 0xe2, 0xc0, 0xeb, 0x79, 0xfd, 0x96, 0x70,
	0x19, 0x63, 0xd0, 0xc8, 0xa4, 0x5e, 0x04, 0x35, 0x7a, 0xa5, 0x98, 0x5d, 0xc1, 0xbf, 0x05, 0xfd,
	0x2a, 0x82, 0x7a, 0xaf, 0xde, 0x6f, 0x0f, 0x4f, 0xf8, 0x8e, 0x30, 0xb7, 0x44, 0xb1, 0xed, 0x89,
	0x2e, 0xa0, 0xed, 0xb4, 0x26, 0x2a, 0xde, 0x18, 0xe2, 0x4c, 0xc5, 0x1b, 0xd2, 0xe9, 0x08, 0x8a,
	0xa3, 0x27, 0x38, 0x16, 0x58, 0x64, 0x2a, 0x2d, 0xf0, 0xc7, 0x4f, 0xa1, 0xa5, 0x2e, 0x0b, 0xea,
	0x6b, 0x0a, 0x97, 0xed, 0x6a, 0xd7, 0xfe, 0xa0, 0x1d, 0x41, 0x67, 0x0b, 0xae, 0x14, 0xbf, 0x34,
	0xbb, 0xb0, 0x3d, 0xd3, 0x3c, 0x57, 0x39, 0x3b, 0x85, 0x26, 0x9a, 0xc0, 0xad, 0xc2, 0x26, 0xd1,
	0x2d, 0xf8, 0xce, 0x1b, 0x83, 0x46, 0x2a, 0x13, 0x74, 0x65, 0x8a, 0x8d, 0xdf, 0xb5, 0x5c, 0x95,
	0x68, 0x6d, 0xb5, 0x84, 0xcb, 0x86, 0x9f, 0x35, 0xf0, 0x1f, 0xc9, 0x1b, 0xbb, 0x33, 0x3a, 0x2f,
	0xb8, 0x5c, 0x7f, 0xcf, 0xc8, 0xed, 0x89, 0xf8, 0xf6, 0x44, 0x7c, 0x6a, 0x4e, 0x14, 0x86, 0x7b,
	0x23, 0xed, 0xdf, 0x69, 0x6c, 0x96, 0x49, 0x10, 0x9a, 0xa7, 0x0a, 0x11, 0x1c, 0x42, 0x98, 0x1f,
	0xd7, 0x1e, 0x1b, 0x03, 0x3c, 0x60, 0x1a, 0x3b, 0x60, 0xf7, 0x57, 0xe7, 0xee, 0x15, 0xc2, 0x0a,
	0x3c, 0x1b, 0xc1, 0x7f, 0x83, 0x20, 0x0b, 0xe7, 0x07, 0x01, 0xa6, 0x54, 0xf5, 0xbd, 0xef, 0xb1,
	0x11, 0xb4, 0x0c, 0xc0, 0xee, 0x3b, 0x3c, 0x48, 0xa0, 0x5a, 0x15, 0x62, 0xd2, 0x7c, 0xae, 0xcb,
	0x6c, 0x39, 0xf3, 0xe9, 0xf9, 0xe6, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x69, 0x22, 0xbc, 0x02, 0x08,
	0x03, 0x00, 0x00,
}
