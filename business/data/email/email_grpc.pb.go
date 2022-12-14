// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: business/data/email/email.proto

package email

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EmailClient is the client API for Email service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmailClient interface {
	Send(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error)
}

type emailClient struct {
	cc grpc.ClientConnInterface
}

func NewEmailClient(cc grpc.ClientConnInterface) EmailClient {
	return &emailClient{cc}
}

func (c *emailClient) Send(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, "/email.Email/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmailServer is the server API for Email service.
// All implementations must embed UnimplementedEmailServer
// for forward compatibility
type EmailServer interface {
	Send(context.Context, *EmailRequest) (*EmailResponse, error)
	mustEmbedUnimplementedEmailServer()
}

// UnimplementedEmailServer must be embedded to have forward compatible implementations.
type UnimplementedEmailServer struct {
}

func (UnimplementedEmailServer) Send(context.Context, *EmailRequest) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedEmailServer) mustEmbedUnimplementedEmailServer() {}

// UnsafeEmailServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmailServer will
// result in compilation errors.
type UnsafeEmailServer interface {
	mustEmbedUnimplementedEmailServer()
}

func RegisterEmailServer(s grpc.ServiceRegistrar, srv EmailServer) {
	s.RegisterService(&Email_ServiceDesc, srv)
}

func _Email_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/email.Email/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServer).Send(ctx, req.(*EmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Email_ServiceDesc is the grpc.ServiceDesc for Email service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Email_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "email.Email",
	HandlerType: (*EmailServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _Email_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "business/data/email/email.proto",
}
