// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: feedpb/subscriber.proto

package feedpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Subscriber_Create_FullMethodName = "/feedpb.Subscriber/Create"
	Subscriber_Update_FullMethodName = "/feedpb.Subscriber/Update"
)

// SubscriberClient is the client API for Subscriber service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SubscriberClient interface {
	Create(ctx context.Context, in *CreateSubscriberRequest, opts ...grpc.CallOption) (*CreateSubscriberResponse, error)
	Update(ctx context.Context, in *UpdateSubscriberRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type subscriberClient struct {
	cc grpc.ClientConnInterface
}

func NewSubscriberClient(cc grpc.ClientConnInterface) SubscriberClient {
	return &subscriberClient{cc}
}

func (c *subscriberClient) Create(ctx context.Context, in *CreateSubscriberRequest, opts ...grpc.CallOption) (*CreateSubscriberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateSubscriberResponse)
	err := c.cc.Invoke(ctx, Subscriber_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *subscriberClient) Update(ctx context.Context, in *UpdateSubscriberRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Subscriber_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SubscriberServer is the server API for Subscriber service.
// All implementations must embed UnimplementedSubscriberServer
// for forward compatibility.
type SubscriberServer interface {
	Create(context.Context, *CreateSubscriberRequest) (*CreateSubscriberResponse, error)
	Update(context.Context, *UpdateSubscriberRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedSubscriberServer()
}

// UnimplementedSubscriberServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSubscriberServer struct{}

func (UnimplementedSubscriberServer) Create(context.Context, *CreateSubscriberRequest) (*CreateSubscriberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedSubscriberServer) Update(context.Context, *UpdateSubscriberRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedSubscriberServer) mustEmbedUnimplementedSubscriberServer() {}
func (UnimplementedSubscriberServer) testEmbeddedByValue()                    {}

// UnsafeSubscriberServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SubscriberServer will
// result in compilation errors.
type UnsafeSubscriberServer interface {
	mustEmbedUnimplementedSubscriberServer()
}

func RegisterSubscriberServer(s grpc.ServiceRegistrar, srv SubscriberServer) {
	// If the following call pancis, it indicates UnimplementedSubscriberServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Subscriber_ServiceDesc, srv)
}

func _Subscriber_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSubscriberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriberServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Subscriber_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriberServer).Create(ctx, req.(*CreateSubscriberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Subscriber_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateSubscriberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SubscriberServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Subscriber_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SubscriberServer).Update(ctx, req.(*UpdateSubscriberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Subscriber_ServiceDesc is the grpc.ServiceDesc for Subscriber service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Subscriber_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "feedpb.Subscriber",
	HandlerType: (*SubscriberServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Subscriber_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Subscriber_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "feedpb/subscriber.proto",
}
