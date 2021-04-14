// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pbtodoapp

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

// TodoServiceClient is the client API for TodoService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TodoServiceClient interface {
	Create(ctx context.Context, in *TodoCreateRequest, opts ...grpc.CallOption) (*Todo, error)
	CreateMultiple(ctx context.Context, opts ...grpc.CallOption) (TodoService_CreateMultipleClient, error)
	Read(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Todo, error)
	ReadAll(ctx context.Context, in *ReadAllRequest, opts ...grpc.CallOption) (TodoService_ReadAllClient, error)
	Update(ctx context.Context, in *Todo, opts ...grpc.CallOption) (*Todo, error)
	Delete(ctx context.Context, in *Id, opts ...grpc.CallOption) (*DeleteResponse, error)
	DeleteMultiple(ctx context.Context, opts ...grpc.CallOption) (TodoService_DeleteMultipleClient, error)
}

type todoServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTodoServiceClient(cc grpc.ClientConnInterface) TodoServiceClient {
	return &todoServiceClient{cc}
}

func (c *todoServiceClient) Create(ctx context.Context, in *TodoCreateRequest, opts ...grpc.CallOption) (*Todo, error) {
	out := new(Todo)
	err := c.cc.Invoke(ctx, "/todoapp.TodoService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) CreateMultiple(ctx context.Context, opts ...grpc.CallOption) (TodoService_CreateMultipleClient, error) {
	stream, err := c.cc.NewStream(ctx, &TodoService_ServiceDesc.Streams[0], "/todoapp.TodoService/CreateMultiple", opts...)
	if err != nil {
		return nil, err
	}
	x := &todoServiceCreateMultipleClient{stream}
	return x, nil
}

type TodoService_CreateMultipleClient interface {
	Send(*TodoCreateRequest) error
	Recv() (*Todo, error)
	grpc.ClientStream
}

type todoServiceCreateMultipleClient struct {
	grpc.ClientStream
}

func (x *todoServiceCreateMultipleClient) Send(m *TodoCreateRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *todoServiceCreateMultipleClient) Recv() (*Todo, error) {
	m := new(Todo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *todoServiceClient) Read(ctx context.Context, in *Id, opts ...grpc.CallOption) (*Todo, error) {
	out := new(Todo)
	err := c.cc.Invoke(ctx, "/todoapp.TodoService/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) ReadAll(ctx context.Context, in *ReadAllRequest, opts ...grpc.CallOption) (TodoService_ReadAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &TodoService_ServiceDesc.Streams[1], "/todoapp.TodoService/ReadAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &todoServiceReadAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TodoService_ReadAllClient interface {
	Recv() (*Todo, error)
	grpc.ClientStream
}

type todoServiceReadAllClient struct {
	grpc.ClientStream
}

func (x *todoServiceReadAllClient) Recv() (*Todo, error) {
	m := new(Todo)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *todoServiceClient) Update(ctx context.Context, in *Todo, opts ...grpc.CallOption) (*Todo, error) {
	out := new(Todo)
	err := c.cc.Invoke(ctx, "/todoapp.TodoService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) Delete(ctx context.Context, in *Id, opts ...grpc.CallOption) (*DeleteResponse, error) {
	out := new(DeleteResponse)
	err := c.cc.Invoke(ctx, "/todoapp.TodoService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *todoServiceClient) DeleteMultiple(ctx context.Context, opts ...grpc.CallOption) (TodoService_DeleteMultipleClient, error) {
	stream, err := c.cc.NewStream(ctx, &TodoService_ServiceDesc.Streams[2], "/todoapp.TodoService/DeleteMultiple", opts...)
	if err != nil {
		return nil, err
	}
	x := &todoServiceDeleteMultipleClient{stream}
	return x, nil
}

type TodoService_DeleteMultipleClient interface {
	Send(*Id) error
	CloseAndRecv() (*DeleteResponse, error)
	grpc.ClientStream
}

type todoServiceDeleteMultipleClient struct {
	grpc.ClientStream
}

func (x *todoServiceDeleteMultipleClient) Send(m *Id) error {
	return x.ClientStream.SendMsg(m)
}

func (x *todoServiceDeleteMultipleClient) CloseAndRecv() (*DeleteResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(DeleteResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TodoServiceServer is the server API for TodoService service.
// All implementations must embed UnimplementedTodoServiceServer
// for forward compatibility
type TodoServiceServer interface {
	Create(context.Context, *TodoCreateRequest) (*Todo, error)
	CreateMultiple(TodoService_CreateMultipleServer) error
	Read(context.Context, *Id) (*Todo, error)
	ReadAll(*ReadAllRequest, TodoService_ReadAllServer) error
	Update(context.Context, *Todo) (*Todo, error)
	Delete(context.Context, *Id) (*DeleteResponse, error)
	DeleteMultiple(TodoService_DeleteMultipleServer) error
	mustEmbedUnimplementedTodoServiceServer()
}

// UnimplementedTodoServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTodoServiceServer struct {
}

func (UnimplementedTodoServiceServer) Create(context.Context, *TodoCreateRequest) (*Todo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedTodoServiceServer) CreateMultiple(TodoService_CreateMultipleServer) error {
	return status.Errorf(codes.Unimplemented, "method CreateMultiple not implemented")
}
func (UnimplementedTodoServiceServer) Read(context.Context, *Id) (*Todo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedTodoServiceServer) ReadAll(*ReadAllRequest, TodoService_ReadAllServer) error {
	return status.Errorf(codes.Unimplemented, "method ReadAll not implemented")
}
func (UnimplementedTodoServiceServer) Update(context.Context, *Todo) (*Todo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedTodoServiceServer) Delete(context.Context, *Id) (*DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedTodoServiceServer) DeleteMultiple(TodoService_DeleteMultipleServer) error {
	return status.Errorf(codes.Unimplemented, "method DeleteMultiple not implemented")
}
func (UnimplementedTodoServiceServer) mustEmbedUnimplementedTodoServiceServer() {}

// UnsafeTodoServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TodoServiceServer will
// result in compilation errors.
type UnsafeTodoServiceServer interface {
	mustEmbedUnimplementedTodoServiceServer()
}

func RegisterTodoServiceServer(s grpc.ServiceRegistrar, srv TodoServiceServer) {
	s.RegisterService(&TodoService_ServiceDesc, srv)
}

func _TodoService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TodoCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todoapp.TodoService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).Create(ctx, req.(*TodoCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_CreateMultiple_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TodoServiceServer).CreateMultiple(&todoServiceCreateMultipleServer{stream})
}

type TodoService_CreateMultipleServer interface {
	Send(*Todo) error
	Recv() (*TodoCreateRequest, error)
	grpc.ServerStream
}

type todoServiceCreateMultipleServer struct {
	grpc.ServerStream
}

func (x *todoServiceCreateMultipleServer) Send(m *Todo) error {
	return x.ServerStream.SendMsg(m)
}

func (x *todoServiceCreateMultipleServer) Recv() (*TodoCreateRequest, error) {
	m := new(TodoCreateRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _TodoService_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todoapp.TodoService/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).Read(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_ReadAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ReadAllRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TodoServiceServer).ReadAll(m, &todoServiceReadAllServer{stream})
}

type TodoService_ReadAllServer interface {
	Send(*Todo) error
	grpc.ServerStream
}

type todoServiceReadAllServer struct {
	grpc.ServerStream
}

func (x *todoServiceReadAllServer) Send(m *Todo) error {
	return x.ServerStream.SendMsg(m)
}

func _TodoService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Todo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todoapp.TodoService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).Update(ctx, req.(*Todo))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Id)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TodoServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/todoapp.TodoService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TodoServiceServer).Delete(ctx, req.(*Id))
	}
	return interceptor(ctx, in, info, handler)
}

func _TodoService_DeleteMultiple_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TodoServiceServer).DeleteMultiple(&todoServiceDeleteMultipleServer{stream})
}

type TodoService_DeleteMultipleServer interface {
	SendAndClose(*DeleteResponse) error
	Recv() (*Id, error)
	grpc.ServerStream
}

type todoServiceDeleteMultipleServer struct {
	grpc.ServerStream
}

func (x *todoServiceDeleteMultipleServer) SendAndClose(m *DeleteResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *todoServiceDeleteMultipleServer) Recv() (*Id, error) {
	m := new(Id)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TodoService_ServiceDesc is the grpc.ServiceDesc for TodoService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TodoService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "todoapp.TodoService",
	HandlerType: (*TodoServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _TodoService_Create_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _TodoService_Read_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _TodoService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _TodoService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CreateMultiple",
			Handler:       _TodoService_CreateMultiple_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "ReadAll",
			Handler:       _TodoService_ReadAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DeleteMultiple",
			Handler:       _TodoService_DeleteMultiple_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "protobuffer-files/todo.proto",
}
