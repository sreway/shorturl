// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: proto/shorturl/v1/shorturl.proto

package v1

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

const (
	ShortURLService_CreateURL_FullMethodName    = "/shorturl.ShortURLService/CreateURL"
	ShortURLService_BatchURL_FullMethodName     = "/shorturl.ShortURLService/BatchURL"
	ShortURLService_GetURL_FullMethodName       = "/shorturl.ShortURLService/GetURL"
	ShortURLService_GetUserURLs_FullMethodName  = "/shorturl.ShortURLService/GetUserURLs"
	ShortURLService_DeleteURL_FullMethodName    = "/shorturl.ShortURLService/DeleteURL"
	ShortURLService_StorageCheck_FullMethodName = "/shorturl.ShortURLService/StorageCheck"
)

// ShortURLServiceClient is the client API for ShortURLService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ShortURLServiceClient interface {
	CreateURL(ctx context.Context, in *AddURLRequest, opts ...grpc.CallOption) (*AddURLResponse, error)
	BatchURL(ctx context.Context, in *BatchAddURLRequest, opts ...grpc.CallOption) (*BatchAddURLResponse, error)
	GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error)
	GetUserURLs(ctx context.Context, in *GetUserURLRequest, opts ...grpc.CallOption) (*GetUserURLResponse, error)
	DeleteURL(ctx context.Context, in *DeleteURLRequest, opts ...grpc.CallOption) (*DeleteURLResponse, error)
	StorageCheck(ctx context.Context, in *StorageCheckRequest, opts ...grpc.CallOption) (*StorageCheckResponse, error)
}

type shortURLServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewShortURLServiceClient(cc grpc.ClientConnInterface) ShortURLServiceClient {
	return &shortURLServiceClient{cc}
}

func (c *shortURLServiceClient) CreateURL(ctx context.Context, in *AddURLRequest, opts ...grpc.CallOption) (*AddURLResponse, error) {
	out := new(AddURLResponse)
	err := c.cc.Invoke(ctx, ShortURLService_CreateURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLServiceClient) BatchURL(ctx context.Context, in *BatchAddURLRequest, opts ...grpc.CallOption) (*BatchAddURLResponse, error) {
	out := new(BatchAddURLResponse)
	err := c.cc.Invoke(ctx, ShortURLService_BatchURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLServiceClient) GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error) {
	out := new(GetURLResponse)
	err := c.cc.Invoke(ctx, ShortURLService_GetURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLServiceClient) GetUserURLs(ctx context.Context, in *GetUserURLRequest, opts ...grpc.CallOption) (*GetUserURLResponse, error) {
	out := new(GetUserURLResponse)
	err := c.cc.Invoke(ctx, ShortURLService_GetUserURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLServiceClient) DeleteURL(ctx context.Context, in *DeleteURLRequest, opts ...grpc.CallOption) (*DeleteURLResponse, error) {
	out := new(DeleteURLResponse)
	err := c.cc.Invoke(ctx, ShortURLService_DeleteURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shortURLServiceClient) StorageCheck(ctx context.Context, in *StorageCheckRequest, opts ...grpc.CallOption) (*StorageCheckResponse, error) {
	out := new(StorageCheckResponse)
	err := c.cc.Invoke(ctx, ShortURLService_StorageCheck_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ShortURLServiceServer is the server API for ShortURLService service.
// All implementations must embed UnimplementedShortURLServiceServer
// for forward compatibility
type ShortURLServiceServer interface {
	CreateURL(context.Context, *AddURLRequest) (*AddURLResponse, error)
	BatchURL(context.Context, *BatchAddURLRequest) (*BatchAddURLResponse, error)
	GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error)
	GetUserURLs(context.Context, *GetUserURLRequest) (*GetUserURLResponse, error)
	DeleteURL(context.Context, *DeleteURLRequest) (*DeleteURLResponse, error)
	StorageCheck(context.Context, *StorageCheckRequest) (*StorageCheckResponse, error)
	mustEmbedUnimplementedShortURLServiceServer()
}

// UnimplementedShortURLServiceServer must be embedded to have forward compatible implementations.
type UnimplementedShortURLServiceServer struct {
}

func (UnimplementedShortURLServiceServer) CreateURL(context.Context, *AddURLRequest) (*AddURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateURL not implemented")
}
func (UnimplementedShortURLServiceServer) BatchURL(context.Context, *BatchAddURLRequest) (*BatchAddURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchURL not implemented")
}
func (UnimplementedShortURLServiceServer) GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURL not implemented")
}
func (UnimplementedShortURLServiceServer) GetUserURLs(context.Context, *GetUserURLRequest) (*GetUserURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserURLs not implemented")
}
func (UnimplementedShortURLServiceServer) DeleteURL(context.Context, *DeleteURLRequest) (*DeleteURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURL not implemented")
}
func (UnimplementedShortURLServiceServer) StorageCheck(context.Context, *StorageCheckRequest) (*StorageCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StorageCheck not implemented")
}
func (UnimplementedShortURLServiceServer) mustEmbedUnimplementedShortURLServiceServer() {}

// UnsafeShortURLServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ShortURLServiceServer will
// result in compilation errors.
type UnsafeShortURLServiceServer interface {
	mustEmbedUnimplementedShortURLServiceServer()
}

func RegisterShortURLServiceServer(s grpc.ServiceRegistrar, srv ShortURLServiceServer) {
	s.RegisterService(&ShortURLService_ServiceDesc, srv)
}

func _ShortURLService_CreateURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServiceServer).CreateURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortURLService_CreateURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServiceServer).CreateURL(ctx, req.(*AddURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURLService_BatchURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchAddURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServiceServer).BatchURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortURLService_BatchURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServiceServer).BatchURL(ctx, req.(*BatchAddURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURLService_GetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServiceServer).GetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortURLService_GetURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServiceServer).GetURL(ctx, req.(*GetURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURLService_GetUserURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServiceServer).GetUserURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortURLService_GetUserURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServiceServer).GetUserURLs(ctx, req.(*GetUserURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURLService_DeleteURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServiceServer).DeleteURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortURLService_DeleteURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServiceServer).DeleteURL(ctx, req.(*DeleteURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ShortURLService_StorageCheck_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StorageCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ShortURLServiceServer).StorageCheck(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ShortURLService_StorageCheck_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ShortURLServiceServer).StorageCheck(ctx, req.(*StorageCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ShortURLService_ServiceDesc is the grpc.ServiceDesc for ShortURLService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ShortURLService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shorturl.ShortURLService",
	HandlerType: (*ShortURLServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateURL",
			Handler:    _ShortURLService_CreateURL_Handler,
		},
		{
			MethodName: "BatchURL",
			Handler:    _ShortURLService_BatchURL_Handler,
		},
		{
			MethodName: "GetURL",
			Handler:    _ShortURLService_GetURL_Handler,
		},
		{
			MethodName: "GetUserURLs",
			Handler:    _ShortURLService_GetUserURLs_Handler,
		},
		{
			MethodName: "DeleteURL",
			Handler:    _ShortURLService_DeleteURL_Handler,
		},
		{
			MethodName: "StorageCheck",
			Handler:    _ShortURLService_StorageCheck_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/shorturl/v1/shorturl.proto",
}
