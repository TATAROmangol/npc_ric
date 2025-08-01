// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: table.proto

package table

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TableService_GetTable_FullMethodName = "/table.TableService/GetTable"
)

// TableServiceClient is the client API for TableService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TableServiceClient interface {
	GetTable(ctx context.Context, in *GetTableRequest, opts ...grpc.CallOption) (*GetTableResponse, error)
}

type tableServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTableServiceClient(cc grpc.ClientConnInterface) TableServiceClient {
	return &tableServiceClient{cc}
}

func (c *tableServiceClient) GetTable(ctx context.Context, in *GetTableRequest, opts ...grpc.CallOption) (*GetTableResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTableResponse)
	err := c.cc.Invoke(ctx, TableService_GetTable_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TableServiceServer is the server API for TableService service.
// All implementations must embed UnimplementedTableServiceServer
// for forward compatibility.
type TableServiceServer interface {
	GetTable(context.Context, *GetTableRequest) (*GetTableResponse, error)
	mustEmbedUnimplementedTableServiceServer()
}

// UnimplementedTableServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTableServiceServer struct{}

func (UnimplementedTableServiceServer) GetTable(context.Context, *GetTableRequest) (*GetTableResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTable not implemented")
}
func (UnimplementedTableServiceServer) mustEmbedUnimplementedTableServiceServer() {}
func (UnimplementedTableServiceServer) testEmbeddedByValue()                      {}

// UnsafeTableServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TableServiceServer will
// result in compilation errors.
type UnsafeTableServiceServer interface {
	mustEmbedUnimplementedTableServiceServer()
}

func RegisterTableServiceServer(s grpc.ServiceRegistrar, srv TableServiceServer) {
	// If the following call pancis, it indicates UnimplementedTableServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TableService_ServiceDesc, srv)
}

func _TableService_GetTable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTableRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TableServiceServer).GetTable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TableService_GetTable_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TableServiceServer).GetTable(ctx, req.(*GetTableRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TableService_ServiceDesc is the grpc.ServiceDesc for TableService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TableService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "table.TableService",
	HandlerType: (*TableServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTable",
			Handler:    _TableService_GetTable_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "table.proto",
}
